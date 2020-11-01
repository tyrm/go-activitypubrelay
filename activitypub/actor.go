package activitypub

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/go-fed/httpsig"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type PublicKey struct {
	ID           string `json:"id,omitempty"`
	Owner        string `json:"owner,omitempty"`
	PublicKeyPem string `json:"publicKeyPem,omitempty"`
}

type Endpoints struct {
	SharedInbox string `json:"sharedInbox,omitempty"`
}

type Image struct {
	URL string `json:"url,omitempty"`
}

type Actor struct {
	Context           interface{} `json:"@context,omitempty"`
	Endpoints         Endpoints   `json:"endpoints,omitempty"`
	Followers         string      `json:"followers,omitempty"`
	Following         string      `json:"following,omitempty"`
	Icon              *Image      `json:"icon,omitempty"`
	Image             *Image      `json:"image,omitempty"`
	ID                string      `json:"id,omitempty"`
	Inbox             string      `json:"inbox,omitempty"`
	Name              string      `json:"name,omitempty"`
	Type              string      `json:"type,omitempty"`
	PublicKey         PublicKey   `json:"publicKey,omitempty"`
	Summary           string      `json:"summary,omitempty"`
	PreferredUsername string      `json:"preferredUsername,omitempty"`
	URL               string      `json:"url,omitempty"`
}

func (a *Actor) GetPublicKey() (*rsa.PublicKey, error) {
	pubPem, _ := pem.Decode([]byte(a.PublicKey.PublicKeyPem))
	if pubPem == nil {
		return nil, ErrPEMDecode
	}

	var parsedKey interface{}
	var err error
	if parsedKey, err = x509.ParsePKIXPublicKey(pubPem.Bytes); err != nil {
		return nil, err
	}

	var pubKey *rsa.PublicKey
	var ok bool
	if pubKey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, ErrPEMParse
	}

	return pubKey, nil
}

func (a *Actor) PushActivity(activity *Activity) error {
	// init signer
	prefs := []httpsig.Algorithm{httpsig.RSA_SHA512, httpsig.RSA_SHA256}
	digestAlgorithm := httpsig.DigestSha256
	headersToSign := []string{httpsig.RequestTarget, "date", "digest", "host"}

	signer, chosenAlgo, err := httpsig.NewSigner(prefs, digestAlgorithm, headersToSign, httpsig.Signature, 1800)
	if err != nil {
		return err
	}
	logger.Tracef("chosen signing algorithm: %s", chosenAlgo)

	// create body
	body, err := json.Marshal(activity)
	if err != nil {
		return err
	}

	// create http request
	req, err := http.NewRequest("POST", a.Inbox, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/activity+json")

	// add date header
	location,_ := time.LoadLocation("GMT")
	currentTime := time.Now().In(location)
	req.Header.Set("Date", currentTime.Format(time.RFC1123))

	// add host header
	inbox, err := url.Parse(a.Inbox)
	if err != nil {
		return err
	}
	req.Header.Set("Host", inbox.Host)

	// sign request
	err = signer.SignRequest(myPrivateKey, fmt.Sprintf("https://%s/actor#main-key", myAPHost), req, body)
	if err != nil {
		return err
	}

	// do request
	logger.Debugf("sending actor (%s): %s", a.Inbox, string(body))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	for k, v := range req.Header {
		logger.Tracef("Header field %q, Value %q", k, v)
	}

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	logger.Debugf("actor (%s) returned %d: %s", a.Inbox, res.StatusCode, string(body))

	return nil
}

func FetchActor(uri string, force bool) (*Actor, error) {
	// Check Cache
	if a, found := cActors.Get(uri); found {
		actor := a.(*Actor)
		return actor, nil
	}

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	actor := Actor{}
	err = json.Unmarshal(body, &actor)
	if err != nil {
		return nil, err
	}

	// Set Actor
	cActors.Set(uri, &actor, cache.DefaultExpiration)
	return &actor, nil
}

package activitypub

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"net/http"
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
	Icon              Image       `json:"icon,omitempty"`
	Image             Image       `json:"image,omitempty"`
	ID                string      `json:"id,omitempty"`
	Inbox             string      `json:"inbox,omitempty"`
	Name              string      `json:"name,omitempty"`
	Type              string      `json:"type,omitempty"`
	PublicKey         PublicKey   `json:"publicKey,omitempty"`
	Summary           string      `json:"summary,omitempty"`
	PreferredUsername string      `json:"preferredUsername,omitempty"`
	URL               string      `json:"url,omitempty"`
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
	err = json.Unmarshal([]byte(body), &actor)
	if err != nil {
		return nil, err
	}

	// Set Actor
	cActors.Set(uri, &actor, cache.DefaultExpiration)
	return &actor, nil
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

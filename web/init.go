package web

import (
	"crypto/rsa"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
	"litepub1/models"
	"net/http"
)

var (
	actor     models.Actor
	apHost string
	logger    *loggo.Logger
	webfinger models.WebFinger
)

func Init(APHost, APServiceName string, rsaKey *rsa.PrivateKey) error {
	newLogger := loggo.GetLogger("web")
	logger = &newLogger

	// Store Config
	apHost = APHost

	// Get RSA Text
	asn1Bytes, err := asn1.Marshal(rsaKey.PublicKey)
	if err != nil {
		return err
	}

	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: asn1Bytes,
		},
	)

	// Init actor
	actor = models.Actor{
		Context: "https://www.w3.org/ns/activitystreams",
		Endpoints: models.Endpoints{
			SharedInbox: fmt.Sprintf("https://%s/inbox", APHost),
		},
		Followers: fmt.Sprintf("https://%s/followers", APHost),
		Following: fmt.Sprintf("https://%s/following", APHost),
		Inbox:     fmt.Sprintf("https://%s/inbox", APHost),
		Name:      APServiceName,
		Type:      "Application",
		ID:        fmt.Sprintf("https://%s/actor", APHost),
		PublicKey: models.PublicKey{
			ID:           fmt.Sprintf("https://%s/actor#main-key", APHost),
			Owner:        fmt.Sprintf("https://%s/actor", APHost),
			PublicKeyPem: fmt.Sprintf("%s", pemdata),
		},
		Summary:           "GoActivityRelay bot",
		PreferredUsername: "relay",
		URL:               fmt.Sprintf("https://%s/actor", APHost),
	}

	// Init webfinger
	webfinger = models.WebFinger{
		Aliases: []string{fmt.Sprintf("https://%s/actor", APHost)},
		Links: []models.Link{
			{
				HRef: fmt.Sprintf("https://%s/actor", APHost),
				Rel: "self",
				Type: "application/activity+json",
			},
			{
				HRef: fmt.Sprintf("https://%s/actor", APHost),
				Rel: "self",
				Type: "application/activity+json",
			},
		},
		Subject: fmt.Sprintf("acct:relay@%s", APHost),
	}

	// Setup Router
	r := mux.NewRouter()
	r.Use(MiddlewareLogRequest)

	r.HandleFunc("/actor", HandleActor).Methods("GET")
	r.HandleFunc("/.well-known/webfinger", HandleWebFinger).Methods("GET")

	r.PathPrefix("/").HandlerFunc(HandleCatchAll) // Workaround to log all requests
	go func() {
		err := http.ListenAndServe(":9000", r)
		if err != nil {
			logger.Errorf("Could not start web server %s", err.Error())
		}
	}()

	return nil
}

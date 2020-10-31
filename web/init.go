package web

import (
	"crypto/rsa"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
	"litepub1/activitypub"
	"net/http"
)

var (
	actor             activitypub.Actor
	apHost            string
	logger            *loggo.Logger
	nodeinfoTemplate  Nodeinfo
	webfinger         WebFinger
	wellknownNodeinfo WellknownNodeinfo
)

type Link struct {
	HRef string `json:"href,omitempty"`
	Rel  string `json:"rel,omitempty"`
	Type string `json:"type,omitempty"`
}

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
	actor = activitypub.Actor{
		Context: "https://www.w3.org/ns/activitystreams",
		Endpoints: activitypub.Endpoints{
			SharedInbox: fmt.Sprintf("https://%s/inbox", APHost),
		},
		Followers: fmt.Sprintf("https://%s/followers", APHost),
		Following: fmt.Sprintf("https://%s/following", APHost),
		Inbox:     fmt.Sprintf("https://%s/inbox", APHost),
		Name:      APServiceName,
		Type:      "Application",
		ID:        fmt.Sprintf("https://%s/actor", APHost),
		PublicKey: activitypub.PublicKey{
			ID:           fmt.Sprintf("https://%s/actor#main-key", APHost),
			Owner:        fmt.Sprintf("https://%s/actor", APHost),
			PublicKeyPem: fmt.Sprintf("%s", pemdata),
		},
		Summary:           "GoActivityRelay bot",
		PreferredUsername: "relay",
		URL:               fmt.Sprintf("https://%s/actor", APHost),
	}

	// Init nodeinfo
	nodeinfoTemplate = Nodeinfo{
		OpenRegistration: true,
		Protocols:        []string{"activitypub"},
		Services: Services{
			Inbound:  []string{},
			Outbound: []string{},
		},
		Software: Software{
			Name:    "goactivityrelay",
			Version: "0.0",
		},
		Usage: Usage{
			LocalPosts: 0,
			Users: UsageUsers{
				Total: 1,
			},
		},
		Version: "2.0",
	}

	// Init webfinger
	webfinger = WebFinger{
		Aliases: []string{fmt.Sprintf("https://%s/actor", APHost)},
		Links: []Link{
			{
				HRef: fmt.Sprintf("https://%s/actor", APHost),
				Rel:  "self",
				Type: "application/activity+json",
			},
			{
				HRef: fmt.Sprintf("https://%s/actor", APHost),
				Rel:  "self",
				Type: "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"",
			},
		},
		Subject: fmt.Sprintf("acct:relay@%s", APHost),
	}

	// Init wellknownNodeinfo
	wellknownNodeinfo = WellknownNodeinfo{
		Links: []Link{
			{
				Rel:  "http://nodeinfo.diaspora.software/ns/schema/2.0",
				HRef: fmt.Sprintf("https://%s/nodeinfo/2.0.json", APHost),
			},
		},
	}

	// Setup Router
	r := mux.NewRouter()
	r.Use(MiddlewareHttpSignatures)
	r.Use(MiddlewareLogRequest)

	r.HandleFunc("/actor", HandleActor).Methods("GET")
	r.HandleFunc("/inbox", HandleInbox).Methods("POST")
	r.HandleFunc("/nodeinfo/2.0.json", HandleNodeinfo20).Methods("GET")
	r.HandleFunc("/.well-known/nodeinfo", HandleWellKnownNodeInfo).Methods("GET")
	r.HandleFunc("/.well-known/webfinger", HandleWellKnownWebFinger).Methods("GET")

	r.PathPrefix("/").HandlerFunc(HandleCatchAll) // Workaround to log all requests
	go func() {
		err := http.ListenAndServe(":9000", r)
		if err != nil {
			logger.Errorf("Could not start web server %s", err.Error())
		}
	}()

	return nil
}

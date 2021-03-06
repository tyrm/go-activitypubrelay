package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//noinspection GoNameStartsWithPackageName
type WebFinger struct {
	Aliases []string `json:"aliases,omitempty"`
	Links   []Link   `json:"links,omitempty"`
	Subject string   `json:"subject,omitempty"`
}

func HandleWellKnownWebFinger(w http.ResponseWriter, r *http.Request) {
	resource, ok := r.URL.Query()["resource"]
	if !ok || len(resource[0]) < 1 {
		http.Error(w, "Url Param 'key' is missing", http.StatusBadRequest)
		return
	}

	if resource[0] != fmt.Sprintf("acct:relay@%s", apHost) {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	actor, err := json.Marshal(&myWebfinger)
	if err != nil {
		logger.Warningf("Could not marshal JSON: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/activity+json")
	_, err = w.Write(actor)
	if err != nil {
		logger.Warningf("Could not write response: %s", err.Error())
		return
	}
}

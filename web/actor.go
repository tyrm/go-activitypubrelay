package web

import (
	"encoding/json"
	"net/http"
)

func HandleActor(w http.ResponseWriter, r *http.Request) {
	actor, err := json.Marshal(&myActor)
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

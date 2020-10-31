package web

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	"io/ioutil"
	"litepub1/activitypub"
	"net/http"
)

func HandleInbox(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("could not read body: %s", err.Error())
		logger.Debugf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	logger.Tracef("request body: %s", reqBody)

	// check for validation
	validated := context.Get(r, "validated")
	if validated == false {
		msg := "signature validation failed"
		logger.Debugf(msg)
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}


	// decode json
	decoder := json.NewDecoder(r.Body)
	var req activitypub.Activity
	err = decoder.Decode(&req)
	if err != nil {
		msg := fmt.Sprintf("could not decode json: %s", err.Error())
		logger.Debugf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// check for actor in body
	if req.Actor == "" {
		msg := "no actor in message"
		logger.Debugf(msg)
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}


	//go activitypub.ProcessInbox()


	nodeinfo, err := json.Marshal(&wellknownNodeinfo)
	if err != nil {
		logger.Warningf("Could not marshal JSON: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/activity+json")
	_, err = w.Write(nodeinfo)
	if err != nil {
		logger.Warningf("Could not write response: %s", err.Error())
		return
	}
}

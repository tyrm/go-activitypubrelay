package web

import (
	"fmt"
	"litepub1/activitypub"
	"litepub1/models"
	"net/http"
	"net/url"
)

func HandleInbox(w http.ResponseWriter, r *http.Request) {
	// check for validation
	validated := r.Context().Value(SignatureValidKey).(bool)
	if validated != true {
		msg := "signature validation failed"
		logger.Debugf(msg)
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}

	// get activity
	activity := r.Context().Value(ActivityKey).(*activitypub.Activity)
	actor, err := url.Parse(activity.Actor)
	if err != nil {
		msg := fmt.Sprintf("could not parse myActor url (%s): %s", activity.Actor, err.Error())
		logger.Debugf(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// check that instance is a
	if activity.Type != "Follow" && models.GetFollowedInstanceExists(actor.Host) {
		msg := fmt.Sprintf("instance (%s) not following relay", actor.Host)
		logger.Debugf(msg)
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}

	go activitypub.ProcessInbox(activity)

	w.Header().Add("Content-Type", "application/activity+json")
	_, err = w.Write([]byte("{}"))
	if err != nil {
		logger.Warningf("Could not write response: %s", err.Error())
		return
	}
}

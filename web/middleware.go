package web

import (
	"encoding/json"
	"fmt"
	"github.com/go-fed/httpsig"
	"github.com/gorilla/context"
	"litepub1/activitypub"
	"net/http"
)

func MiddlewareLogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		client := r.RemoteAddr
		if r.Header.Get("X-Forwarded-For") != "" {
			client = r.Header.Get("X-Forwarded-For")
		}

		validated := context.Get(r, "validated")

		logger.Debugf("%s \"%s %s\" (%s) validated: %v", client, r.Method, r.RequestURI, r.Header.Get("User-Agent"), validated)

		next.ServeHTTP(w, r)
	})
}

func MiddlewareHttpSignatures(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context.Set(r, "validated", false)

		if r.Header.Get("signature") != "" && r.Method == "POST" {
			logger.Tracef("http signature detected. parsing json body")

			// sdecode json
			decoder := json.NewDecoder(r.Body)
			var req activitypub.Activity
			err := decoder.Decode(&req)
			if err != nil {
				msg := fmt.Sprintf("could not decode json: %s", err.Error())
				logger.Debugf(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}

			// check for actor in body
			if req.Actor == "" {
				msg := "signature check failed, no actor in message"
				logger.Debugf(msg)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}

			// get actor data
			actorData, err := activitypub.FetchActor(req.Actor, false)
			if err != nil {
				msg := fmt.Sprintf("could not retrieve actor: %s", err.Error())
				logger.Debugf(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}

			logger.Tracef("found actor '%s'", req.Actor)

			verifier, err := httpsig.NewVerifier(r)
			if err != nil {
				msg := fmt.Sprintf("could not initiate verifier")
				logger.Warningf(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}

			pk, err := actorData.GetPublicKey()
			if err != nil {
				msg := fmt.Sprintf("could not get actor's key: %s", err.Error())
				logger.Warningf(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}

			var algo httpsig.Algorithm = httpsig.RSA_SHA256
			err = verifier.Verify(pk, algo)
			if err != nil {
				msg := fmt.Sprintf("message signature verification failed: %s", err.Error())
				logger.Warningf(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}

			context.Set(r, "validated", true)
		}

		next.ServeHTTP(w, r)
	})
}
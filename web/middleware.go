package web

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-fed/httpsig"
	"litepub1/activitypub"
	"net/http"
)

type contextKey int

func MiddlewareLogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		client := r.RemoteAddr
		if r.Header.Get("X-Forwarded-For") != "" {
			client = r.Header.Get("X-Forwarded-For")
		}

		//validated := context.Get(r, "validated")
		validated := r.Context().Value(SignatureValidKey).(bool)

		logger.Debugf("%s \"%s %s\" (%s) validated: %v", client, r.Method, r.RequestURI, r.Header.Get("User-Agent"), validated)

		next.ServeHTTP(w, r)
	})
}

func MiddlewareHttpSignatures(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), SignatureValidKey, false)

		if r.Header.Get("signature") != "" && r.Method == "POST" {
			logger.Tracef("http signature detected. parsing json body")

			// decode json
			decoder := json.NewDecoder(r.Body)
			var activity activitypub.Activity
			err := decoder.Decode(&activity)
			if err != nil {
				msg := fmt.Sprintf("could not decode json: %s", err.Error())
				logger.Debugf(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}

			// Save activity to context
			ctx = context.WithValue(ctx, ActivityKey, &activity)

			// check for myActor in body
			if activity.Actor == "" {
				msg := "signature check failed, no actor in message"
				logger.Debugf(msg)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}

			// get actor data
			actorData, err := activitypub.FetchActor(activity.Actor, false)
			if err != nil {
				msg := fmt.Sprintf("could not retrieve actor: %s", err.Error())
				logger.Warningf(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}

			ctx = context.WithValue(ctx, ActorKey, actorData)

			verifier, err := httpsig.NewVerifier(r)
			if err != nil {
				msg := fmt.Sprintf("could not initiate verifier: %s", err.Error())
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

			var algo = httpsig.RSA_SHA256
			if err := verifier.Verify(pk, algo); err != nil {
				msg := fmt.Sprintf("message signature verification failed: %s", err.Error())
				logger.Warningf(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}

			ctx = context.WithValue(ctx, SignatureValidKey, true)
		}

		rWithSignature := r.WithContext(ctx)

		next.ServeHTTP(w, rWithSignature)
	})
}

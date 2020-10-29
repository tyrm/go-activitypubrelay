package activitypub

import (
	"encoding/json"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"litepub1/models"
	"net/http"
)

func FetchActor(uri string, force bool) (*models.Actor, error) {
	// Check Cache
	if a, found := cActors.Get(uri); found {
		actor := a.(*models.Actor)
		return actor, nil
	}

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	actor := models.Actor{}
	err = json.Unmarshal([]byte(body), &actor)
	if err != nil {
		return nil, err
	}

	// Set Actor
	cActors.Set(uri, &actor, cache.DefaultExpiration)
	return &actor, nil
}
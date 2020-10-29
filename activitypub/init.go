package activitypub

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var cActors *cache.Cache

func Init() {
	// init cache
	cActors = cache.New(1*time.Hour, 10*time.Minute)
}
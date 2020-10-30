package activitypub

import (
	"errors"
	"github.com/patrickmn/go-cache"
	"time"
)

var cActors *cache.Cache

var (
	ErrPEMDecode = errors.New("unable to decode pem")
	ErrPEMParse = errors.New("unable to parse pem")
)

func Init() {
	// init cache
	cActors = cache.New(1*time.Hour, 10*time.Minute)
}

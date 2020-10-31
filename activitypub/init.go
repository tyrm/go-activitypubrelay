package activitypub

import (
	"errors"
	"github.com/juju/loggo"
	"github.com/patrickmn/go-cache"
	"time"
)

var cActors *cache.Cache

var logger *loggo.Logger

var (
	ErrPEMDecode = errors.New("unable to decode pem")
	ErrPEMParse = errors.New("unable to parse pem")
)

func Init() {
	newLogger := loggo.GetLogger("models")
	logger = &newLogger

	// init cache
	cActors = cache.New(1*time.Hour, 10*time.Minute)
}

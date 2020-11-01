package activitypub

import (
	"crypto/rsa"
	"errors"
	"github.com/juju/loggo"
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	cActors      *cache.Cache
	logger       *loggo.Logger
	myPrivateKey *rsa.PrivateKey
	myAPHost   string
)

// Errors
var (
	ErrPEMDecode = errors.New("unable to decode pem")
	ErrPEMParse  = errors.New("unable to parse pem")
)

func Init(apHost string, rsaKey *rsa.PrivateKey) {
	newLogger := loggo.GetLogger("activitypub")
	logger = &newLogger

	// save config
	myPrivateKey = rsaKey
	myAPHost = apHost

	// init cache
	cActors = cache.New(1*time.Hour, 10*time.Minute)
}

package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"io/ioutil"
	"litepub1/activitypub"
	"litepub1/models"
	"litepub1/web"
	"os"
	"os/signal"
	"syscall"
)

var logger *loggo.Logger

func main() {
	config := CollectConfig()

	// Init Logging
	newLogger := loggo.GetLogger("main")
	logger = &newLogger

	err := loggo.ConfigureLoggers(config.LoggerConfig)
	if err != nil {
		logger.Errorf("Error configuring Logger: %s", err.Error())
		return
	}
	_, err = loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr))
	if err != nil {
		logger.Errorf("Error configuring Color Logger: %s", err.Error())
		return
	}

	logger.Infof("Starting LitePub Relay")

	// Load RSA Key
	logger.Debugf("Loading RSA Key")
	rsaKey, err := LoadRSAKey(config.RSAPrivateKey, config.RSAPublicKey)
	if err != nil {
		logger.Errorf("Could not read RSA Key: %s", err.Error())
		return
	}

	err = models.Init(config.DBType, config.DBConnString)
	if err != nil {
		logger.Errorf("Could init models: %s", err.Error())
		return
	}

	activitypub.Init(config.APHost, rsaKey)

	err = web.Init(config.APHost, config.APServiceName, rsaKey, config.AllowlistEnabled)
	if err != nil {
		logger.Errorf("Could init web: %s", err.Error())
		return
	}

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	nch := make(chan os.Signal)
	signal.Notify(nch, syscall.SIGINT, syscall.SIGTERM)
	logger.Infof("%s", <-nch)
}

func LoadRSAKey(privKeyPath, publicKeyPath string) (*rsa.PrivateKey, error) {
	// Read Private Key
	priv, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		return nil, err
	}
	privPem, _ := pem.Decode(priv)
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPem.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPem.Bytes); err != nil { // note this returns type `interface{}`
			return nil, err
		}
	}
	var privateKey *rsa.PrivateKey
	var ok bool
	privateKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, err
	}

	// Read Public Key
	pub, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}
	pubPem, _ := pem.Decode(pub)
	if pubPem == nil {
		return nil, err
	}
	if parsedKey, err = x509.ParsePKIXPublicKey(pubPem.Bytes); err != nil {
		return nil, err
	}

	var pubKey *rsa.PublicKey
	if pubKey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, err
	}

	// Combine
	privateKey.PublicKey = *pubKey

	return privateKey, nil

}

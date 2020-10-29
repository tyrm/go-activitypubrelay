package main

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	APHost        string
	APSerivceName string

	LoggerConfig string

	RSAPrivateKey string
	RSAPublicKey  string
}

func CollectConfig() *Config {
	var missingEnv []string
	var config Config

	// AP_HOST
	config.APHost = os.Getenv("AP_HOST")
	if config.APHost == "" {
		missingEnv = append(missingEnv, "AP_HOST")
	}

	// AP_SERVICE_NAME
	var envAPServiceName = os.Getenv("AP_SERVICE_NAME")
	if envAPServiceName == "" {
		config.APSerivceName = "GoActivityRelay"
	} else {
		config.APSerivceName = envAPServiceName
	}

	// LOG_LEVEL
	var envLoggerLevel = os.Getenv("LOG_LEVEL")
	if envLoggerLevel == "" {
		config.LoggerConfig = "<root>=INFO"
	} else {
		config.LoggerConfig = fmt.Sprintf("<root>=%s", strings.ToUpper(envLoggerLevel))
	}

	// RSA_PRIVATE_KEY
	var envRSAPrivateKey = os.Getenv("RSA_PRIVATE_KEY")
	if envRSAPrivateKey == "" {
		config.RSAPrivateKey = "privatekey.pem"
	} else {
		config.RSAPrivateKey = envRSAPrivateKey
	}

	// RSA_PRIVATE_KEY
	var envRSAPublicKey = os.Getenv("RSA_PUBLIC_KEY")
	if envRSAPublicKey == "" {
		config.RSAPublicKey = "publickey.pem"
	} else {
		config.RSAPublicKey = envRSAPublicKey
	}

	// Validation
	if len(missingEnv) > 0 {
		msg := fmt.Sprintf("Environment variables missing: %v", missingEnv)
		panic(fmt.Sprint(msg))
	}

	return &config
}

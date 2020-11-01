package main

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	AllowlistEnabled bool

	APHost        string
	APServiceName string

	DBType       string
	DBConnString string

	LoggerConfig string

	RSAPrivateKey string
	RSAPublicKey  string
}

func CollectConfig() *Config {
	var missingEnv []string
	var config Config

	// ALLOWLIST_ENABLED
	var envAllowlistEnabled = os.Getenv("ALLOWLIST_ENABLED")
	if strings.ToLower(envAllowlistEnabled) == "true" {
		config.AllowlistEnabled = true
	} else {
		config.AllowlistEnabled = false
	}

	// AP_HOST
	config.APHost = os.Getenv("AP_HOST")
	if config.APHost == "" {
		missingEnv = append(missingEnv, "AP_HOST")
	}

	// AP_SERVICE_NAME
	var envAPServiceName = os.Getenv("AP_SERVICE_NAME")
	if envAPServiceName == "" {
		config.APServiceName = "GoActivityRelay"
	} else {
		config.APServiceName = envAPServiceName
	}

	// DB_TYPE
	var envDBType = os.Getenv("DB_TYPE")
	if envDBType == "" {
		config.DBType = "sqlite"
	} else {
		config.DBType = envDBType
	}

	// DB_CONN_STRING
	var envDBConnString = os.Getenv("DB_CONN_STRING")
	if envDBConnString == "" {
		config.DBConnString = "litepub.db"
	} else {
		config.DBConnString = envDBConnString
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

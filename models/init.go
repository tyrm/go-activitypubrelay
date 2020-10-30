package models

import (
	"errors"
	"github.com/juju/loggo"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var logger *loggo.Logger

var (
	ErrInvalidDBType = errors.New("does not exist")
)

func Init(dbType, dbConnString string) error {
	newLogger := loggo.GetLogger("models")
	logger = &newLogger

	if dbType == "sqlite" {
		dbClient, err := gorm.Open(sqlite.Open(dbConnString), &gorm.Config{})
		if err != nil {
			return err
		}

		db = dbClient
	} else if dbType == "postgres" {
		dbClient, err := gorm.Open(postgres.Open(dbConnString), &gorm.Config{})
		if err != nil {
			return err
		}

		db = dbClient
	} else if dbType == "mysql" {
		dbClient, err := gorm.Open(mysql.Open(dbConnString), &gorm.Config{})
		if err != nil {
			return err
		}

		db = dbClient
	} else {
		return ErrInvalidDBType
	}

	// Migrate the schema
	err := db.AutoMigrate(&Peer{})
	if err != nil {
		return err
	}

	return nil
}

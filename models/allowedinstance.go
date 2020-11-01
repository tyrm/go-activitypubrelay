package models

import (
	"gorm.io/gorm"
)

type AllowedInstance struct {
	gorm.Model
	Hostname string
}
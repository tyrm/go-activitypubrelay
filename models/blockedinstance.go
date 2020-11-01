package models

import (
	"gorm.io/gorm"
)

type BlockedInstance struct {
	gorm.Model
	Hostname string
}

func BlockedInstanceExists(hostname string) bool {
	var count int64
	db.Model(&BlockedInstance{}).Where("hostname = ?", hostname).Count(&count)
	if count > 0 {
		return true
	}

	return false
}
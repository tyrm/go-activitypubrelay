package models

import (
	"gorm.io/gorm"
)

type AllowedInstance struct {
	gorm.Model
	Hostname string
}

func GetAllowedPeers() (*[]AllowedInstance, error) {
	var allowedPeers []AllowedInstance

	result := db.Order("url asc").Find(&allowedPeers)
	if result.Error != nil {
		return nil, result.Error
	}

	return &allowedPeers, nil
}

func GetAllowedPeersExists(url string) bool {
	var count int64
	db.Model(&AllowedInstance{}).Where("hostname = ?", url).Count(&count)
	if count > 0 {
		return true
	}

	return false
}
package models

import (
	"gorm.io/gorm"
)

type AllowedPeer struct {
	gorm.Model
	URL string
}

func GetAllowedPeers() (*[]AllowedPeer, error) {
	var allowedPeers []AllowedPeer

	result := db.Order("url asc").Find(&allowedPeers)
	if result.Error != nil {
		return nil, result.Error
	}

	return &allowedPeers, nil
}

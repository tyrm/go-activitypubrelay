package models

import (
	"gorm.io/gorm"
)

type BlockedInstance struct {
	gorm.Model
	Hostname string
}

func GetBlockedPeers() (*[]BlockedInstance, error) {
	var blockedPeers []BlockedInstance

	result := db.Order("url asc").Find(&blockedPeers)
	if result.Error != nil {
		return nil, result.Error
	}

	return &blockedPeers, nil
}

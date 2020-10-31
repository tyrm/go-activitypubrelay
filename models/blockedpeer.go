package models

import (
	"gorm.io/gorm"
)

type BlockedPeer struct {
	gorm.Model
	Hostname string
}

func GetBlockedPeers() (*[]BlockedPeer, error) {
	var blockedPeers []BlockedPeer

	result := db.Order("url asc").Find(&blockedPeers)
	if result.Error != nil {
		return nil, result.Error
	}

	return &blockedPeers, nil
}

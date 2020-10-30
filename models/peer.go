package models

import (
	"gorm.io/gorm"
)

type Peer struct {
	gorm.Model
	URL string
}

func GetPeers() (*[]Peer, error) {
	var peers []Peer

	result := db.Order("url asc").Find(&peers)
	if result.Error != nil {
		return nil, result.Error
	}

	return &peers, nil
}

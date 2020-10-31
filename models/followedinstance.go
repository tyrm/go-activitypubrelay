package models

import (
	"gorm.io/gorm"
)

type FollowedInstance struct {
	gorm.Model
	Hostname string
}

func GetFollowedInstance() (*[]FollowedInstance, error) {
	var peers []FollowedInstance

	result := db.Order("url asc").Find(&peers)
	if result.Error != nil {
		return nil, result.Error
	}

	return &peers, nil
}

func GetFollowedInstanceExists(url string) bool {
	var count int64
	db.Model(&FollowedInstance{}).Where("hostname = ?", url).Count(&count)
	if count > 0 {
		return true
	}

	return false
}
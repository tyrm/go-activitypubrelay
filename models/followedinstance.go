package models

import (
	"gorm.io/gorm"
)

type FollowedInstance struct {
	gorm.Model
	Hostname string
}

func CreateFollowedInstance(instance *FollowedInstance) error {
	result := db.Create(&instance)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func ReadFollowedInstances() (*[]FollowedInstance, error) {
	var peers []FollowedInstance

	result := db.Order("hostname asc").Find(&peers)
	if result.Error != nil {
		return nil, result.Error
	}

	return &peers, nil
}

func FollowedInstanceExists(hostname string) bool {
	var count int64
	db.Model(&FollowedInstance{}).Where("hostname = ?", hostname).Count(&count)
	if count > 0 {
		return true
	}

	return false
}
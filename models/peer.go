package models

import "gorm.io/gorm"

type Peer struct {
	gorm.Model
	URL string
}

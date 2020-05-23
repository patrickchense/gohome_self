package models

import (
	"github.com/jinzhu/gorm"
)

type Site struct {
	gorm.Model
	SiteStatus int    `gorm:"default:0"`
	Name       string `gorm:"not null"`
	Url        string `gorm:"not null"`
	UserId     uint32 `gorm:"default:null"`
}

func (site *Site) GetSiteStatusAsString() string {
	switch site.SiteStatus {
	case 0:
		return "ok"
	case 1:
		return "stopped"
	default:
		return "unknown"
	}
}

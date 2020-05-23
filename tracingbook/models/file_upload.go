package models

import "github.com/jinzhu/gorm"

type FileUpload struct {
	gorm.Model
	Filename     string
	FilePath     string
	OriginalName string
	FileSize     uint
}

// Scopes, not used
func TagImages(db *gorm.DB) *gorm.DB {
	return db.Where("type = ?", "TagImage")
}

// db.Scopes(CategoryImages, ProductImages).Find(&images)

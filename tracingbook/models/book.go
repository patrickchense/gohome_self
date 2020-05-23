package models

import (
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type Book struct {
	gorm.Model
	Name        string `gorm:"size:280;not null"`
	Description string `gorm:"not null"`
	Slug        string `gorm:"unique_index;not null"`
	Author      string `gorm:"not null"`
	Latest      string `gorm:"not null"`
	UserId      uint32 `gorm:"default null"`
}

func (book *Book) BeforeSave() (err error) {
	book.Slug = slug.Make(book.Name)
	return
}

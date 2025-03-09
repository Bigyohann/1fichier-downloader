package models

import (
	"gorm.io/gorm"
)

type File struct {
	Model
	Url         string `json:"url"`
	Filename    string `json:"filename"`
	Size        int    `json:"size"`
	Date        string `json:"date"`
	ContentType string `json:"content-type"`
	Checksum    string `json:"checksum"`
	Downloaded  bool   `json:"downloaded"`
}

func InitDbFileModel(db *gorm.DB) {
	db.AutoMigrate(&File{})
}

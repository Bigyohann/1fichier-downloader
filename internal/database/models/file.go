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
	ContentType string `json:"contentType"`
	Checksum    string `json:"checksum"`
	Downloaded  bool   `json:"downloaded"`
	Status      string `json:"status"`
}

func InitDbFileModel(db *gorm.DB) {
	db.AutoMigrate(&File{})
}

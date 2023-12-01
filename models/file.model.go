package models

import (
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Filename     string
	FilePath     string
	OriginalName string
	FileSize     uint
}

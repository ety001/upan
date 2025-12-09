package models

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Code      string         `gorm:"type:varchar(6);uniqueIndex;not null" json:"code"`
	Path      string         `gorm:"type:varchar(255);not null" json:"path"`
	Filename  string         `gorm:"type:varchar(255);not null" json:"filename"`
	Mime      string         `gorm:"type:varchar(100);not null" json:"mime"`
	Size      int64          `gorm:"not null" json:"size"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (File) TableName() string {
	return "files"
}


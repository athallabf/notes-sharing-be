package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Note struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;not null" json:"id"`
	Title     string         `gorm:"not null" json:"title"`
	Content   string         `json:"content"`
	FilePath  string         `json:"file_path,omitempty"`
	FileURL   string         `gorm:"-" json:"file_url,omitempty"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"-"`
	User      *User          `gorm:"foreignKey:UserID;references:ID" json:"-"`
	CreatedAt time.Time      `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (note *Note) BeforeCreate(_ *gorm.DB) error {
	note.ID = uuid.New()
	return nil
}

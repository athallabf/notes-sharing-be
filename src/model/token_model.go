package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Token struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
	Token     string    `gorm:"not null"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_tokens_user_id_type" json:"-"`
	Type      string    `gorm:"not null;index:idx_tokens_user_id_type" json:"-"`
	Expires   time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli"`
	UpdatedAt time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli"`
	User      *User     `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

func (token *Token) BeforeCreate(_ *gorm.DB) error {
	token.ID = uuid.New()
	return nil
}

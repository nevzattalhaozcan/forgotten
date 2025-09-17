package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Annotation struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	BookID     uint           `json:"book_id" gorm:"not null;foreignKey:BookID"`
	AuthorID   uint           `json:"author_id" gorm:"not null;foreignKey:AuthorID"`
	Quote      string         `json:"quote" gorm:"type:text;not null"`
	PageNumber *int           `json:"page,omitempty"`
	Thoughts   string         `json:"thoughts,omitempty" gorm:"type:text"`
	IsPublic   bool           `json:"is_public" gorm:"default:true"`
	LikesCount int            `json:"likes_count" gorm:"default:0"`
	Tags       pq.StringArray `json:"tags" gorm:"type:text[]"`
	Book       Book           `json:"book" gorm:"foreignKey:BookID"`
	Author     User           `json:"author" gorm:"foreignKey:AuthorID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

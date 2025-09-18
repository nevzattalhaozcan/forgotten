package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID            uint     `json:"id" gorm:"primaryKey"`
	Title         string   `json:"title" gorm:"size:255;not null"`
	User          *string  `json:"user,omitempty" gorm:"size:255"`
	CoverURL      *string  `json:"cover_url,omitempty" gorm:"type:text"`
	Genre         *string  `json:"genre,omitempty" gorm:"size:100"`
	Pages         *int     `json:"pages,omitempty"`
	PublishedYear *int     `json:"published_year,omitempty"`
	ISBN          *string  `json:"isbn,omitempty" gorm:"size:20;uniqueIndex"`
	Description   *string  `json:"description,omitempty" gorm:"type:text"`
	Rating        *float32 `json:"rating,omitempty" gorm:"type:decimal(2,1)"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

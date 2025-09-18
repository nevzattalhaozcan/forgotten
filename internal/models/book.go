package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID            uint     `json:"id" gorm:"primaryKey"`
	Title         string   `json:"title" gorm:"size:255;not null"`
	Author       *string  `json:"author,omitempty" gorm:"size:255"`
	CoverURL      *string  `json:"cover_url,omitempty" gorm:"type:text"`
	Genre         *string  `json:"genre,omitempty" gorm:"size:100"`
	Pages         *int     `json:"pages,omitempty"`
	PublishedYear *int     `json:"published_year,omitempty"`
	ISBN          *string  `json:"isbn,omitempty" gorm:"size:20;uniqueIndex"`
	Description   *string  `json:"description,omitempty" gorm:"type:text"`
	Rating        *float32 `json:"rating,omitempty" gorm:"type:decimal(2,1)" validate:"omitempty,gte=0,lte=5"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateBookRequest struct {
	Title         string  `json:"title" validate:"required,min=1,max=255"`
	Author       *string `json:"author,omitempty" validate:"omitempty,min=1,max=255"`
	CoverURL      *string `json:"cover_url,omitempty" validate:"omitempty,url"`
	Genre         *string `json:"genre,omitempty" validate:"omitempty,min=1,max=100"`
	Pages         *int    `json:"pages,omitempty" validate:"omitempty,gte=1"`
	PublishedYear *int    `json:"published_year,omitempty" validate:"omitempty,gte=0,lte=2100"`
	ISBN          *string `json:"isbn,omitempty" validate:"omitempty,isbn"`
	Description   *string `json:"description,omitempty" validate:"omitempty,min=1"`
}

type UpdateBookRequest struct {
	Title         *string  `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Author       *string  `json:"author,omitempty" validate:"omitempty,min=1,max=255"`
	CoverURL      *string  `json:"cover_url,omitempty" validate:"omitempty,url"`
	Genre         *string  `json:"genre,omitempty" validate:"omitempty,min=1,max=100"`
	Pages         *int     `json:"pages,omitempty" validate:"omitempty,gte=1"`
	PublishedYear *int     `json:"published_year,omitempty" validate:"omitempty,gte=0,lte=2100"`
	ISBN          *string  `json:"isbn,omitempty" validate:"omitempty,isbn"`
	Description   *string  `json:"description,omitempty" validate:"omitempty,min=1"`
	Rating        *float32 `json:"rating,omitempty" validate:"omitempty,gte=0,lte=5"`
}

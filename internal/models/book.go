package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID uint `json:"id" gorm:"primaryKey"`

	// External API fields
	ExternalID *string `json:"external_id,omitempty" gorm:"index"`
	Source     *string `json:"source,omitempty" gorm:"size:50"`

	// Book metadata
	Title         string   `json:"title" gorm:"size:255;not null"`
	Author        *string  `json:"author,omitempty" gorm:"size:255"`
	CoverURL      *string  `json:"cover_url,omitempty" gorm:"type:text"`
	Genre         *string  `json:"genre,omitempty" gorm:"size:100"`
	Pages         *int     `json:"pages,omitempty"`
	PublishedYear *int     `json:"published_year,omitempty"`
	ISBN          *string  `json:"isbn,omitempty" gorm:"size:20;uniqueIndex"`
	Description   *string  `json:"description,omitempty" gorm:"type:text"`
	Rating        *float32 `json:"rating,omitempty" gorm:"type:decimal(2,1)" validate:"omitempty,gte=0,lte=5"`

	// Platform-specific analytics (simplified)
	ReadCount      int      `json:"read_count" gorm:"default:0"`
	LocalRating    *float32 `json:"local_rating,omitempty" gorm:"type:decimal(2,1)"`
	RatingCount    int      `json:"rating_count" gorm:"default:0"`
	IsClubFavorite bool     `json:"is_club_favorite" gorm:"default:false"`
	IsTrending     bool     `json:"is_trending" gorm:"default:false"`

	// Cache management
	CachedAt     *time.Time `json:"cached_at,omitempty"`
	LastAccessed *time.Time `json:"last_accessed,omitempty"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type ExternalBook struct {
	ExternalID    string   `json:"external_id"`
	Source        string   `json:"source"`
	Title         string   `json:"title"`
	Author        *string  `json:"author,omitempty"`
	CoverURL      *string  `json:"cover_url,omitempty"`
	Genre         *string  `json:"genre,omitempty"`
	Pages         *int     `json:"pages,omitempty"`
	PublishedYear *int     `json:"published_year,omitempty"`
	ISBN          *string  `json:"isbn,omitempty"`
	Description   *string  `json:"description,omitempty"`
	Rating        *float32 `json:"rating,omitempty"`
}

func (eb *ExternalBook) ToBook() *Book {
	now := time.Now()
	return &Book{
		ExternalID:    &eb.ExternalID,
		Source:        &eb.Source,
		Title:         eb.Title,
		Author:        eb.Author,
		CoverURL:      eb.CoverURL,
		Genre:         eb.Genre,
		Pages:         eb.Pages,
		PublishedYear: eb.PublishedYear,
		ISBN:          eb.ISBN,
		Description:   eb.Description,
		Rating:        eb.Rating,
		CachedAt:      &now,
		LastAccessed:  &now,
	}
}

type CreateBookRequest struct {
	Title         string  `json:"title" validate:"required,min=1,max=255"`
	Author        *string `json:"author,omitempty" validate:"omitempty,min=1,max=255"`
	CoverURL      *string `json:"cover_url,omitempty" validate:"omitempty,url"`
	Genre         *string `json:"genre,omitempty" validate:"omitempty,min=1,max=100"`
	Pages         *int    `json:"pages,omitempty" validate:"omitempty,gte=1"`
	PublishedYear *int    `json:"published_year,omitempty" validate:"omitempty,gte=0,lte=2100"`
	ISBN          *string `json:"isbn,omitempty" validate:"omitempty,isbn"`
	Description   *string `json:"description,omitempty" validate:"omitempty,min=1"`
	ExternalID    *string `json:"external_id,omitempty"`
	Source        *string `json:"source,omitempty"`
}

type UpdateBookRequest struct {
	Title         *string  `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Author        *string  `json:"author,omitempty" validate:"omitempty,min=1,max=255"`
	CoverURL      *string  `json:"cover_url,omitempty" validate:"omitempty,url"`
	Genre         *string  `json:"genre,omitempty" validate:"omitempty,min=1,max=100"`
	Pages         *int     `json:"pages,omitempty" validate:"omitempty,gte=1"`
	PublishedYear *int     `json:"published_year,omitempty" validate:"omitempty,gte=0,lte=2100"`
	ISBN          *string  `json:"isbn,omitempty" validate:"omitempty,isbn"`
	Description   *string  `json:"description,omitempty" validate:"omitempty,min=1"`
	Rating        *float32 `json:"rating,omitempty" validate:"omitempty,gte=0,lte=5"`
}

type BookSearchRequest struct {
	Query      string  `json:"query" form:"q" validate:"required,min=1"`
	Limit     int     `json:"limit" form:"limit" validate:"omitempty,gte=1,lte=100"`
	Source string `json:"source" form:"source" validate:"omitempty,oneof=local external all"`
} 

type BookResponse struct {
	ID            uint     `json:"id"`
	ExternalID    *string  `json:"external_id,omitempty"`
	Source        *string  `json:"source,omitempty"`
	Title         string   `json:"title"`
	Author        *string  `json:"author,omitempty"`
	CoverURL      *string  `json:"cover_url,omitempty"`
	Genre         *string  `json:"genre,omitempty"`
	Pages         *int     `json:"pages,omitempty"`
	PublishedYear *int     `json:"published_year,omitempty"`
	ISBN          *string  `json:"isbn,omitempty"`
	Description   *string  `json:"description,omitempty"`
	Rating        *float32 `json:"rating,omitempty"`
	LocalRating   *float32 `json:"local_rating,omitempty"`

	ReadCount      int  `json:"read_count"`
	RatingCount    int  `json:"rating_count"`
	IsClubFavorite bool `json:"is_club_favorite"`
	IsTrending     bool `json:"is_trending"`

	ClubCount *int `json:"club_count,omitempty"`

	UserRating    *float32 `json:"user_rating,omitempty"`
	ReadingStatus *string  `json:"reading_status,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b *Book) ToResponse() BookResponse {
	return BookResponse{
		ID:             b.ID,
		ExternalID:     b.ExternalID,
		Source:         b.Source,
		Title:          b.Title,
		Author:         b.Author,
		CoverURL:       b.CoverURL,
		Genre:          b.Genre,
		Pages:          b.Pages,
		PublishedYear:  b.PublishedYear,
		ISBN:           b.ISBN,
		Description:    b.Description,
		Rating:         b.Rating,
		LocalRating:    b.LocalRating,
		ReadCount:      b.ReadCount,
		RatingCount:    b.RatingCount,
		IsClubFavorite: b.IsClubFavorite,
		IsTrending:     b.IsTrending,
		CreatedAt:      b.CreatedAt,
		UpdatedAt:      b.UpdatedAt,
	}
}

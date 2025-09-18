package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AnnotationLike struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	UserID       uint       `json:"user_id"`
	AnnotationID uint       `json:"annotation_id"`
	User         User       `json:"user" gorm:"foreignKey:UserID"`
	Annotation   Annotation `json:"annotation" gorm:"foreignKey:AnnotationID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Annotation struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	BookID     uint           `json:"book_id" gorm:"not null;foreignKey:BookID"`
	UserID     uint           `json:"user_id" gorm:"not null;foreignKey:UserID"`
	Quote      string         `json:"quote" gorm:"type:text;not null"`
	PageNumber *int           `json:"page,omitempty"`
	Thoughts   string         `json:"thoughts,omitempty" gorm:"type:text"`
	IsPublic   bool           `json:"is_public" gorm:"default:true"`
	LikesCount int            `json:"likes_count" gorm:"default:0"`
	Likes      []AnnotationLike `json:"likes,omitempty" gorm:"foreignKey:AnnotationID"`
	Tags       pq.StringArray `json:"tags" gorm:"type:text[]"`
	Book       Book           `json:"book" gorm:"foreignKey:BookID"`
	User       User           `json:"user" gorm:"foreignKey:UserID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateAnnotationRequest struct {
	BookID     uint     `json:"book_id" validate:"required"`
	Quote      string   `json:"quote" validate:"required,min=1"`
	PageNumber *int     `json:"page,omitempty" validate:"omitempty,gte=1"`
	Thoughts   *string  `json:"thoughts,omitempty" validate:"omitempty,min=1"`
	IsPublic   *bool    `json:"is_public,omitempty"`
	Tags       []string `json:"tags,omitempty" validate:"omitempty,dive,max=50"`
}

type UpdateAnnotationRequest struct {
	Quote      *string  `json:"quote,omitempty" validate:"omitempty,min=1"`
	PageNumber *int     `json:"page,omitempty" validate:"omitempty,gte=1"`
	Thoughts   *string  `json:"thoughts,omitempty" validate:"omitempty,min=1"`
	IsPublic   *bool    `json:"is_public,omitempty"`
	Tags       []string `json:"tags,omitempty" validate:"omitempty,dive,max=50"`
}

type LikeAnnotationRequest struct {
	AnnotationID uint `json:"annotation_id" validate:"required"`
}

type UnlikeAnnotationRequest struct {
	AnnotationID uint `json:"annotation_id" validate:"required"`
}
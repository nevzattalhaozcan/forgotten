package models

import (
	"time"

	"gorm.io/gorm"
)

type CommentLike struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	UserID    uint    `json:"user_id"`
	CommentID uint    `json:"comment_id"`
	User      User    `json:"user" gorm:"foreignKey:UserID"`
	Comment   Comment `json:"comment" gorm:"foreignKey:CommentID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Comment struct {
	ID         uint          `json:"id" gorm:"primaryKey"`
	PostID     uint          `json:"post_id" gorm:"not null; foreignKey:PostID"`
	UserID     uint          `json:"user_id" gorm:"not null;foreignKey:UserID"`
	Content    string        `json:"content" gorm:"type:text;not null"`
	LikesCount int           `json:"likes_count" gorm:"default:0"`
	User       User          `json:"user" gorm:"foreignKey:UserID"`
	Likes      []CommentLike `json:"likes,omitempty" gorm:"foreignKey:CommentID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateCommentRequest struct {
	PostID  uint   `json:"post_id" validate:"required"`
	Content string `json:"content" validate:"required,min=1"`
}

type UpdateCommentRequest struct {
	Content *string `json:"content,omitempty" validate:"omitempty,min=1"`
}

type LikeCommentRequest struct {
	CommentID uint `json:"comment_id" validate:"required"`
}

type UnlikeCommentRequest struct {
	CommentID uint `json:"comment_id" validate:"required"`
}
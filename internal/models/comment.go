package models

import (
	"time"

	"gorm.io/gorm"
)

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

type CommentLike struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	UserID    uint    `json:"user_id"`
	CommentID uint    `json:"comment_id"`
	User      User    `json:"user" gorm:"foreignKey:UserID"`
	Comment   Comment `json:"comment" gorm:"foreignKey:CommentID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

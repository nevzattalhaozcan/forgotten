package models

import (
	"time"

	"gorm.io/gorm"
)

type CommentLike struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	UserID    uint    `json:"user_id"`
	CommentID uint    `json:"comment_id"`
	User      User    `json:"user" gorm:"foreignKey:UserID" swaggerignore:"true"`
	Comment   Comment `json:"comment" gorm:"foreignKey:CommentID" swaggerignore:"true"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Comment struct {
	ID         uint          `json:"id" gorm:"primaryKey"`
	PostID     uint          `json:"post_id" gorm:"not null; foreignKey:PostID"`
	UserID     uint          `json:"user_id" gorm:"not null;foreignKey:UserID"`
	Content    string        `json:"content" gorm:"type:text;not null"`
	LikesCount int           `json:"likes_count" gorm:"default:0"`
	User       User          `json:"user" gorm:"foreignKey:UserID" swaggerignore:"true"`
	Likes      []CommentLike `json:"likes,omitempty" gorm:"foreignKey:CommentID" swaggerignore:"true"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1"`
}

type UpdateCommentRequest struct {
	Content *string `json:"content,omitempty" validate:"omitempty,min=1"`
}

type CommentResponse struct {
	ID         uint          `json:"id"`
	PostID     uint          `json:"post_id"`
	UserID     uint          `json:"user_id"`
	Content    string        `json:"content"`
	LikesCount int           `json:"likes_count"`
	User       User          `json:"user" swaggerignore:"true"`
	Likes      []CommentLike `json:"likes,omitempty" swaggerignore:"true"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentLikeResponse struct {
	ID        uint         `json:"id"`
	User      UserResponse `json:"user"`
	CreatedAt time.Time    `json:"created_at"`
}

func (cl CommentLike) ToResponse() CommentLikeResponse {
	return CommentLikeResponse{
		ID:        cl.ID,
		User:      cl.User.ToResponse(),
		CreatedAt: cl.CreatedAt,
	}
}

func (c *Comment) ToResponse() CommentResponse {
	return CommentResponse{
		ID:         c.ID,
		PostID:     c.PostID,
		UserID:     c.UserID,
		Content:    c.Content,
		LikesCount: c.LikesCount,
		User:       c.User,
		Likes:      c.Likes,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}

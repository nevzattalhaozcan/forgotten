package models

import (
	"time"

	"gorm.io/gorm"
)

type PostLike struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id"`
	PostID uint `json:"post_id"`
	User   User `json:"user" gorm:"foreignKey:UserID"`
	Post   Post `json:"post" gorm:"foreignKey:PostID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Post struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	Title         string     `json:"title" gorm:"size:255;not null"`
	Content       string     `json:"content" gorm:"type:text;not null"`
	Type          string     `json:"type" gorm:"not null" validate:"required,oneof=discussion announcement event poll review" default:"discussion"`
	IsPinned      bool       `json:"is_pinned" gorm:"default:false"`
	LikesCount    int        `json:"likes_count" gorm:"default:0"`
	CommentsCount int        `json:"comments_count" gorm:"default:0"`
	ViewsCount    int        `json:"views_count" gorm:"default:0"`
	UserID        uint       `json:"user_id"`
	ClubID        uint       `json:"club_id"`
	User          User       `json:"user" gorm:"foreignKey:UserID"`
	Comments      []Comment  `json:"comments,omitempty" gorm:"foreignKey:PostID"`
	Likes         []PostLike `json:"likes,omitempty" gorm:"foreignKey:PostID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreatePostRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=255"`
	Content string `json:"content" validate:"required,min=1"`
	Type    string `json:"type" validate:"required,oneof=discussion announcement event poll review"`
	ClubID  uint   `json:"club_id" validate:"required"`
}

type UpdatePostRequest struct {
	Title    *string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Content  *string `json:"content,omitempty" validate:"omitempty,min=1"`
	Type     *string `json:"type,omitempty" validate:"omitempty,oneof=discussion announcement event poll review"`
	IsPinned *bool   `json:"is_pinned,omitempty"`
}

type LikePostRequest struct {
	PostID uint `json:"post_id" validate:"required"`
}

type UnlikePostRequest struct {
	PostID uint `json:"post_id" validate:"required"`
}

type PostResponse struct {
	ID            uint       `json:"id"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	Type          string     `json:"type"`
	IsPinned      bool       `json:"is_pinned"`
	LikesCount    int        `json:"likes_count"`
	CommentsCount int        `json:"comments_count"`
	ViewsCount    int        `json:"views_count"`
	UserID        uint       `json:"user_id"`
	ClubID        uint       `json:"club_id"`
	User          User       `json:"user"`
	Comments      []Comment  `json:"comments,omitempty"`
	Likes         []PostLike `json:"likes,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Post) ToResponse() PostResponse {
	return PostResponse{
		ID:            p.ID,
		Title:         p.Title,
		Content:       p.Content,
		Type:          p.Type,
		IsPinned:      p.IsPinned,
		LikesCount:    p.LikesCount,
		CommentsCount: p.CommentsCount,
		ViewsCount:    p.ViewsCount,
		UserID:        p.UserID,
		ClubID:        p.ClubID,
		User:          p.User,
		Comments:      p.Comments,
		Likes:         p.Likes,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

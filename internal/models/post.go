package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID            uint         `json:"id" gorm:"primaryKey"`
	Title         string       `json:"title" gorm:"size:255;not null"`
	Content       string       `json:"content" gorm:"type:text;not null"`
	Type          string       `json:"type" gorm:"not null" validate:"required,oneof=discussion announcement post poll review annotation" default:"discussion"`
	TypeData      PostTypeData `json:"type_data,omitempty" gorm:"type:jsonb"`
	IsPinned      bool         `json:"is_pinned" gorm:"default:false"`
	LikesCount    int          `json:"likes_count" gorm:"default:0"`
	CommentsCount int          `json:"comments_count" gorm:"default:0"`
	ViewsCount    int          `json:"views_count" gorm:"default:0"`
	UserID        uint         `json:"user_id"`
	ClubID        uint         `json:"club_id"`

	Club     Club       `json:"club" gorm:"foreignKey:ClubID" swaggerignore:"true"`
	User     User       `json:"user" gorm:"foreignKey:UserID" swaggerignore:"true"`
	Comments []Comment  `json:"comments,omitempty" gorm:"foreignKey:PostID" swaggerignore:"true"`
	Likes    []PostLike `json:"likes,omitempty" gorm:"foreignKey:PostID" swaggerignore:"true"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type PostTypeData json.RawMessage

func (ptd *PostTypeData) Scan(value interface{}) error {
	if value == nil {
		*ptd = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("cannot scan into PostTypeData")
	}
	*ptd = bytes
	return nil
}

func (ptd PostTypeData) Value() (driver.Value, error) {
	if len(ptd) == 0 {
		return nil, nil
	}
	return []byte(ptd), nil
}

type PostLike struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id" gorm:"uniqueIndex:idx_user_post_like"`
	PostID uint `json:"post_id" gorm:"uniqueIndex:idx_user_post_like"`
	User   User `json:"user" gorm:"foreignKey:UserID" swaggerignore:"true"`
	Post   Post `json:"post" gorm:"foreignKey:PostID" swaggerignore:"true"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PollVote struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	PostID   uint   `json:"post_id" gorm:"index"`
	UserID   uint   `json:"user_id" gorm:"index"`
	OptionID string `json:"option_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReviewData struct {
	BookID     uint    `json:"book_id" validate:"required"`
	Rating     float32 `json:"rating" validate:"required,gte=1,lte=5"`
	BookTitle  string  `json:"book_title,omitempty"`
	BookAuthor string  `json:"book_author,omitempty"`
}

type PollData struct {
	Question      string       `json:"question" validate:"required,min=1,max=500"`
	Options       []PollOption `json:"options" validate:"required,min=2,max=10,dive"`
	AllowMultiple bool         `json:"allow_multiple"`
	ExpiresAt     *time.Time   `json:"expires_at,omitempty"`
}

type PollOption struct {
	ID    string `json:"id"`
	Text  string `json:"text" validate:"required,min=1,max=200"`
	Votes int    `json:"votes"`
}

type AnnotationData struct {
	BookID     uint   `json:"book_id" validate:"required"`
	Page       *int   `json:"page,omitempty" validate:"omitempty,gte=1"`
	Chapter    *int   `json:"chapter,omitempty" validate:"omitempty,gte=1"`
	Quote      string `json:"quote,omitempty" validate:"omitempty,max=1000"`
	BookTitle  string `json:"book_title,omitempty"`
	BookAuthor string `json:"book_author,omitempty"`
}

type PostData struct {
	PostID      uint   `json:"post_id,omitempty"`
	PostTitle   string `json:"post_title,omitempty"`
	PostContent string `json:"post_content,omitempty"`
}

type UserSummary struct {
	ID        uint    `json:"id"`
	Username  string  `json:"username"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type ClubSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type PostSummary struct {
	ID            uint         `json:"id" gorm:"column:id"`
	Title         string       `json:"title" gorm:"column:title"`
	Content       string       `json:"content" gorm:"column:content"`
	Type          string       `json:"type" gorm:"column:type"`
	TypeData      PostTypeData `json:"type_data,omitempty" gorm:"type:jsonb"`
	IsPinned      bool         `json:"is_pinned" gorm:"column:is_pinned"`
	LikesCount    int          `json:"likes_count" gorm:"column:likes_count"`
	CommentsCount int          `json:"comments_count" gorm:"column:comments_count"`
	ViewsCount    int          `json:"views_count" gorm:"column:views_count"`
	UserID        uint         `json:"user_id" gorm:"column:post_user_id"`
	ClubID        *uint        `json:"club_id" gorm:"column:post_club_id"`
	User          UserSummary  `json:"user"`
	Club          *ClubSummary `json:"club,omitempty"`
	CreatedAt     time.Time    `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time    `json:"updated_at" gorm:"column:updated_at"`
}

func (p *Post) GetReviewData() (*ReviewData, error) {
	if p.Type != "review" || len(p.TypeData) == 0 {
		return nil, nil
	}
	var data ReviewData
	err := json.Unmarshal(p.TypeData, &data)
	return &data, err
}

func (p *Post) GetPollData() (*PollData, error) {
	if p.Type != "poll" || len(p.TypeData) == 0 {
		return nil, nil
	}
	var data PollData
	err := json.Unmarshal(p.TypeData, &data)
	return &data, err
}

func (p *Post) GetAnnotationData() (*AnnotationData, error) {
	if p.Type != "annotation" || len(p.TypeData) == 0 {
		return nil, nil
	}
	var data AnnotationData
	err := json.Unmarshal(p.TypeData, &data)
	return &data, err
}

func (p *Post) GetPostData() (*PostData, error) {
	if p.Type != "post" || len(p.TypeData) == 0 {
		return nil, nil
	}
	var data PostData
	err := json.Unmarshal(p.TypeData, &data)
	return &data, err
}

type CreatePostRequest struct {
	Title    string      `json:"title" validate:"required,min=1,max=255"`
	Content  string      `json:"content" validate:"required,min=1"`
	Type     string      `json:"type" validate:"required,oneof=discussion announcement post poll review annotation"`
	ClubID   uint        `json:"club_id" validate:"required"`
	TypeData interface{} `json:"type_data,omitempty"`
}

type UpdatePostRequest struct {
	Title    *string     `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Content  *string     `json:"content,omitempty" validate:"omitempty,min=1"`
	Type     *string     `json:"type,omitempty" validate:"omitempty,oneof=discussion announcement post poll review annotation"`
	ClubID   *uint       `json:"club_id,omitempty" validate:"omitempty"`
	IsPinned *bool       `json:"is_pinned,omitempty"`
	TypeData interface{} `json:"type_data,omitempty"`
}

type PostResponse struct {
	ID            uint        `json:"id"`
	Title         string      `json:"title"`
	Content       string      `json:"content"`
	Type          string      `json:"type"`
	TypeData      interface{} `json:"type_data,omitempty"`
	IsPinned      bool        `json:"is_pinned"`
	LikesCount    int         `json:"likes_count"`
	CommentsCount int         `json:"comments_count"`
	ViewsCount    int         `json:"views_count"`
	UserID        uint        `json:"user_id"`
	ClubID        uint        `json:"club_id"`
	User          User        `json:"user" swaggerignore:"true"`
	Comments      []Comment   `json:"comments,omitempty" swaggerignore:"true"`
	Likes         []PostLike  `json:"likes,omitempty" swaggerignore:"true"`

	UserVoted bool     `json:"user_voted,omitempty"`
	UserVotes []string `json:"user_votes,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PollVoteRequest struct {
	OptionIDs []string `json:"option_ids" validate:"required,min=1"`
}

type PostLikeResponse struct {
	ID        uint         `json:"id"`
	User      UserResponse `json:"user"`
	CreatedAt time.Time    `json:"created_at"`
}

func (pl PostLike) ToResponse() PostLikeResponse {
	return PostLikeResponse{
		ID:        pl.ID,
		User:      pl.User.ToResponse(),
		CreatedAt: pl.CreatedAt,
	}
}

func (p *Post) ToResponse() PostResponse {
	response := PostResponse{
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

	if len(p.TypeData) > 0 {
		var typeData interface{}
		if err := json.Unmarshal(p.TypeData, &typeData); err == nil {
			response.TypeData = typeData
		} else {
			response.TypeData = nil
		}
	}
	return response
}

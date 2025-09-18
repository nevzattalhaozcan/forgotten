package models

import (
	"encoding/json"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CurrentBook struct {
	Title    string  `json:"title"`
	Author     *string `json:"author,omitempty"`
	CoverURL *string `json:"cover_url,omitempty"`
	BookID   *uint   `json:"book_id,omitempty"`
	Progress *int    `json:"progress,omitempty"`
}

type NextMeeting struct {
	Date     *time.Time `json:"date,omitempty"`
	Location *string    `json:"location,omitempty"`
	Topic    *string    `json:"topic,omitempty"`
}

type ClubMembership struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	UserID   uint      `json:"user_id"`
	ClubID   uint      `json:"club_id"`
	Role     string    `json:"role" gorm:"default:'member'"`
	JoinedAt time.Time `json:"joined_at" gorm:"autoCreateTime"`
	User     User      `json:"user" gorm:"foreignKey:UserID"`
}

type Club struct {
	ID            uint             `json:"id" gorm:"primaryKey"`
	Name          string           `json:"name" gorm:"size:100;not null;unique"`
	Description   string           `json:"description" gorm:"type:text"`
	Location      *string          `json:"location" gorm:"size:255"`
	Genre         *string          `json:"genre" gorm:"size:100"`
	CoverImageURL *string          `json:"cover_image_url" gorm:"type:text"`
	IsPrivate     bool             `json:"is_private" gorm:"default:false"`
	MaxMembers    int              `json:"max_members" gorm:"default:100"`
	MembersCount  int              `json:"members_count" gorm:"default:0"`
	Rating        float32          `json:"rating" gorm:"default:0"`
	Tags          pq.StringArray   `json:"tags" gorm:"type:text[]"`
	OwnerID       uint             `json:"owner_id"`
	CurrentBook   json.RawMessage  `json:"current_book" gorm:"type:jsonb"`
	NextMeeting   json.RawMessage  `json:"next_meeting" gorm:"type:jsonb"`
	Owner         User             `json:"owner" gorm:"foreignKey:OwnerID"`
	Moderators    []User           `json:"moderators" gorm:"many2many:club_moderators;"`
	Members       []ClubMembership `json:"members" gorm:"foreignKey:ClubID"`
	Posts         []Post           `json:"posts,omitempty" gorm:"foreignKey:ClubID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateClubRequest struct {
	Name          string         `json:"name" validate:"required,min=3,max=100"`
	Description   string         `json:"description" validate:"max=1000"`
	Location      *string        `json:"location" validate:"omitempty,max=255"`
	Genre         *string        `json:"genre" validate:"omitempty,max=100"`
	CoverImageURL *string        `json:"cover_image_url" validate:"omitempty,url"`
	IsPrivate     bool           `json:"is_private"`
	MaxMembers    int            `json:"max_members" validate:"gte=1,lte=1000"`
	Tags          pq.StringArray `json:"tags" validate:"dive,max=50"`
}

type UpdateClubRequest struct {
	Name          *string         `json:"name" validate:"omitempty,min=3,max=100"`
	Description   *string         `json:"description" validate:"omitempty,max=1000"`
	Location      *string         `json:"location" validate:"omitempty,max=255"`
	Genre         *string         `json:"genre" validate:"omitempty,max=100"`
	CoverImageURL *string         `json:"cover_image_url" validate:"omitempty,url"`
	IsPrivate     *bool           `json:"is_private"`
	MaxMembers    *int            `json:"max_members" validate:"omitempty,gte=1,lte=1000"`
	Tags          *pq.StringArray `json:"tags" validate:"omitempty,dive,max=50"`
	CurrentBook   *CurrentBook    `json:"current_book"`
	NextMeeting   *NextMeeting    `json:"next_meeting"`
}

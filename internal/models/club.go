package models

import (
	"encoding/json"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CurrentBook struct {
	Title    string  `json:"title"`
	Author   *string `json:"author,omitempty"`
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
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id"`
	ClubID     uint      `json:"club_id"`
	Role       string    `json:"role" gorm:"default:'member'"`
	IsApproved bool      `json:"is_approved" gorm:"default:false"`
	JoinedAt   time.Time `json:"joined_at" gorm:"autoCreateTime"`
	User       User      `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Club       Club      `json:"-" gorm:"foreignKey:ClubID;constraint:OnDelete:CASCADE"`
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
	RatingsCount  int              `json:"ratings_count" gorm:"default:0"`
	Tags          pq.StringArray   `json:"tags" gorm:"type:text[]"`
	OwnerID       *uint            `json:"owner_id"`
	CurrentBook   json.RawMessage  `json:"current_book" gorm:"type:jsonb"`
	NextMeeting   json.RawMessage  `json:"next_meeting" gorm:"type:jsonb"`
	Owner         User             `json:"owner" gorm:"foreignKey:OwnerID;constraint:OnDelete:SET NULL"`
	Moderators    []User           `json:"moderators" gorm:"many2many:club_moderators;"`
	Members       []ClubMembership `json:"members" gorm:"foreignKey:ClubID;constraint:OnDelete:CASCADE"`
	Posts         []Post           `json:"posts,omitempty" gorm:"foreignKey:ClubID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type ClubRating struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ClubID    uint      `json:"club_id" gorm:"index;not null"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	Rating    float32   `json:"rating" gorm:"not null"`
	Comment   *string   `json:"comment" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type RateClubRequest struct {
	Rating  float32 `json:"rating" validate:"required,gte=1,lte=5"`
	Comment *string `json:"comment" validate:"omitempty,max=1000"`
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

type UpdateClubMembershipRequest struct {
	UserID     *uint   `json:"user_id" validate:"omitempty"`
	Role       *string `json:"role" validate:"omitempty,oneof=member moderator admin"`
	IsApproved *bool   `json:"is_approved" validate:"omitempty"`
}

type UpdateClubRatingRequest struct {
	Rating float32 `json:"rating" validate:"required,gte=0,lte=5"`
}

type ClubMembershipResponse struct {
	ID         uint         `json:"id"`
	UserID     uint         `json:"user_id"`
	ClubID     uint         `json:"club_id"`
	Role       string       `json:"role"`
	IsApproved bool         `json:"is_approved"`
	JoinedAt   time.Time    `json:"joined_at"`
	User       UserResponse `json:"user"`
}

type ClubResponse struct {
	ID            uint             `json:"id"`
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	Location      *string          `json:"location,omitempty"`
	Genre         *string          `json:"genre,omitempty"`
	CoverImageURL *string          `json:"cover_image_url,omitempty"`
	IsPrivate     bool             `json:"is_private"`
	MaxMembers    int              `json:"max_members"`
	MembersCount  int              `json:"members_count"`
	Rating        float32          `json:"rating"`
	RatingsCount  int              `json:"ratings_count"`
	Tags          pq.StringArray   `json:"tags"`
	OwnerID       uint             `json:"owner_id"`
	Owner         UserResponse     `json:"owner"`
	CurrentBook   *CurrentBook     `json:"current_book,omitempty"`
	NextMeeting   *NextMeeting     `json:"next_meeting,omitempty"`
	Members       []ClubMembership `json:"members,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

func (c *Club) ToResponse() ClubResponse {
	var currentBook *CurrentBook
	if len(c.CurrentBook) > 0 {
		var cb CurrentBook
		if err := json.Unmarshal(c.CurrentBook, &cb); err == nil {
			currentBook = &cb
		}
	}

	var nextMeeting *NextMeeting
	if len(c.NextMeeting) > 0 {
		var nm NextMeeting
		if err := json.Unmarshal(c.NextMeeting, &nm); err == nil {
			nextMeeting = &nm
		}
	}

	members := make([]ClubMembership, len(c.Members))
	copy(members, c.Members)

	return ClubResponse{
		ID:            c.ID,
		Name:          c.Name,
		Description:   c.Description,
		Location:      c.Location,
		Genre:         c.Genre,
		CoverImageURL: c.CoverImageURL,
		IsPrivate:     c.IsPrivate,
		MaxMembers:    c.MaxMembers,
		MembersCount:  c.MembersCount,
		Rating:        c.Rating,
		Tags:          c.Tags,
		OwnerID: func() uint {
			if c.OwnerID != nil {
				return *c.OwnerID
			}
			return 0
		}(),
		Owner:       c.Owner.ToResponse(),
		CurrentBook: currentBook,
		NextMeeting: nextMeeting,
		Members:     members,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func (cm *ClubMembership) ToResponse() ClubMembershipResponse {
	return ClubMembershipResponse{
		ID:         cm.ID,
		UserID:     cm.UserID,
		ClubID:     cm.ClubID,
		Role:       cm.Role,
		IsApproved: cm.IsApproved,
		JoinedAt:   cm.JoinedAt,
		User:       cm.User.ToResponse(),
	}
}

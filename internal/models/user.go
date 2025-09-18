package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Username     string `json:"username" gorm:"uniqueIndex;not null" validate:"required,min=3,max=50"`
	Email        string `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	PasswordHash string `json:"-" gorm:"not null"`
	FirstName    string `json:"first_name" validate:"min=2,max=50"`
	LastName     string `json:"last_name" validate:"min=2,max=50"`
	IsActive     bool   `json:"is_active" gorm:"default:true"`
	Role         string `json:"role" validate:"required,oneof=admin user moderator support superuser" gorm:"default:'user'"`

	AvatarURL      *string        `json:"avatar_url" gorm:"type:text"`
	Location       *string        `json:"location" gorm:"size:255"`
	FavoriteGenres pq.StringArray `json:"favorite_genres" gorm:"type:text[]"`
	Bio            *string        `json:"bio" gorm:"type:text"`
	ReadingGoal    int            `json:"reading_goal" gorm:"default:0"`
	BooksRead      int            `json:"books_read" gorm:"default:0"`
	Badges         pq.StringArray `json:"badges" gorm:"type:text[]"`
	IsOnline       bool           `json:"is_online" gorm:"default:false"`
	LastSeen       *time.Time     `json:"last_seen"`

	OwnedClubs      []Club           `json:"owned_clubs,omitempty" gorm:"foreignKey:OwnerID"`
	ClubMemberships []ClubMembership `json:"club_memberships,omitempty" gorm:"foreignKey:UserID"`
	Posts           []Post           `json:"posts,omitempty" gorm:"foreignKey:UserID"`
	Comments        []Comment        `json:"comments,omitempty" gorm:"foreignKey:UserID"`
	Annotations     []Annotation     `json:"annotations,omitempty" gorm:"foreignKey:UserID"`
	PostLikes       []PostLike       `json:"post_likes,omitempty" gorm:"foreignKey:UserID"`
	CommentLikes    []CommentLike    `json:"comment_likes,omitempty" gorm:"foreignKey:UserID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsActive  bool   `json:"is_active"`
	Role      string `json:"role"`

	AvatarURL      *string        `json:"avatar_url"`
	Location       *string        `json:"location"`
	FavoriteGenres pq.StringArray `json:"favorite_genres"`
	Bio            *string        `json:"bio"`
	ReadingGoal    *int           `json:"reading_goal"`
	BooksRead      int            `json:"books_read"`
	Badges         pq.StringArray `json:"badges"`
	IsOnline       bool           `json:"is_online"`
	LastSeen       *time.Time     `json:"last_seen"`

	OwnedClubs      []Club           `json:"owned_clubs,omitempty"`
	ClubMemberships []ClubMembership `json:"club_memberships,omitempty"`
	Posts           []Post           `json:"posts,omitempty"`
	Comments        []Comment        `json:"comments,omitempty"`
	Annotations     []Annotation     `json:"annotations,omitempty"`
	PostLikes       []PostLike       `json:"post_likes,omitempty"`
	CommentLikes    []CommentLike    `json:"comment_likes,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Role      string `json:"role" validate:"omitempty,oneof=admin user moderator support superuser" gorm:"default:'user'"`

	AvatarURL      string   `json:"avatar_url"`
	Location       string   `json:"location"`
	FavoriteGenres []string `json:"favorite_genres"`
	Bio            string   `json:"bio"`
	ReadingGoal    int      `json:"reading_goal"`
}

type UpdateUserRequest struct {
	Username  *string `json:"username" validate:"omitempty,min=3,max=50"`
	Email     *string `json:"email" validate:"omitempty,email"`
	Password  *string `json:"password" validate:"omitempty,min=6"`
	FirstName *string `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName  *string `json:"last_name" validate:"omitempty,min=2,max=50"`
	Role      *string `json:"role" validate:"omitempty,oneof=admin user moderator support superuser"`
	IsActive  *bool   `json:"is_active"`

	AvatarURL      *string   `json:"avatar_url" validate:"omitempty,url"`
	Location       *string   `json:"location" validate:"omitempty,max=255"`
	FavoriteGenres *[]string `json:"favorite_genres"`
	Bio            *string   `json:"bio" validate:"omitempty"`
	ReadingGoal    *int      `json:"reading_goal" validate:"omitempty,gte=0"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:              u.ID,
		Username:        u.Username,
		Email:           u.Email,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		IsActive:        u.IsActive,
		Role:            u.Role,
		AvatarURL:       u.AvatarURL,
		Location:        u.Location,
		FavoriteGenres:  u.FavoriteGenres,
		Bio:             u.Bio,
		ReadingGoal:     &u.ReadingGoal,
		BooksRead:       u.BooksRead,
		Badges:          u.Badges,
		IsOnline:        u.IsOnline,
		LastSeen:        u.LastSeen,
		OwnedClubs:      u.OwnedClubs,
		ClubMemberships: u.ClubMemberships,
		Posts:           u.Posts,
		Comments:        u.Comments,
		Annotations:     u.Annotations,
		PostLikes:       u.PostLikes,
		CommentLikes:    u.CommentLikes,
		CreatedAt:       u.CreatedAt,
	}
}

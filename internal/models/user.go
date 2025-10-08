package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserPreferences map[string]interface{}

func (up *UserPreferences) Scan(value interface{}) error {
    if value == nil {
        *up = make(UserPreferences)
        return nil
    }

    var bytes []byte
    switch v := value.(type) {
    case []byte:
        bytes = v
    case string:
        bytes = []byte(v)
    default:
        *up = make(UserPreferences)
        return nil
    }

    if len(bytes) == 0 {
        *up = make(UserPreferences)
        return nil
    }

    return json.Unmarshal(bytes, up)
}

func (up UserPreferences) Value() (driver.Value, error) {
    if len(up) == 0 {
        return "{}", nil
    }
    return json.Marshal(up)
}

const (
	PREF_SHOW_LOCATION  = "privacy.show_location"
	PREF_SHOW_LAST_SEEN = "privacy.show_last_seen"
	PREF_ALLOW_SEARCH   = "privacy.allow_search"

	PREF_LANGUAGE            = "app.language"
	PREF_THEME               = "app.theme"
	PREF_TIMEZONE            = "app.timezone"
	PREF_NOTIFICATIONS       = "notifications.enabled"
	PREF_EMAIL_NOTIFICATIONS = "notifications.email"
	PREF_PUSH_NOTIFICATIONS  = "notifications.push"
)

func DefaultUserPreferences() UserPreferences {
	return UserPreferences{
		PREF_SHOW_LOCATION:       true,
		PREF_SHOW_LAST_SEEN:      true,
		PREF_ALLOW_SEARCH:        true,
		PREF_LANGUAGE:            "tr",
		PREF_THEME:               "auto",
		PREF_TIMEZONE:            "Europe/Istanbul",
		PREF_NOTIFICATIONS:       true,
		PREF_EMAIL_NOTIFICATIONS: false,
		PREF_PUSH_NOTIFICATIONS:  true,
	}
}

func (up UserPreferences) GetBool(key string, defaultVal bool) bool {
	if val, exists := up[key]; exists {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return defaultVal
}

func (up UserPreferences) GetString(key string, defaultVal string) string {
	if val, exists := up[key]; exists {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return defaultVal
}

func (up UserPreferences) Set(key string, value interface{}) {
	up[key] = value
}

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
	FavoriteGenres pq.StringArray `json:"favorite_genres" gorm:"type:text[]" swaggertype:"array,string"`
	Bio            *string        `json:"bio" gorm:"type:text"`
	ReadingGoal    int            `json:"reading_goal" gorm:"default:0"`
	BooksRead      int            `json:"books_read" gorm:"default:0"`
	Badges         pq.StringArray `json:"badges" gorm:"type:text[]" swaggertype:"array,string"`
	IsOnline       bool           `json:"is_online" gorm:"default:false"`
	LastSeen       *time.Time     `json:"last_seen"`

	Preferences UserPreferences `json:"preferences" gorm:"type:jsonb;default:'{}'"`

	OwnedClubs      []Club           `json:"owned_clubs,omitempty" gorm:"foreignKey:OwnerID" swaggerignore:"true"`
	ClubMemberships []ClubMembership `json:"club_memberships,omitempty" gorm:"foreignKey:UserID" swaggerignore:"true"`
	Posts           []Post           `json:"posts,omitempty" gorm:"foreignKey:UserID" swaggerignore:"true"`
	Comments        []Comment        `json:"comments,omitempty" gorm:"foreignKey:UserID" swaggerignore:"true"`
	PostLikes       []PostLike       `json:"post_likes,omitempty" gorm:"foreignKey:UserID" swaggerignore:"true"`
	CommentLikes    []CommentLike    `json:"comment_likes,omitempty" gorm:"foreignKey:UserID" swaggerignore:"true"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type PublicUserProfile struct {
	ID             uint           `json:"id"`
	Username       string         `json:"username"`
	FirstName      string         `json:"first_name"`
	LastName       string         `json:"last_name"`
	AvatarURL      *string        `json:"avatar_url"`
	Location       *string        `json:"location"`
	FavoriteGenres pq.StringArray `json:"favorite_genres"`
	Bio            *string        `json:"bio"`
	BooksRead      int            `json:"books_read"`
	Badges         pq.StringArray `json:"badges"`
	IsOnline       bool           `json:"is_online"`
	LastSeen       *time.Time     `json:"last_seen,omitempty"`
	JoinedAt       time.Time      `json:"joined_at"`

	TotalPosts    int `json:"total_posts"`
	TotalComments int `json:"total_comments"`
	ClubsCount    int `json:"clubs_count"`
	ReadingStreak int `json:"reading_streak,omitempty"`
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
	FavoriteGenres pq.StringArray `json:"favorite_genres" swaggertype:"array,string"`
	Bio            *string        `json:"bio"`
	ReadingGoal    *int           `json:"reading_goal"`
	BooksRead      int            `json:"books_read"`
	Badges         pq.StringArray `json:"badges" swaggertype:"array,string"`
	IsOnline       bool           `json:"is_online"`
	LastSeen       *time.Time     `json:"last_seen"`

	Preferences UserPreferences `json:"preferences"`

	OwnedClubs      []Club           `json:"owned_clubs,omitempty" swaggerignore:"true"`
	ClubMemberships []ClubMembership `json:"club_memberships,omitempty" swaggerignore:"true"`
	Posts           []Post           `json:"posts,omitempty" swaggerignore:"true"`
	Comments        []Comment        `json:"comments,omitempty" swaggerignore:"true"`
	PostLikes       []PostLike       `json:"post_likes,omitempty" swaggerignore:"true"`
	CommentLikes    []CommentLike    `json:"comment_likes,omitempty" swaggerignore:"true"`

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

type UpdatePasswordRequest struct {
	Password    string `json:"password" validate:"required,min=6"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type UpdateProfileRequest struct {
	Bio            *string   `json:"bio,omitempty" validate:"omitempty"`
	Location       *string   `json:"location,omitempty" validate:"omitempty,max=255"`
	FavoriteGenres *[]string `json:"favorite_genres,omitempty"`
	ReadingGoal    *int      `json:"reading_goal,omitempty" validate:"omitempty,gte=0"`
}

type UpdateAvatarRequest struct {
	AvatarURL string `json:"avatar_url" validate:"required,url"`
}

type UpdateAccountRequest struct {
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=2,max=50"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=2,max=50"`
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
	Username  *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
}

type UpdatePreferencesRequest struct {
    Preferences UserPreferences `json:"preferences" validate:"required"`
}

func (u *User) ToResponse() UserResponse {
	if len(u.Preferences) == 0 {
		u.Preferences = DefaultUserPreferences()
	}

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
		Preferences:     u.Preferences,
		OwnedClubs:      u.OwnedClubs,
		ClubMemberships: u.ClubMemberships,
		Posts:           u.Posts,
		Comments:        u.Comments,
		PostLikes:       u.PostLikes,
		CommentLikes:    u.CommentLikes,
		CreatedAt:       u.CreatedAt,
	}
}

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a standard success response
type SuccessResponse struct {
	Message string `json:"message"`
}

func (u *User) ToPublicProfile() PublicUserProfile {
	prefs := u.Preferences
	if len(prefs) == 0 {
		prefs = DefaultUserPreferences()
	}

	profile := PublicUserProfile{
		ID:        u.ID,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		JoinedAt:  u.CreatedAt,
		BooksRead: u.BooksRead,
		Badges:    u.Badges,
		IsOnline:  u.IsOnline,
	}

	if prefs.GetBool(PREF_SHOW_LOCATION, true) {
		profile.Location = u.Location
	}

	if prefs.GetBool(PREF_SHOW_LAST_SEEN, true) && u.LastSeen != nil && u.IsOnline {
		since := time.Since(*u.LastSeen)
		if since <= 15*time.Minute {
			profile.LastSeen = u.LastSeen
		}
	}

	profile.AvatarURL = u.AvatarURL

	return profile
}

package models

import "time"

type Event struct {
	ID           uint        `json:"id" gorm:"primaryKey"`
	Title        string      `json:"title" gorm:"not null"`
	Description  string      `json:"description"`
	ClubID       uint        `json:"club_id" gorm:"not null"`
	EventType    EventType   `json:"event_type" gorm:"type:varchar(20)"`
	StartTime    time.Time   `json:"start_time" gorm:"not null"`
	EndTime      time.Time   `json:"end_time" gorm:"not null"`
	Location     string      `json:"location,omitempty"`
	OnlineLink   string      `json:"online_link,omitempty"`
	MaxAttendees *int        `json:"max_attendees,omitempty"`
	IsPublic	 bool        `json:"is_public" gorm:"default:false"`
	CreatedAt    time.Time   `json:"created_at" gorm:"autoCreateTime"`
	Club         Club        `json:"club" gorm:"foreignKey:ClubID"`
	RSVPs        []EventRSVP `json:"rsvps,omitempty"`
}

type EventType string

const (
	EventInPerson EventType = "in_person"
	EventOnline   EventType = "online"
)

type EventRSVP struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id" gorm:"not null"`
	EventID   uint       `json:"event_id" gorm:"not null"`
	Status    RSVPStatus `json:"status" gorm:"type:varchar(20);not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	User      User       `json:"user" gorm:"foreignKey:UserID"`
	Event     Event      `json:"event" gorm:"foreignKey:EventID"`
}

type RSVPStatus string

const (
	RSVPGoing    RSVPStatus = "going"
	RSVPMaybe    RSVPStatus = "maybe"
	RSVPNotGoing RSVPStatus = "not_going"
)

type CreateEventRequest struct {
	Title        string    `json:"title" binding:"required"`
	Description  string    `json:"description"`
	EventType    EventType `json:"event_type" binding:"required,oneof=in_person online"`
	StartTime    time.Time `json:"start_time" binding:"required"`
	EndTime      time.Time `json:"end_time" binding:"required,gtfield=StartTime"`
	Location     string    `json:"location,omitempty"`
	OnlineLink   string    `json:"online_link,omitempty"`
	MaxAttendees *int      `json:"max_attendees,omitempty"`
	IsPublic	 bool      `json:"is_public" gorm:"default:false"`
}

type UpdateEventRequest struct {
	Title        *string    `json:"title,omitempty"`
	Description  *string    `json:"description,omitempty"`
	EventType    *EventType `json:"event_type,omitempty" binding:"omitempty,oneof=in_person online"`
	StartTime    *time.Time `json:"start_time,omitempty" binding:"omitempty"`
	EndTime      *time.Time `json:"end_time,omitempty" binding:"omitempty,gtfield=StartTime"`
	Location     *string    `json:"location,omitempty"`
	OnlineLink   *string    `json:"online_link,omitempty"`
	MaxAttendees *int       `json:"max_attendees,omitempty"`
	IsPublic	 *bool      `json:"is_public,omitempty"`
}

type RSVPRequest struct {
	Status RSVPStatus `json:"status" binding:"required,oneof=going maybe not_going"`
}

type EventResponse struct {
	ID           uint        `json:"id"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	ClubID       uint        `json:"club_id"`
	EventType    EventType   `json:"event_type"`
	StartTime    string      `json:"start_time"`
	EndTime      string      `json:"end_time"`
	Location     string      `json:"location,omitempty"`
	OnlineLink   string      `json:"online_link,omitempty"`
	MaxAttendees *int        `json:"max_attendees,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	RSVPs        []EventRSVP `json:"rsvps,omitempty"`
	IsPublic     *bool       `json:"is_public,omitempty"`
}

func (e *Event) ToResponse() EventResponse {
	return EventResponse{
		ID:           e.ID,
		Title:        e.Title,
		Description:  e.Description,
		ClubID:       e.ClubID,
		EventType:    e.EventType,
		StartTime:    e.StartTime.Format(time.RFC3339),
		EndTime:      e.EndTime.Format(time.RFC3339),
		Location:     e.Location,
		OnlineLink:   e.OnlineLink,
		MaxAttendees: e.MaxAttendees,
		CreatedAt:    e.CreatedAt,
		RSVPs:        e.RSVPs,
		IsPublic:     &e.IsPublic,
	}
}

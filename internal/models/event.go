package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type DateYMD struct{ time.Time } // accepts "2006-01-02"
type TimeHM struct{ time.Time }  // accepts "15:04"
type DBTime struct{ time.Time }

func (t *DBTime) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		t.Time = v
		return nil
	case string:
		parsed, err := time.Parse("15:04:05", v)
		if err != nil {
			return err
		}
		t.Time = parsed
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into DBTime", value)
	}
}

func (t DBTime) Value() (driver.Value, error) {
	return t.Time.Format("15:04:05"), nil
}

func (d *DateYMD) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("date must be a string YYYY-MM-DD: %w", err)
	}
	s = strings.TrimSpace(s)
	tt, err := time.ParseInLocation("2006-01-02", s, time.Local) // or a fixed TZ
	if err != nil {
		return fmt.Errorf("date must be YYYY-MM-DD: %w", err)
	}
	d.Time = tt
	return nil
}

func (t *TimeHM) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("time must be a string HH:MM (24h): %w", err)
	}
	s = strings.TrimSpace(s)
	tt, err := time.ParseInLocation("15:04", s, time.Local) // or a fixed TZ
	if err != nil {
		return fmt.Errorf("time must be HH:MM (24h): %w", err)
	}
	t.Time = tt
	return nil
}

type Event struct {
	ID           uint        `json:"id" gorm:"primaryKey"`
	Title        string      `json:"title" gorm:"not null"`
	Description  string      `json:"description"`
	ClubID       uint        `json:"club_id" gorm:"not null"`
	EventType    EventType   `json:"event_type" gorm:"type:varchar(20)"`
	EventDate    time.Time   `json:"event_date" gorm:"type:date;not null"`
	EventTime    DBTime      `json:"event_time" gorm:"type:time;not null"`
	Location     string      `json:"location,omitempty"`
	OnlineLink   string      `json:"online_link,omitempty"`
	MaxAttendees *int        `json:"max_attendees,omitempty"`
	IsPublic     bool        `json:"is_public" gorm:"default:false"`
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
	EventDate    DateYMD   `json:"event_date" binding:"required"`
	EventTime    TimeHM    `json:"event_time" binding:"required"`
	Location     string    `json:"location,omitempty"`
	OnlineLink   string    `json:"online_link,omitempty"`
	MaxAttendees *int      `json:"max_attendees,omitempty"`
	IsPublic     bool      `json:"is_public" gorm:"default:false"`
}

type UpdateEventRequest struct {
	Title        *string    `json:"title,omitempty"`
	Description  *string    `json:"description,omitempty"`
	EventType    *EventType `json:"event_type,omitempty" binding:"omitempty,oneof=in_person online"`
	EventDate    *DateYMD   `json:"event_date,omitempty" binding:"omitempty"`
	EventTime    *TimeHM    `json:"event_time,omitempty" binding:"omitempty"`
	Location     *string    `json:"location,omitempty"`
	OnlineLink   *string    `json:"online_link,omitempty"`
	MaxAttendees *int       `json:"max_attendees,omitempty"`
	IsPublic     *bool      `json:"is_public,omitempty"`
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
	EventDate    string      `json:"event_date"`
	EventTime    string      `json:"event_time"`
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
		EventDate:    e.EventDate.Format("2006-01-02"),
		EventTime:    e.EventTime.Format("15:04"),
		Location:     e.Location,
		OnlineLink:   e.OnlineLink,
		MaxAttendees: e.MaxAttendees,
		CreatedAt:    e.CreatedAt,
		RSVPs:        e.RSVPs,
		IsPublic:     &e.IsPublic,
	}
}

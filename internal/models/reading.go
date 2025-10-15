package models

import (
    "time"
)

type ReadingStatus string
const (
    ReadingNotStarted ReadingStatus = "not_started"
    ReadingActive     ReadingStatus = "reading"
    ReadingPaused     ReadingStatus = "paused"
    ReadingFinished   ReadingStatus = "finished"
)

type UserBookProgress struct {
    ID          uint          `json:"id" gorm:"primaryKey"`
    UserID      uint          `json:"user_id" gorm:"index;not null"`
    BookID      uint          `json:"book_id" gorm:"index;not null"`
    Status      ReadingStatus `json:"status" gorm:"type:text;default:'reading'"`
    CurrentPage *int          `json:"current_page,omitempty"`
    Percent     *float32      `json:"percent,omitempty" gorm:"type:numeric(5,2)"`
    StartedAt   time.Time     `json:"started_at" gorm:"autoCreateTime"`
    FinishedAt  *time.Time    `json:"finished_at,omitempty"`
    UpdatedAt   time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}

func (p *UserBookProgress) ToResponse(book *Book) UserBookProgressResponse {
    var percent float32
    if p.Percent != nil {
        percent = *p.Percent
    } else if book.Pages != nil && p.CurrentPage != nil && *book.Pages > 0 {
        percent = float32(*p.CurrentPage) * 100 / float32(*book.Pages)
    }
    return UserBookProgressResponse{
        UserID:      p.UserID,
        BookID:      p.BookID,
        Status:      string(p.Status),
        CurrentPage: p.CurrentPage,
        Percent:     &percent,
        StartedAt:   p.StartedAt,
        FinishedAt:  p.FinishedAt,
        UpdatedAt:   p.UpdatedAt,
    }
}

type ReadingLog struct {
    ID           uint       `json:"id" gorm:"primaryKey"`
    UserID       uint       `json:"user_id" gorm:"index;not null"`
    BookID       uint       `json:"book_id" gorm:"index;not null"`
    ClubID       *uint      `json:"club_id,omitempty" gorm:"index"`
    AssignmentID *uint      `json:"assignment_id,omitempty"`
    PagesDelta   *int       `json:"pages_delta,omitempty"`
    FromPage     *int       `json:"from_page,omitempty"`
    ToPage       *int       `json:"to_page,omitempty"`
    Minutes      *int       `json:"minutes,omitempty"`
    Note         *string    `json:"note,omitempty" gorm:"type:text"`
    CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

type ClubBookAssignmentStatus string
const (
    ClubAssignmentActive    ClubBookAssignmentStatus = "active"
    ClubAssignmentCompleted ClubBookAssignmentStatus = "completed"
    ClubAssignmentArchived  ClubBookAssignmentStatus = "archived"
)

type ClubBookAssignment struct {
    ID          uint                        `json:"id" gorm:"primaryKey"`
    ClubID      uint                        `json:"club_id" gorm:"index;not null"`
    BookID      uint                        `json:"book_id" gorm:"index;not null"`
    Status      ClubBookAssignmentStatus    `json:"status" gorm:"type:text;default:'active'"`
    StartDate   *time.Time                  `json:"start_date,omitempty"`
    DueDate     *time.Time                  `json:"due_date,omitempty"`
    CompletedAt *time.Time                  `json:"completed_at,omitempty"`
    CreatedAt   time.Time                   `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time                   `json:"updated_at" gorm:"autoUpdateTime"`
}

type StartReadingRequest struct {
    BookID uint `json:"book_id" validate:"required"`
}

type UpdateReadingProgressRequest struct {
    CurrentPage *int     `json:"current_page" validate:"omitempty,gte=0"`
    Percent     *float32 `json:"percent" validate:"omitempty,gte=0,lte=100"`
    PagesDelta  *int     `json:"pages_delta" validate:"omitempty"`
    Minutes     *int     `json:"minutes" validate:"omitempty,gte=0"`
    Note        *string  `json:"note" validate:"omitempty,max=500"`
}

type CompleteReadingRequest struct {
    Note *string `json:"note,omitempty" validate:"omitempty,max=500"`
}

type AssignBookRequest struct {
    BookID     uint       `json:"book_id" validate:"required"`
    StartDate  *time.Time `json:"start_date,omitempty"`
    DueDate    *time.Time `json:"due_date,omitempty"`
}

type UserBookProgressResponse struct {
    UserID      uint       `json:"user_id"`
    BookID      uint       `json:"book_id"`
    Status      string     `json:"status"`
    CurrentPage *int       `json:"current_page,omitempty"`
    Percent     *float32   `json:"percent,omitempty"`
    StartedAt   time.Time  `json:"started_at"`
    FinishedAt  *time.Time `json:"finished_at,omitempty"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

type UserReadingHistoryItem struct {
    Book       Book       `json:"book"`
    FinishedAt time.Time  `json:"finished_at"`
    Logs       []ReadingLog `json:"logs,omitempty"`
}

type ClubAssignmentResponse struct {
    ID         uint       `json:"id"`
    ClubID     uint       `json:"club_id"`
    Book       Book       `json:"book"`
    Status     string     `json:"status"`
    StartDate  *time.Time `json:"start_date,omitempty"`
    DueDate    *time.Time `json:"due_date,omitempty"`
}

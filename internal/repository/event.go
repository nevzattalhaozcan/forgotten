package repository

import (
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"gorm.io/gorm"
)

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *eventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(event *models.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) GetClubEvents(clubID uint) ([]models.Event, error) {
	var events []models.Event
	err := r.db.
		Where("club_id = ?", clubID).
		Preload("RSVPs").
		Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *eventRepository) GetByID(id uint) (*models.Event, error) {
	var event models.Event
	err := r.db.
		Preload("Club").
		Preload("RSVPs").
		Preload("RSVPs.User").
		First(&event, id).Error
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *eventRepository) Update(event *models.Event) error {
	return r.db.Save(event).Error
}

func (r *eventRepository) Delete(id uint) error {
	return r.db.Delete(&models.Event{}, id).Error
}

func (r *eventRepository) RSVP(eventID uint, rsvp *models.EventRSVP) error {
	var existingRSVP models.EventRSVP
	err := r.db.
		Where("user_id = ? AND event_id = ?", rsvp.UserID, eventID).
		First(&existingRSVP).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return r.db.Create(rsvp).Error
		}
		return err
	}
	existingRSVP.Status = rsvp.Status
	return r.db.Save(&existingRSVP).Error
}

func (r *eventRepository) GetEventAttendees(eventID uint) ([]models.EventRSVP, error) {
	var rsvps []models.EventRSVP
	err := r.db.
		Where("event_id = ?", eventID).
		Preload("User").
		Find(&rsvps).Error
	if err != nil {
		return nil, err
	}
	return rsvps, nil
}

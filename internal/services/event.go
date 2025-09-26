package services

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/nevzattalhaozcan/forgotten/internal/config"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/internal/repository"
	"gorm.io/gorm"
)

type EventService struct {
	eventRepo repository.EventRepository
	clubRepo  repository.ClubRepository
	config    *config.Config
}

func NewEventService(eventRepo repository.EventRepository, clubRepo repository.ClubRepository, config *config.Config) *EventService {
	return &EventService{
		eventRepo: eventRepo,
		clubRepo:  clubRepo,
		config:    config,
	}
}

func (s *EventService) CreateEvent(clubID uint, req *models.CreateEventRequest) (*models.EventResponse, error) {
	_, err := s.clubRepo.GetByID(clubID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("club not found")
		}
		return nil, err
	}

	event := &models.Event{
		Title:        req.Title,
		Description:  req.Description,
		ClubID:       clubID,
		EventType:    req.EventType,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		Location:     req.Location,
		OnlineLink:   req.OnlineLink,
		MaxAttendees: req.MaxAttendees,
	}

	if err := s.eventRepo.Create(event); err != nil {
		return nil, err
	}

	response := event.ToResponse()
	_ = s.refreshClubNextMeeting(clubID)
	return &response, nil
}

func (s *EventService) GetClubEvents(clubID uint) ([]models.EventResponse, error) {
	_, err := s.clubRepo.GetByID(clubID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("club not found")
		}
		return nil, err
	}

	events, err := s.eventRepo.GetClubEvents(clubID)
	if err != nil {
		return nil, err
	}

	var responses []models.EventResponse
	for _, event := range events {
		response := event.ToResponse()
		responses = append(responses, response)
	}

	return responses, nil
}

func (s *EventService) GetEventByID(id uint) (*models.EventResponse, error) {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("event not found")
		}
		return nil, err
	}
	response := event.ToResponse()
	return &response, nil
}

func (s *EventService) UpdateEvent(id uint, req *models.UpdateEventRequest) (*models.EventResponse, error) {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("event not found")
		}
		return nil, err
	}

	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.EventType != nil {
		event.EventType = *req.EventType
	}
	if req.StartTime != nil {
		event.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		event.EndTime = *req.EndTime
	}
	if req.Location != nil {
		event.Location = *req.Location
	}
	if req.OnlineLink != nil {
		event.OnlineLink = *req.OnlineLink
	}
	if req.MaxAttendees != nil {
		event.MaxAttendees = req.MaxAttendees
	}

	if err := s.eventRepo.Update(event); err != nil {
		return nil, err
	}

	response := event.ToResponse()
	_ = s.refreshClubNextMeeting(event.ClubID)
	return &response, nil
}

func (s *EventService) DeleteEvent(id uint) error {
	event, err := s.eventRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("event not found")
		}
		return err
	}
	if err := s.eventRepo.Delete(id); err != nil {
		return err
	}

	_ = s.refreshClubNextMeeting(event.ClubID)

	return nil
}

func (s *EventService) RSVPToEvent(id uint, rsvp *models.EventRSVP) error {
	_, err := s.eventRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("event not found")
		}
		return err
	}
	return s.eventRepo.RSVP(id, rsvp)
}

func (s *EventService) GetEventAttendees(eventID uint) ([]models.EventRSVP, error) {
	_, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("event not found")
		}
		return nil, err
	}
	return s.eventRepo.GetEventAttendees(eventID)
}

func (s *EventService) refreshClubNextMeeting(clubID uint) error {
	events, err := s.eventRepo.GetClubEvents(clubID)
	if err != nil {
		return err
	}

	now := time.Now()
	var next *models.Event
	for i := range events {
		e := events[i]
		if e.StartTime.Before(now) {
			continue
		}
		if next == nil || e.StartTime.Before(next.StartTime) {
			next = &e
		}
	}

	club, err := s.clubRepo.GetByID(clubID)
	if err != nil {
		return err
	}

	if next == nil {
		club.NextMeeting = nil
		return s.clubRepo.Update(club)
	}

	var loc *string
	switch next.EventType {
	case models.EventOnline:
		if next.OnlineLink != "" {
			l := next.OnlineLink
			loc = &l
		}
	default:
		if next.Location != "" {
			l := next.Location
			loc = &l
		}	
	}

	topic := next.Title
	nm := models.NextMeeting{
		Date: &next.StartTime,
		Location: loc,
		Topic: &topic,
	}
	if b, merr := json.Marshal(&nm); merr == nil {
		club.NextMeeting = b
	}

	return s.clubRepo.Update(club)
}

func (s *EventService) ClubRepo() repository.ClubRepository {
    return s.clubRepo
}
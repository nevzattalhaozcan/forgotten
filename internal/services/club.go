package services

import (
	"encoding/json"
	"errors"

	"github.com/nevzattalhaozcan/forgotten/internal/config"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/internal/repository"
	"gorm.io/gorm"
)

type ClubService struct {
	clubRepo repository.ClubRepository
	config   *config.Config
}

func NewClubService(clubRepo repository.ClubRepository, config *config.Config) *ClubService {
	return &ClubService{
		clubRepo: clubRepo,
		config:   config,
	}
}

func (s *ClubService) CreateClub(req *models.CreateClubRequest) (*models.ClubResponse, error) {
	_, err := s.clubRepo.GetByName(req.Name)
	if err == nil {
		return nil, errors.New("club name already exists")
	}

	club := &models.Club{
		Name:        req.Name,
		Description: req.Description,
		Location:   req.Location,
		Genre:      req.Genre,
		CoverImageURL: req.CoverImageURL,
		IsPrivate:   req.IsPrivate,
		MaxMembers: req.MaxMembers,
		Tags:       req.Tags,
	}

	if err := s.clubRepo.Create(club); err != nil {
		return nil, err
	}

	response := club.ToResponse()
	return &response, nil
}

func (s *ClubService) GetClubByID(id uint) (*models.ClubResponse, error) {
	club, err := s.clubRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("club not found")
		}
		return nil, err
	}
	
	response := club.ToResponse()
	return &response, nil
}

func (s *ClubService) GetAllClubs(limit, offset int) ([]*models.ClubResponse, error) {
	clubs, err := s.clubRepo.List(limit, offset)
	if err != nil {
		return nil, err
	}

	var responses []*models.ClubResponse
	for _, club := range clubs {
		resp := club.ToResponse()
		responses = append(responses, &resp)
	}

	return responses, nil
}

func (s *ClubService) UpdateClub(id uint, req *models.UpdateClubRequest) (*models.ClubResponse, error) {
	club, err := s.clubRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("club not found")
		}
		return nil, err
	}

	if req.Name != nil && *req.Name != club.Name {
		existingClub, err := s.clubRepo.GetByName(*req.Name)
		if err == nil && existingClub.ID != club.ID {
			return nil, errors.New("club name already exists")
		}
		club.Name = *req.Name
	}

	if req.Description != nil {
		club.Description = *req.Description
	}
	if req.Location != nil {
		club.Location = req.Location
	}
	if req.Genre != nil {
		club.Genre = req.Genre
	}
	if req.CoverImageURL != nil {
		club.CoverImageURL = req.CoverImageURL
	}
	if req.IsPrivate != nil {
		club.IsPrivate = *req.IsPrivate
	}
	if req.MaxMembers != nil {
		club.MaxMembers = *req.MaxMembers
	}
	if req.Tags != nil {
		club.Tags = *req.Tags
	}
	if req.CurrentBook != nil {
		cb, err := json.Marshal(req.CurrentBook)
		if err == nil {
			club.CurrentBook = cb
		}
	}
	if req.NextMeeting != nil {
		nm, err := json.Marshal(req.NextMeeting)
		if err == nil {
			club.NextMeeting = nm
		}
	}

	if err := s.clubRepo.Update(club); err != nil {
		return nil, err
	}

	response := club.ToResponse()
	return &response, nil
}

func (s *ClubService) DeleteClub(id uint) error {
	club, err := s.clubRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("club not found")
		}
		return err
	}

	return s.clubRepo.Delete(club.ID)
}

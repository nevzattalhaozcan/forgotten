package services

import (
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"github.com/nevzattalhaozcan/forgotten/internal/config"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrClubNotFound   = errors.New("club not found")
	ErrMemberNotFound = errors.New("member not found")
	ErrClubNameExists = errors.New("club name already exists")
)

type ClubService struct {
	clubRepo       repository.ClubRepository
	clubRatingRepo repository.ClubRatingRepository
	config         *config.Config
}

func NewClubService(clubRepo repository.ClubRepository, clubRatingRepo repository.ClubRatingRepository, config *config.Config) *ClubService {
	return &ClubService{
		clubRepo: clubRepo,
		clubRatingRepo: clubRatingRepo,
		config:   config,
	}
}

func (s *ClubService) CanManageClub(clubID, userID uint) bool {
	club, err := s.clubRepo.GetByID(clubID)
	if err != nil {
		return false
	}

	if club.OwnerID != nil && *club.OwnerID == userID {
		return true
	}

	member, err := s.clubRepo.GetClubMemberByUserID(clubID, userID)
	if err == nil {
		return false
	}
	if !member.IsApproved {
		return false
	}
	return member.Role == "moderator" || member.Role == "admin"
}

func (s *ClubService) CreateClub(ownerID uint, req *models.CreateClubRequest) (*models.ClubResponse, error) {
	_, err := s.clubRepo.GetByName(req.Name)
	if err == nil {
		return nil, ErrClubNameExists
	}

	club := &models.Club{
		Name:          req.Name,
		Description:   req.Description,
		Location:      req.Location,
		Genre:         req.Genre,
		CoverImageURL: req.CoverImageURL,
		IsPrivate:     req.IsPrivate,
		MaxMembers:    req.MaxMembers,
		Tags:          req.Tags,
		OwnerID:       &ownerID,
	}

	if err := s.clubRepo.Create(club); err != nil {
		if isUniqueViolation(err) {
			return nil, ErrClubNameExists
		}
		return nil, err
	}

	ownerMembership := &models.ClubMembership{
		ClubID:     club.ID,
		UserID:     ownerID,
		Role:       "admin",
		IsApproved: true,
	}
	if err := s.clubRepo.JoinClub(ownerMembership); err != nil {
		_ = s.clubRepo.Delete(club.ID)
		return nil, err
	}
	club.MembersCount = 1
	if err := s.clubRepo.Update(club); err != nil {
		return nil, err
	}

	created, err := s.clubRepo.GetByID(club.ID)
    if err != nil {
        return nil, err
    }
    response := created.ToResponse()
    return &response, nil
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && string(pqErr.Code) == "23505" {
		return true
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}
	return false
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

func (s *ClubService) GetAllClubs() ([]*models.ClubResponse, error) {
	clubs, err := s.clubRepo.List(50, 0)
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

	created, err := s.clubRepo.GetByID(club.ID)
    if err != nil {
        return nil, err
    }
    response := created.ToResponse()
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

func (s *ClubService) JoinClub(clubID, userID uint) (*models.ClubMembershipResponse, error) {
	club, err := s.clubRepo.GetByID(clubID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("club not found")
		}
		return nil, err
	}

	if _, err := s.clubRepo.GetClubMemberByUserID(clubID, userID); err == nil {
		return nil, errors.New("user is already a member of the club")
	}

	//TODO: Handle invitations for private clubs or approval process
	if club.IsPrivate {
		return nil, errors.New("cannot join a private club without an invitation")
	}

	approve := !club.IsPrivate
	if approve {
		if club.MaxMembers > 0 && club.MembersCount >= club.MaxMembers {
			return nil, errors.New("club is full")
		}
	}

	membership := &models.ClubMembership{
		ClubID:     clubID,
		UserID:     userID,
		Role:       "member",
		IsApproved: approve,
	}
	if err := s.clubRepo.JoinClub(membership); err != nil {
		return nil, err
	}

	if approve {
		club.MembersCount++
		if err := s.clubRepo.Update(club); err != nil {
			return nil, err
		}
	}

	withUser, err := s.clubRepo.GetClubMemberByUserID(clubID, userID)
	if err != nil {
		return nil, err
	}
	response := withUser.ToResponse()
	return &response, nil
}

func (s *ClubService) LeaveClub(clubID, userID uint) error {
	club, err := s.clubRepo.GetByID(clubID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("club not found")
		}
		return err
	}
	m, err := s.clubRepo.GetClubMemberByUserID(clubID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("not a member")
		}
		return err
	}
	if m.IsApproved && club.MembersCount > 0 {
		club.MembersCount--
		if err := s.clubRepo.Update(club); err != nil {
			return err
		}
	}
	return s.clubRepo.LeaveClub(clubID, userID)
}

func (s *ClubService) ListClubMembers(clubID uint) ([]*models.ClubMembership, error) {
	_, err := s.clubRepo.GetByID(clubID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("club not found")
		}
		return nil, err
	}
	return s.clubRepo.ListClubMembers(clubID)
}

func (s *ClubService) UpdateClubMember(membership *models.ClubMembership) error {
	return s.clubRepo.UpdateClubMember(membership)
}

func (s *ClubService) GetClubMemberByUserID(clubID, userID uint) (*models.ClubMembership, error) {
	return s.clubRepo.GetClubMemberByUserID(clubID, userID)
}

func (s *ClubService) UpdateClubMemberFields(clubID, userID uint, req *models.UpdateClubMembershipRequest) (*models.ClubMembershipResponse, error) {
	club, err := s.clubRepo.GetByID(clubID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClubNotFound
		}
		return nil, err
	}
	member, err := s.clubRepo.GetClubMemberByUserID(clubID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}

	wasApproved := member.IsApproved
	newApproved := member.IsApproved

	if req.Role != nil {
		member.Role = *req.Role
	}
	if req.IsApproved != nil {
		newApproved = *req.IsApproved
		member.IsApproved = *req.IsApproved
	}

	if !wasApproved && newApproved {
		if club.MaxMembers > 0 && club.MembersCount >= club.MaxMembers {
			return nil, errors.New("club is full, cannot approve")
		}
	}

	if err := s.clubRepo.UpdateClubMember(&models.ClubMembership{
		ClubID:     clubID,
		UserID:     userID,
		Role:       member.Role,
		IsApproved: member.IsApproved,
	}); err != nil {
		return nil, err
	}

	if !wasApproved && newApproved {
		club.MembersCount++
		if err := s.clubRepo.Update(club); err != nil {
			return nil, err
		}
	} else if wasApproved && !newApproved && club.MembersCount > 0 {
		club.MembersCount--
		if err := s.clubRepo.Update(club); err != nil {
			return nil, err
		}
	}

	updated, _ := s.clubRepo.GetClubMemberByUserID(clubID, userID)
	response := updated.ToResponse()
	return &response, nil
}

func (s *ClubService) RateClub(userID, clubID uint, req *models.RateClubRequest) (*models.ClubResponse, error) {
	if req.Rating < 1 || req.Rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}

	_, err := s.clubRepo.GetByID(clubID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClubNotFound
		}
		return nil, err
	}

	cr := &models.ClubRating{
		ClubID: clubID,
		UserID: userID,
		Rating: req.Rating,
		Comment: req.Comment,
	}
	if err := s.clubRatingRepo.UpsertRating(cr); err != nil {
		return nil, err
	}

	avg, count, err := s.clubRatingRepo.GetAggregateForClub(clubID)
	if err != nil { return nil, err }
	if err := s.clubRepo.UpdateRatingAggregate(clubID, avg, count); err != nil {
		return nil, err
	}

	club, err := s.clubRepo.GetByID(clubID)
	if err != nil {
		return nil, err
	}

	response := club.ToResponse()
	return &response, nil
}

func (s *ClubService) ListClubRatings(clubID uint, limit, offset int) ([]models.ClubRating, error) {
	_, err := s.clubRepo.GetByID(clubID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClubNotFound
		}
		return nil, err
	}
	return s.clubRatingRepo.ListByClub(clubID, limit, offset)
}
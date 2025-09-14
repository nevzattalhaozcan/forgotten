package services

import (
	"errors"

	"github.com/nevzattalhaozcan/forgotten/internal/config"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/internal/repository"
	"github.com/nevzattalhaozcan/forgotten/pkg/metrics"
	"github.com/nevzattalhaozcan/forgotten/pkg/utils"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo repository.UserRepository
	config  *config.Config
}

func NewUserService(userRepo repository.UserRepository, config *config.Config) *UserService {
	return &UserService{
		userRepo: userRepo,
		config:   config,
	}
}

func (s *UserService) Register(req *models.RegisterRequest) (*models.UserResponse, error) {
	_, err := s.userRepo.GetByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	_, err = s.userRepo.GetByUsername(req.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	if req.Role == "" {
		req.Role = "user"
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		PasswordHash: hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive: true,
		Role: req.Role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// increment user count metric
	metrics.IncrementUserCount(1)

	response := user.ToResponse()
	return &response, nil
}

func (s *UserService) Login(req *models.LoginRequest) (string, *models.UserResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			metrics.RecordAuthAttempt(false)
			return "", nil, errors.New("invalid email or password")
		}
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		metrics.RecordAuthAttempt(false)
		return "", nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateJWT(
		user.ID, 
		user.Role, 
		s.config.JWT.Secret, 
		s.config.JWT.ExpirationHours,
	)
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	metrics.RecordAuthAttempt(true)
	response := user.ToResponse()
	return token, &response, nil
}

func (s *UserService) GetUserByID(id uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *UserService) GetAllUsers() ([]models.UserResponse, error) {
	users, err := s.userRepo.List(50, 0)
	if err != nil {
		return nil, err
	}

	var responses []models.UserResponse
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	return responses, nil
}
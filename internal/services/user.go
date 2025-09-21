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
		AvatarURL: &req.AvatarURL,
		Location: &req.Location,
		FavoriteGenres: req.FavoriteGenres,
		Bio: &req.Bio,
		ReadingGoal: req.ReadingGoal,
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
		user.Email, 
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

func (s *UserService) UpdateUser(id uint, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if req.Username != nil && *req.Username != user.Username {
		existingUser, err := s.userRepo.GetByUsername(*req.Username)
		if err == nil && existingUser.ID != user.ID {
			return nil, errors.New("username already exists")
		}
		user.Username = *req.Username
	}

	if req.Email != nil && *req.Email != user.Email {
		existingUser, err := s.userRepo.GetByEmail(*req.Email)
		if err == nil && existingUser.ID != user.ID {
			return nil, errors.New("email already exists")
		}
		user.Email = *req.Email
	}

	if req.Password != nil {
		hashedPassword, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		user.PasswordHash = hashedPassword
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		user.LastName = *req.LastName
	}

	if req.Role != nil {
		user.Role = *req.Role
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	if req.Location != nil {
		user.Location = req.Location
	}

	if req.FavoriteGenres != nil {
		user.FavoriteGenres = *req.FavoriteGenres
	}

	if req.Bio != nil {
		user.Bio = req.Bio
	}

	if req.ReadingGoal != nil {
		user.ReadingGoal = *req.ReadingGoal
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update user")
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *UserService) DeleteUser(id uint) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	user.IsActive = false
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("failed to deactivate user")
	}
	
	err = s.userRepo.Delete(id)
	if err != nil {
		return err
	}

	// decrement user count metric
	metrics.IncrementUserCount(-1)

	return nil
}
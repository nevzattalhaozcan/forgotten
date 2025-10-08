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
	config   *config.Config
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
		Username:       req.Username,
		Email:          req.Email,
		PasswordHash:   hashedPassword,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		IsActive:       true,
		Role:           req.Role,
		AvatarURL:      &req.AvatarURL,
		Location:       &req.Location,
		FavoriteGenres: req.FavoriteGenres,
		Bio:            &req.Bio,
		ReadingGoal:    req.ReadingGoal,
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

func (s *UserService) GetPublicUserProfile(userID uint, viewerID *uint) (*models.PublicUserProfile, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("user profile is not available")
	}

	profile := user.ToPublicProfile()

	profile.TotalPosts = len(user.Posts)
	profile.TotalComments = len(user.Comments)
	profile.ClubsCount = len(user.ClubMemberships)

	return &profile, nil
}

func (s *UserService) GetUserProfileStats(userID uint) (map[string]interface{}, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	stats := map[string]interface{}{
		"total_posts":    len(user.Posts),
		"total_comments": len(user.Comments),
		"clubs_joined":   len(user.ClubMemberships),
		"books_read":     user.BooksRead,
		"profile_views":  0,
		"member_since":   user.CreatedAt.Format("January 2006"),
	}

	return stats, nil
}

func (s *UserService) SearchPublicUsers(query string, limit int) ([]models.PublicUserProfile, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	users, err := s.userRepo.SearchByUsernameOrName(query, limit)
	if err != nil {
		return nil, err
	}

	var profiles []models.PublicUserProfile
	for _, user := range users {
		if user.IsActive {
			profiles = append(profiles, user.ToPublicProfile())
		}
	}

	return profiles, nil
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

func (s *UserService) UpdatePassword(id uint, req *models.UpdatePasswordRequest) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return errors.New("current password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	user.PasswordHash = hashedPassword
	return s.userRepo.Update(user)
}

func (s *UserService) UpdateProfile(id uint, req *models.UpdateProfileRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if req.Bio != nil {
		user.Bio = req.Bio
	}

	if req.Location != nil {
		user.Location = req.Location
	}

	if req.FavoriteGenres != nil {
		user.FavoriteGenres = *req.FavoriteGenres
	}

	if req.ReadingGoal != nil {
		user.ReadingGoal = *req.ReadingGoal
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update profile")
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *UserService) UpdateAvatar(id uint, req *models.UpdateAvatarRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.AvatarURL = &req.AvatarURL
	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update avatar")
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *UserService) UpdateAccount(id uint, req *models.UpdateAccountRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		user.LastName = *req.LastName
	}

	if req.Email != nil && *req.Email != user.Email {
		existingUser, err := s.userRepo.GetByEmail(*req.Email)
		if err == nil && existingUser.ID != user.ID {
			return nil, errors.New("email already exists")
		}
		user.Email = *req.Email
	}

	if req.Username != nil && *req.Username != user.Username {
		existingUser, err := s.userRepo.GetByUsername(*req.Username)
		if err == nil && existingUser.ID != user.ID {
			return nil, errors.New("username already exists")
		}
		user.Username = *req.Username
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update account")
	}

	response := user.ToResponse()
	return &response, nil
}

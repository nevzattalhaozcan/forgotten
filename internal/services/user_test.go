package services

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nevzattalhaozcan/forgotten/internal/configtest"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/internal/repository/mocks"
	"github.com/nevzattalhaozcan/forgotten/pkg/testutil"
	"github.com/nevzattalhaozcan/forgotten/pkg/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	cfg := configtest.New()
	userService := NewUserService(mockUserRepo, cfg)

	t.Run("successful registration", func(t *testing.T) {
		req, err := test_helpers.CreateRegisterRequest()
		assert.NoError(t, err)

		mockUserRepo.EXPECT().
			GetByEmail(req.Email).
			Return(nil, gorm.ErrRecordNotFound)

		mockUserRepo.EXPECT().
			GetByUsername(req.Username).
			Return(nil, gorm.ErrRecordNotFound)

		mockUserRepo.EXPECT().
			Create(gomock.Any()).
			DoAndReturn(func(user *models.User) error {
				user.ID = 1
				return nil
			})

		result, err := userService.Register(req)
		
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Username, result.Username)
		assert.Equal(t, req.Email, result.Email)
	})

	t.Run("email already exists", func(t *testing.T) {
		req, err := test_helpers.CreateRegisterRequest()
		assert.NoError(t, err)

		existingUser := &models.User{
			ID:       1,
			Email:    "john@example.com",
		}

		mockUserRepo.EXPECT().
			GetByEmail(req.Email).
			Return(existingUser, nil)

		result, err := userService.Register(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestUserService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	cfg := configtest.New()
	userService := NewUserService(mockUserRepo, cfg)

	t.Run("successful login", func(t *testing.T) {
		req, err := test_helpers.CreateLoginRequest()
		assert.NoError(t, err)

		user, err := test_helpers.CreateUser()
		assert.NoError(t, err)

		user.PasswordHash, err = utils.HashPassword(user.PasswordHash)
		assert.NoError(t, err)

		mockUserRepo.EXPECT().
			GetByEmail(req.Email).
			Return(user, nil)

		token, userResp, err := userService.Login(req)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotNil(t, userResp)
		assert.Equal(t, user.Email, userResp.Email)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		req, err := test_helpers.CreateLoginRequest()
		assert.NoError(t, err)

		mockUserRepo.EXPECT().
			GetByEmail(req.Email).
			Return(nil, gorm.ErrRecordNotFound)

		token, userResp, err := userService.Login(req)
		
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, userResp)
		assert.Contains(t, err.Error(), "invalid email or password")
	})
}


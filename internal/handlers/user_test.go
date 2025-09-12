package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nevzattalhaozcan/forgotten/internal/configtest"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/internal/repository"
	"github.com/nevzattalhaozcan/forgotten/internal/services"
	"github.com/nevzattalhaozcan/forgotten/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserHandlerTestSuite struct {
	suite.Suite
	db 		*gorm.DB
	router *gin.Engine
	userHandler *UserHandler
}

func (suite *UserHandlerTestSuite) SetupTest() {
	var err error

	suite.db, err = test_helpers.SetupTestDB()
	suite.Require().NoError(err)

	gin.SetMode(gin.TestMode)

	cfg := configtest.New()
	userRepo := repository.NewUserRepository(suite.db)
	userService := services.NewUserService(userRepo, cfg)
	suite.userHandler = NewUserHandler(userService)

	suite.router = gin.New()
	suite.SetupRoutes()
}

func (suite *UserHandlerTestSuite) SetupRoutes() {
	api := suite.router.Group("/api/v1")
	api.POST("/auth/register", suite.userHandler.Register)
	api.POST("/auth/login", suite.userHandler.Login)
}

func (suite *UserHandlerTestSuite) TearDownTest() {
	err := test_helpers.CleanupTestDB(suite.db)
	suite.Require().NoError(err)
}

func (suite *UserHandlerTestSuite) TestRegisterHandler() {
	t := suite.T()

	t.Run("successful registration", func(t *testing.T) {
		reqBody, err := test_helpers.CreateRegisterRequest()
		suite.Require().NoError(err)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	
		var resp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		
		assert.NoError(t, err)
		assert.Equal(t, "user created successfully", resp["message"])
		assert.NotNil(t, resp["user"])
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		reqBody := models.RegisterRequest{
			Username: "a",
			Email:   "invalid-email",
			Password: "123",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func (suite *UserHandlerTestSuite) TestLoginHandler() {
	t := suite.T()

	userRepo := repository.NewUserRepository(suite.db)
	cfg := configtest.New()
	userService := services.NewUserService(userRepo, cfg)

	registerReq, err := test_helpers.CreateRegisterRequest()
	suite.Require().NoError(err)

	_, err = userService.Register(registerReq)
	suite.Require().NoError(err)

	t.Run("successful login", func(t *testing.T) {
		reqBody, err := test_helpers.CreateLoginRequest()
		suite.Require().NoError(err)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resp)

		assert.NoError(t, err)
		assert.Equal(t, "login successful", resp["message"])
		assert.NotNil(t, resp["token"])
		assert.NotNil(t, resp["user"])
	})

	t.Run("invalid credentials", func(t *testing.T) {
		reqBody := models.LoginRequest{
			Email: "john@example.com",
			Password: "wrongpassword",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
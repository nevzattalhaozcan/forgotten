package repository

import (
	"testing"

	"github.com/nevzattalhaozcan/forgotten/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db       *gorm.DB
	userRepo UserRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	var err error

	suite.db, err = test_helpers.SetupTestDB()
	suite.Require().NoError(err)

	suite.userRepo = NewUserRepository(suite.db)
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	err := test_helpers.CleanupTestDB(suite.db)
	suite.Require().NoError(err)
}

func (suite *UserRepositoryTestSuite) TestCreateUser() {
	user, err := test_helpers.CreateUser()
	suite.Require().NoError(err)

	err = suite.userRepo.Create(user)

	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), user.ID)
	assert.NotZero(suite.T(), user.CreatedAt)
}

func (suite *UserRepositoryTestSuite) TestGetUserByEmail() {
	user, err := test_helpers.CreateUser()
	suite.Require().NoError(err)

	err = suite.userRepo.Create(user)
	suite.Require().NoError(err)

	retrievedUser, err := suite.userRepo.GetByEmail(user.Email)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Email, retrievedUser.Email)
	assert.Equal(suite.T(), user.PasswordHash, retrievedUser.PasswordHash)
}

func (suite *UserRepositoryTestSuite) TestGetUserByEmail_NotFound() {
	_, err := suite.userRepo.GetByEmail("nonexistent@example.com")

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrRecordNotFound, err)
}

func (suite *UserRepositoryTestSuite) TestUpdateUser() {
	user, err := test_helpers.CreateUser()
	suite.Require().NoError(err)

	err = suite.userRepo.Create(user)
	suite.Require().NoError(err)

	user.FirstName = "UpdatedFirstName"
	err = suite.userRepo.Update(user)

	assert.NoError(suite.T(), err)

	retrievedUser, err := suite.userRepo.GetByID(user.ID)
	suite.Require().NoError(err)
	assert.Equal(suite.T(), user.FirstName, retrievedUser.FirstName)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

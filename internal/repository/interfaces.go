package repository

import "github.com/nevzattalhaozcan/forgotten/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List(limit, offset int) ([]*models.User, error)
}

type ClubRepository interface {
	Create(club *models.Club) error
	GetByID(id uint) (*models.Club, error)
	Update(club *models.Club) error
	Delete(id uint) error
	List(limit, offset int) ([]*models.Club, error)
	GetByName(name string) (*models.Club, error)
	JoinClub(membership *models.ClubMembership) error
	LeaveClub(clubID, userID uint) error
	ListClubMembers(clubID uint) ([]*models.ClubMembership, error)
	UpdateClubMember(membership *models.ClubMembership) error
	GetClubMemberByUserID(clubID, userID uint) (*models.ClubMembership, error)
}

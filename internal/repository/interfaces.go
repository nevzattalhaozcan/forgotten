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
	GetByEmailIncludingDeleted(email string) (*models.User, error)
	GetByUsernameIncludingDeleted(username string) (*models.User, error)
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

type EventRepository interface {
	Create(event *models.Event) error
    GetClubEvents(clubID uint) ([]models.Event, error)
    GetByID(id uint) (*models.Event, error)
    Update(event *models.Event) error
    Delete(id uint) error
    RSVP(eventID uint, rsvp *models.EventRSVP) error
    GetEventAttendees(eventID uint) ([]models.EventRSVP, error)
}

type BookRepository interface {
	Create(book *models.Book) error
	GetByID(id uint) (*models.Book, error)
	Update(book *models.Book) error
	Delete(id uint) error
	List(limit, offset int) ([]*models.Book, error)
}

type PostRepository interface {
	Create(post *models.Post) error
	GetByID(id uint) (*models.Post, error)
	Update(post *models.Post) error
	Delete(id uint) error
	ListByUserID(userID uint) ([]models.Post, error)
	ListByClubID(clubID uint) ([]models.Post, error)
	ListAll() ([]models.Post, error)
	AddLike(like *models.PostLike) error
	RemoveLike(userID, postID uint) error
	ListLikesByPostID(postID uint) ([]models.PostLike, error)
	HasUserLiked(userID, postID uint) (bool, error)
}
package repository

import (
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"gorm.io/gorm"
)

type clubRepository struct {
	db *gorm.DB
}

func NewClubRepository(db *gorm.DB) *clubRepository {
	return &clubRepository{db: db}
}

func (r *clubRepository) Create(club *models.Club) error {
	return r.db.Create(club).Error
}

func (r *clubRepository) GetByID(id uint) (*models.Club, error) {
	var club models.Club
	err := r.db.
		Preload("Owner").
		Preload("Moderators").
		Preload("Members").
		Preload("Members.User").
		Preload("Posts").
		First(&club, id).Error
	if err != nil {
		return nil, err
	}
	return &club, nil
}

func (r *clubRepository) GetByName(name string) (*models.Club, error) {
	var club models.Club
	err := r.db.Where("name = ?", name).First(&club).Error
	if err != nil {
		return nil, err
	}
	return &club, nil
}

func (r *clubRepository) Update(club *models.Club) error {
	return r.db.Save(club).Error
}

func (r *clubRepository) Delete(id uint) error {
	return r.db.Delete(&models.Club{}, id).Error
}

func (r *clubRepository) List(limit, offset int) ([]*models.Club, error) {
	var clubs []*models.Club
	err := r.db.
		Preload("Owner").
		Limit(limit).
		Offset(offset).
		Find(&clubs).Error
	if err != nil {
		return nil, err
	}
	if len(clubs) == 0 {
		return []*models.Club{}, nil
	}
	return clubs, nil
}

func (r *clubRepository) JoinClub(membership *models.ClubMembership) error {
	return r.db.Create(membership).Error
}

func (r *clubRepository) LeaveClub(clubID, userID uint) error {
	return r.db.Where("club_id = ? AND user_id = ?", clubID, userID).Delete(&models.ClubMembership{}).Error
}

func (r *clubRepository) ListClubMembers(clubID uint) ([]*models.ClubMembership, error) {
	var members []*models.ClubMembership
	err := r.db.
		Where("club_id = ?", clubID).
		Preload("User").
		Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (r *clubRepository) UpdateClubMember(membership *models.ClubMembership) error {
	return r.db.
		Model(&models.ClubMembership{}).
		Where("club_id = ? AND user_id = ?", membership.ClubID, membership.UserID).
		Updates(map[string]interface{}{
			"role":        membership.Role,
			"is_approved": membership.IsApproved,
		}).Error
}

func (r *clubRepository) GetClubMemberByUserID(clubID, userID uint) (*models.ClubMembership, error) {
	var membership models.ClubMembership
	err := r.db.
		Where("club_id = ? AND user_id = ?", clubID, userID).
		Preload("User").
		First(&membership).Error
	if err != nil {
		return nil, err
	}
	return &membership, nil
}

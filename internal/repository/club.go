package repository

import (
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"gorm.io/gorm"
)

type ClubRepository struct {
	db *gorm.DB
}

func NewClubRepository(db *gorm.DB) *ClubRepository {
	return &ClubRepository{db: db}
}

func (r *ClubRepository) Create(club *models.Club) error {
	return r.db.Create(club).Error
}

func (r *ClubRepository) GetByID(id uint) (*models.Club, error) {
	var club models.Club
	err := r.db.First(&club, id).Error
	if err != nil {
		return nil, err
	}
	return &club, nil
}

func (r *ClubRepository) Update(club *models.Club) error {
	return r.db.Save(club).Error
}

func (r *ClubRepository) Delete(id uint) error {
	return r.db.Delete(&models.Club{}, id).Error
}

func (r *ClubRepository) List(limit, offset int) ([]*models.Club, error) {
	var clubs []*models.Club
	err := r.db.Limit(limit).Offset(offset).Find(&clubs).Error
	if err != nil {
		return nil, err
	}
	return clubs, nil
}
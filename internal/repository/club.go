package repository

import (
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type clubRepository struct {
	db *gorm.DB
}

type clubRatingRepository struct {
	db *gorm.DB
}

func NewClubRepository(db *gorm.DB) *clubRepository {
	return &clubRepository{db: db}
}

func NewClubRatingRepository(db *gorm.DB) *clubRatingRepository {
	return &clubRatingRepository{db: db}
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
	return r.db.Omit(clause.Associations).Save(club).Error
}

func (r *clubRepository) UpdateRatingAggregate(clubID uint, avg float32, count int) error {
	return r.db.Model(&models.Club{}).
		Where("id = ?", clubID).
		Updates(map[string]interface{}{
			"rating":        avg,
			"ratings_count": count,
		}).Error
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
	result := r.db.Where("club_id = ? AND user_id = ?", clubID, userID).Delete(&models.ClubMembership{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *clubRepository) CountApprovedMembers(clubID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.ClubMembership{}).
		Where("club_id = ? AND is_approved = true", clubID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *clubRepository) ListClubMembers(clubID uint) ([]*models.ClubMembership, error) {
	var members []*models.ClubMembership
	err := r.db.
		Where("club_id = ? AND is_approved = true", clubID).
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

func (r *clubRatingRepository) UpsertRating(cr *models.ClubRating) error {
	var existing models.ClubRating
	err := r.db.Where("club_id = ? AND user_id = ?", cr.ClubID, cr.UserID).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.Create(cr).Error
	}
	if err != nil {
		return err
	}

	existing.Rating = cr.Rating
	existing.Comment = cr.Comment
	return r.db.Save(&existing).Error
}

func (r *clubRatingRepository) ListByClub(clubID uint, limit, offset int) ([]models.ClubRating, error) {
	var out []models.ClubRating
	err := r.db.Where("club_id = ?", clubID).
		Order("updated_at desc").
		Limit(limit).Offset(offset).
		Find(&out).Error

	return out, err
}

func (r *clubRatingRepository) GetAggregateForClub(clubID uint) (float32, int, error) {
	type agg struct {
		Avg   float64
		Count int64
	}
	var a agg
	err := r.db.Model(&models.ClubRating{}).
		Select("COALESCE(AVG(rating), 0) AS avg, COUNT(*) as count").
		Where("club_id = ?", clubID).
		Scan(&a).Error

	return float32(a.Avg), int(a.Count), err
}

func (r *clubRepository) UpdateMembership(m *models.ClubMembership) error {
	return r.db.Save(m).Error
}

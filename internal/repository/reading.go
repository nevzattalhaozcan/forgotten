package repository

import (
    "errors"

    "github.com/nevzattalhaozcan/forgotten/internal/models"
    "gorm.io/gorm"
)

type readingRepository struct {
	db *gorm.DB
}

func NewReadingRepository(db *gorm.DB) *readingRepository {
	return &readingRepository{db: db}
}

func (r *readingRepository) UpsertUserProgress(p *models.UserBookProgress) error {
    var existing models.UserBookProgress
    err := r.db.Where("user_id = ? AND book_id = ?", p.UserID, p.BookID).First(&existing).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return r.db.Create(p).Error
    }
    if err != nil {
        return err
    }
    p.ID = existing.ID
    return r.db.Model(&existing).Updates(p).Error
}

func (r *readingRepository) GetUserBookProgress(userID, bookID uint) (*models.UserBookProgress, error) {
    var p models.UserBookProgress
    if err := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&p).Error; err != nil {
        return nil, err
    }
    return &p, nil
}

func (r *readingRepository) ListUserProgress(userID uint) ([]*models.UserBookProgress, error) {
    var out []*models.UserBookProgress
    err := r.db.Where("user_id = ?", userID).Order("updated_at DESC").Find(&out).Error
    return out, err
}

func (r *readingRepository) ListUserFinished(userID uint) ([]*models.UserBookProgress, error) {
    var out []*models.UserBookProgress
    err := r.db.Where("user_id = ? AND status = ?", userID, models.ReadingFinished).
        Order("finished_at DESC").Find(&out).Error
    return out, err
}

func (r *readingRepository) AppendLog(l *models.ReadingLog) error {
    return r.db.Create(l).Error
}

func (r *readingRepository) ListLogsByUserAndBook(userID, bookID uint) ([]models.ReadingLog, error) {
    var logs []models.ReadingLog
    err := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).
        Order("created_at ASC").Find(&logs).Error
    return logs, err
}

type clubReadingRepository struct { db *gorm.DB }
func NewClubReadingRepository(db *gorm.DB) *clubReadingRepository { return &clubReadingRepository{db: db} }

func (r *clubReadingRepository) CreateAssignment(a *models.ClubBookAssignment) error {
    return r.db.Create(a).Error
}

func (r *clubReadingRepository) CompleteAssignment(assignmentID uint) error {
    return r.db.Model(&models.ClubBookAssignment{}).Where("id = ?", assignmentID).
        Updates(map[string]interface{}{"status": models.ClubAssignmentCompleted, "completed_at": gorm.Expr("NOW()")}).Error
}

func (r *clubReadingRepository) GetActiveAssignment(clubID uint) (*models.ClubBookAssignment, error) {
    var a models.ClubBookAssignment
    if err := r.db.Where("club_id = ? AND status = ?", clubID, models.ClubAssignmentActive).
        Order("created_at DESC").First(&a).Error; err != nil {
        return nil, err
    }
    return &a, nil
}

func (r *clubReadingRepository) ListAssignments(clubID uint) ([]models.ClubBookAssignment, error) {
    var out []models.ClubBookAssignment
    err := r.db.Where("club_id = ?", clubID).Order("created_at DESC").Find(&out).Error
    return out, err
}

func (r *clubReadingRepository) UpdateAssignment(a *models.ClubBookAssignment) error {
    return r.db.Save(a).Error
}
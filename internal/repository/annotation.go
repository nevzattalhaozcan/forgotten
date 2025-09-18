package repository

import (
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"gorm.io/gorm"
)

type AnnotationRepository struct {
	db *gorm.DB
}

func NewAnnotationRepository(db *gorm.DB) *AnnotationRepository {
	return &AnnotationRepository{db: db}
}

func (r *AnnotationRepository) Create(annotation *models.Annotation) error {
	return r.db.Create(annotation).Error
}

func (r *AnnotationRepository) GetByID(id uint) (*models.Annotation, error) {
	var annotation models.Annotation
	if err := r.db.First(&annotation, id).Error; err != nil {
		return nil, err
	}
	return &annotation, nil
}

func (r *AnnotationRepository) Update(annotation *models.Annotation) error {
	return r.db.Save(annotation).Error
}

func (r *AnnotationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Annotation{}, id).Error
}

func (r *AnnotationRepository) ListByUserID(userID uint) ([]models.Annotation, error) {
	var annotations []models.Annotation
	if err := r.db.Where("user_id = ?", userID).Find(&annotations).Error; err != nil {
		return nil, err
	}
	return annotations, nil
}

func (r *AnnotationRepository) LikeAnnotation(like *models.AnnotationLike) error {
	return r.db.Create(like).Error
}

func (r *AnnotationRepository) UnlikeAnnotation(userID, annotationID uint) error {
	return r.db.Where("user_id = ? AND annotation_id = ?", userID, annotationID).Delete(&models.AnnotationLike{}).Error
}

func (r *AnnotationRepository) CountLikes(annotationID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.AnnotationLike{}).Where("annotation_id = ?", annotationID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *AnnotationRepository) ListAnnotationLikes(annotationID uint) ([]models.AnnotationLike, error) {
	var likes []models.AnnotationLike
	if err := r.db.Where("annotation_id = ?", annotationID).Find(&likes).Error; err != nil {
		return nil, err
	}
	return likes, nil
}

func (r *AnnotationRepository) ListAll() ([]models.Annotation, error) {
	var annotations []models.Annotation
	if err := r.db.Find(&annotations).Error; err != nil {
		return nil, err
	}
	return annotations, nil
}
package repository

import (
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *CommentRepository) GetByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	if err := r.db.First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepository) Update(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

func (r *CommentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Comment{}, id).Error
}

func (r *CommentRepository) ListByPostID(postID uint) ([]models.Comment, error) {
	var comments []models.Comment
	if err := r.db.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepository) ListByUserID(userID uint) ([]models.Comment, error) {
	var comments []models.Comment
	if err := r.db.Where("user_id = ?", userID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepository) LikeComment(like *models.CommentLike) error {
	return r.db.Create(like).Error
}

func (r *CommentRepository) UnlikeComment(userID, commentID uint) error {
	return r.db.Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&models.CommentLike{}).Error
}

func (r *CommentRepository) CountLikes(commentID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.CommentLike{}).Where("comment_id = ?", commentID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *CommentRepository) ListCommentLikes(commentID uint) ([]models.CommentLike, error) {
	var likes []models.CommentLike
	if err := r.db.Where("comment_id = ?", commentID).Find(&likes).Error; err != nil {
		return nil, err
	}
	return likes, nil
}
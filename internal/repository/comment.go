package repository

import (
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *commentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) GetByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	if err := r.db.
		Preload("User").
		Preload("Likes").
		Preload("Likes.User").
		First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) Update(comment *models.Comment) error {
	return r.db.Omit(clause.Associations).Save(comment).Error
}

func (r *commentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Comment{}, id).Error
}

func (r *commentRepository) ListByPostID(postID uint) ([]models.Comment, error) {
	var comments []models.Comment
	if err := r.db.
		Preload("User").
		Preload("Likes").
		Where("post_id = ?", postID).
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepository) ListByUserID(userID uint) ([]models.Comment, error) {
	var comments []models.Comment
	if err := r.db.
		Preload("User").
		Preload("Likes").
		Where("user_id = ?", userID).
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepository) LikeComment(like *models.CommentLike) error {
	return r.db.Create(like).Error
}

func (r *commentRepository) UnlikeComment(userID, commentID uint) error {
	return r.db.Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&models.CommentLike{}).Error
}

func (r *commentRepository) CountLikes(commentID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.CommentLike{}).Where("comment_id = ?", commentID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *commentRepository) ListCommentLikes(commentID uint) ([]models.CommentLike, error) {
	var likes []models.CommentLike
	if err := r.db.
		Preload("User").
		Preload("Comment").
		Where("comment_id = ?", commentID).
		Find(&likes).Error; err != nil {
		return nil, err
	}
	return likes, nil
}

func (r *commentRepository) HasUserLiked(userID, commentID uint) (bool, error) {
    var count int64
    err := r.db.Model(&models.CommentLike{}).
        Where("user_id = ? AND comment_id = ?", userID, commentID).
        Count(&count).Error
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

func (r *commentRepository) UpdateLikesCount(commentID uint, count int) error {
    return r.db.Model(&models.Comment{}).
        Where("id = ?", commentID).
        Update("likes_count", count).Error
}
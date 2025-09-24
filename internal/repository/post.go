package repository

import (
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *postRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) GetByID(id uint) (*models.Post, error) {
	var post models.Post
	if err := r.db.
		Preload("User").
		Preload("Comments").
		Preload("Likes").
		Preload("Likes.User").
		First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}

func (r *postRepository) ListByUserID(userID uint) ([]models.Post, error) {
	var posts []models.Post
	if err := r.db.
		Preload("User").
		Preload("Comments").
		Preload("Likes").
		Where("user_id = ?", userID).
		Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) ListByClubID(clubID uint) ([]models.Post, error) {
	var posts []models.Post
	if err := r.db.
		Preload("User").
		Preload("Comments").
		Preload("Likes").
		Where("club_id = ?", clubID).
		Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) ListAll() ([]models.Post, error) {
	var posts []models.Post
	if err := r.db.
		Preload("User").
		Preload("Comments").
		Preload("Likes").
		Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) AddLike(like *models.PostLike) error {
	return r.db.Create(like).Error
}

func (r *postRepository) RemoveLike(userID, postID uint) error {
	result := r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&models.PostLike{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *postRepository) CountLikes(postID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.PostLike{}).Where("post_id = ?", postID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *postRepository) ListLikesByPostID(postID uint) ([]models.PostLikeResponse, error) {
	var likes []models.PostLike
	if err := r.db.
		Preload("User").
		Where("post_id = ?", postID).
		Find(&likes).Error; err != nil {
		return nil, err
	}

	res := make([]models.PostLikeResponse, 0, len(likes))
    for _, l := range likes {
        res = append(res, l.ToResponse())
    }
    return res, nil
}

func (r *postRepository) HasUserLiked(userID, postID uint) (bool, error) {
    var count int64
    err := r.db.Model(&models.PostLike{}).
        Where("user_id = ? AND post_id = ?", userID, postID).
        Count(&count).Error
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

func (r *postRepository) UpdateLikesCount(postID uint, count int) error {
	return r.db.Model(&models.Post{}).
		Where("id = ?", postID).
		Update("likes_count", count).Error
}
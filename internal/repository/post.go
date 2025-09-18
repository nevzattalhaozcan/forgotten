package repository

import (
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepository) GetByID(id uint) (*models.Post, error) {
	var post models.Post
	if err := r.db.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

func (r *PostRepository) Delete(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}

func (r *PostRepository) ListByUserID(userID uint) ([]models.Post, error) {
	var posts []models.Post
	if err := r.db.Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepository) ListByClubID(clubID uint) ([]models.Post, error) {
	var posts []models.Post
	if err := r.db.Where("club_id = ?", clubID).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepository) ListAll() ([]models.Post, error) {
	var posts []models.Post
	if err := r.db.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepository) LikePost(like *models.PostLike) error {
	return r.db.Create(like).Error
}

func (r *PostRepository) UnlikePost(userID, postID uint) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&models.PostLike{}).Error
}

func (r *PostRepository) CountLikes(postID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.PostLike{}).Where("post_id = ?", postID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *PostRepository) ListLikesByPostID(postID uint) ([]models.PostLike, error) {
	var likes []models.PostLike
	if err := r.db.Where("post_id = ?", postID).Find(&likes).Error; err != nil {
		return nil, err
	}
	return likes, nil
}
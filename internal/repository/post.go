package repository

import (
	"strconv"
	"time"

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
		Preload("Likes.User").
		Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) ListPublicPosts() ([]models.Post, error) {
    var posts []models.Post
    if err := r.db.
        Preload("User").
        Preload("Club"). // Add this to show which club the post belongs to
        Preload("Comments").
        Preload("Likes").
        Joins("JOIN clubs ON posts.club_id = clubs.id").
        Where("clubs.is_private = ?", false). // Fix: false for public clubs
        Order("posts.likes_count DESC, posts.created_at DESC"). // Popular first
        Limit(20). // Limit for homepage
        Find(&posts).Error; err != nil {
        return nil, err
    }
    return posts, nil
}

func (r *postRepository) ListPopularPublicPosts(limit int) ([]models.Post, error) {
    var posts []models.Post
    if err := r.db.
        Preload("User").
        Preload("Club").
        Preload("Comments").
        Preload("Likes").
        Joins("JOIN clubs ON posts.club_id = clubs.id").
        Where("clubs.is_private = ? AND posts.created_at > ?", false, time.Now().AddDate(0, 0, -30)). // Last 30 days
        Order("posts.likes_count DESC, posts.comments_count DESC, posts.created_at DESC").
        Limit(limit).
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

func (r *postRepository) VoteOnPoll(vote *models.PollVote) error {
    return r.db.Create(vote).Error
}

func (r *postRepository) RemoveVoteFromPoll(postID, userID uint, optionID string) error {
    return r.db.Where("post_id = ? AND user_id = ? AND option_id = ?", postID, userID, optionID).
        Delete(&models.PollVote{}).Error
}

func (r *postRepository) GetUserPollVotes(postID, userID uint) ([]models.PollVote, error) {
    var votes []models.PollVote
    err := r.db.Where("post_id = ? AND user_id = ?", postID, userID).Find(&votes).Error
    return votes, err
}

// TODO: Find a way to update poll vote counts
func (r *postRepository) UpdatePollVoteCounts(postID uint) error {
    return nil
}

func (r *postRepository) GetPostsByType(postType string, limit, offset int) ([]models.Post, error) {
    var posts []models.Post
    err := r.db.Where("type = ?", postType).
        Preload("User").
        Preload("Club").
        Limit(limit).
        Offset(offset).
        Find(&posts).Error
    return posts, err
}

func (r *postRepository) GetReviewPostsByBookID(bookID uint) ([]models.Post, error) {
    var posts []models.Post
    err := r.db.Where("type = ? AND type_data->>'book_id' = ?", "review", strconv.Itoa(int(bookID))).
        Preload("User").
        Find(&posts).Error
    return posts, err
}

func (r *postRepository) GetPollPostsByClubID(clubID uint, includeExpired bool) ([]models.Post, error) {
    query := r.db.Where("type = ? AND club_id = ?", "poll", clubID)
    
    if !includeExpired {
        query = query.Where("(type_data->>'expires_at' IS NULL OR type_data->>'expires_at'::timestamp > NOW())")
    }
    
    var posts []models.Post
    err := query.Preload("User").Find(&posts).Error
    return posts, err
}
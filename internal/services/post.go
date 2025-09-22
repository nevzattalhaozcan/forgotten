package services

import (
	"errors"
	"strings"

	"github.com/lib/pq"
	"github.com/nevzattalhaozcan/forgotten/internal/config"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/internal/repository"
	"gorm.io/gorm"
)

type PostService struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
	clubRepo repository.ClubRepository
	config   *config.Config
}

func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository, clubRepo repository.ClubRepository, config *config.Config) *PostService {
	return &PostService{
		postRepo: postRepo,
		userRepo: userRepo,
		clubRepo: clubRepo,
		config:   config,
	}
}

func (s *PostService) CreatePost(userID uint, req *models.CreatePostRequest) (*models.PostResponse, error) {
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	_, err = s.clubRepo.GetByID(req.ClubID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("club not found")
		}
		return nil, err
	}

	post := &models.Post{
		Title:   req.Title,
		Content: req.Content,
		Type:    req.Type,
		ClubID:  req.ClubID,
		UserID:  userID,
	}

	if err := s.postRepo.Create(post); err != nil {
		return nil, err
	}

	response := post.ToResponse()
	return &response, nil
}

func (s *PostService) GetPostByID(id uint) (*models.PostResponse, error) {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	response := post.ToResponse()
	return &response, nil
}

func (s *PostService) UpdatePost(id uint, req *models.UpdatePostRequest) (*models.PostResponse, error) {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	if req.Title != nil {
		post.Title = *req.Title
	}
	if req.Content != nil {
		post.Content = *req.Content
	}
	if req.Type != nil {
		post.Type = *req.Type
	}
	if req.IsPinned != nil {
		post.IsPinned = *req.IsPinned
	}

	if err := s.postRepo.Update(post); err != nil {
		return nil, err
	}

	response := post.ToResponse()
	return &response, nil
}

func (s *PostService) DeletePost(id uint) error {
	_, err := s.postRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return err
	}

	return s.postRepo.Delete(id)
}

func (s *PostService) ListPostsByUserID(userID uint) ([]models.PostResponse, error) {
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	posts, err := s.postRepo.ListByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no posts found for this user")
		}
		return nil, err
	}

	var responses []models.PostResponse
	for _, post := range posts {
		responses = append(responses, post.ToResponse())
	}
	return responses, nil
}

func (s *PostService) ListPostsByClubID(clubID uint) ([]models.PostResponse, error) {
	_, err := s.clubRepo.GetByID(clubID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("club not found")
		}
		return nil, err
	}

	posts, err := s.postRepo.ListByClubID(clubID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no posts found for this club")
		}
		return nil, err
	}

	var responses []models.PostResponse
	for _, post := range posts {
		responses = append(responses, post.ToResponse())
	}
	return responses, nil
}

func (s *PostService) ListAllPosts() ([]models.PostResponse, error) {
	posts, err := s.postRepo.ListAll()
	if err != nil {
		return nil, err
	}

	var responses []models.PostResponse
	for _, post := range posts {
		responses = append(responses, post.ToResponse())
	}
	return responses, nil
}

func (s *PostService) LikePost(userID, postID uint) error {
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return err
	}

	like := &models.PostLike{
		UserID: userID,
		PostID: postID,
	}

	hasLiked, err := s.postRepo.HasUserLiked(userID, postID)
	if err != nil {
		return err
	}
	if hasLiked {
		return errors.New("user has already liked this post")
	}

	err = s.postRepo.AddLike(like)
	if err != nil {
		if isUniqueConstraintError(err, "idx_user_post_like") {
			return errors.New("user has already liked this post")
		}
		return err
	}

	// Only increment the count if the like was successfully added
	post.LikesCount++
	if err := s.postRepo.Update(post); err != nil {
		// If updating the count fails, we should remove the like to maintain consistency
		s.postRepo.RemoveLike(userID, postID)
		return err
	}

	return nil
}

//TODO: fix success issue for the first unlike action
func (s *PostService) UnlikePost(userID, postID uint) error {
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return err
	}

	hasLiked, err := s.postRepo.HasUserLiked(userID, postID)
	if err != nil {
		return err
	}
	if !hasLiked {
		return errors.New("user has not liked this post")
	}

	err = s.postRepo.RemoveLike(userID, postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user has not liked this post")
		}
		return err
	}

	// Only decrement if the like was successfully removed and count is positive
	if post.LikesCount > 0 {
		post.LikesCount--
		if err := s.postRepo.Update(post); err != nil {
			return err
		}
	}

	return nil
}

func (s *PostService) ListLikesByPostID(postID uint) ([]models.PostLike, error) {
	_, err := s.postRepo.GetByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	likes, err := s.postRepo.ListLikesByPostID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no likes found for this post")
		}
		return nil, err
	}

	return likes, nil
}

func isUniqueConstraintError(err error, constraintName string) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		// PostgreSQL unique constraint violation code is 23505
		if pqErr.Code == "23505" {
			// Check if the constraint name matches
			return strings.Contains(string(pqErr.Constraint), constraintName) ||
				strings.Contains(pqErr.Message, constraintName) ||
				strings.Contains(pqErr.Detail, constraintName)
		}
	}
	return false
}

package services

import (
	"errors"

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

	created, err := s.postRepo.GetByID(post.ID)
    if err != nil {
        return nil, err
    }
    response := created.ToResponse()
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

	created, err := s.postRepo.GetByID(post.ID)
    if err != nil {
        return nil, err
    }
    response := created.ToResponse()
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

func (s *PostService) ListPublicPosts() ([]models.PostResponse, error) {
	posts, err := s.postRepo.ListPublicPosts()
	if err != nil {
		return nil, err
	}

	var responses []models.PostResponse
	for _, post := range posts {
		responses = append(responses, post.ToResponse())
	}
	return responses, nil
}

func (s *PostService) ListPopularPublicPosts(limit int) ([]models.PostResponse, error) {
	posts, err := s.postRepo.ListPopularPublicPosts(limit)
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

	hasLiked, err := s.postRepo.HasUserLiked(userID, postID)
	if err != nil {
		return err
	}
	if hasLiked {
		return errors.New("user has already liked this post")
	}

	err = s.postRepo.AddLike(&models.PostLike{
		UserID: userID,
		PostID: postID,
	})
	if err != nil {
		return err
	}

	count, err := s.postRepo.CountLikes(postID)
	if err != nil {
		return err
	}

	post.LikesCount = int(count)
	if err := s.postRepo.Update(post); err != nil {
		return err
	}

	return s.postRepo.UpdateLikesCount(postID, int(count))
}

func (s *PostService) UnlikePost(userID, postID uint) error {
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	_, err = s.postRepo.GetByID(postID)
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

	if err := s.postRepo.RemoveLike(userID, postID); err != nil {
		return err
	}

	count, err := s.postRepo.CountLikes(postID)
	if err != nil {
		return err
	}

	return s.postRepo.UpdateLikesCount(postID, int(count))
}

func (s *PostService) ListLikesByPostID(postID uint) ([]models.PostLikeResponse, error) {
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

package services

import (
	"errors"

	"github.com/nevzattalhaozcan/forgotten/internal/config"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/internal/repository"
	"gorm.io/gorm"
)

type CommentService struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
	userRepo    repository.UserRepository
	config      *config.Config
}

func NewCommentService(commentRepo repository.CommentRepository, postRepo repository.PostRepository, userRepo repository.UserRepository, config *config.Config) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		userRepo:    userRepo,
		config:      config,
	}
}

func (s *CommentService) CreateComment(postID, userID uint, req *models.CreateCommentRequest) (*models.CommentResponse, error) {
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	_, err = s.postRepo.GetByID(postID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	comment := &models.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	response := comment.ToResponse()
	return &response, nil
}

func (s *CommentService) GetCommentByID(id uint) (*models.CommentResponse, error) {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}

	response := comment.ToResponse()
	return &response, nil
}

func (s *CommentService) UpdateComment(id uint, req *models.UpdateCommentRequest) (*models.CommentResponse, error) {
	comment, err := s.commentRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}

	if req.Content != nil {
		comment.Content = *req.Content
	}

	if err := s.commentRepo.Update(comment); err != nil {
		return nil, err
	}

	response := comment.ToResponse()
	return &response, nil
}

func (s *CommentService) DeleteComment(id uint) error {
	_, err := s.commentRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("comment not found")
		}
		return err
	}

	return s.commentRepo.Delete(id)
}

func (s *CommentService) ListCommentsByPostID(postID uint) ([]models.Comment, error) {
	_, err := s.postRepo.GetByID(postID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	comments, err := s.commentRepo.ListByPostID(postID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *CommentService) ListCommentsByUserID(userID uint) ([]models.Comment, error) {
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	comments, err := s.commentRepo.ListByUserID(userID)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *CommentService) LikeComment(userID, commentID uint) error {
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("user not found")
		}
		return err
	}

	comment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("comment not found")
		}
		return err
	}

	hasLiked, err := s.commentRepo.HasUserLiked(userID, commentID)
	if err != nil {
		return err
	}
	if hasLiked {
		return errors.New("user has already liked this comment")
	}

	if err := s.commentRepo.LikeComment(&models.CommentLike{
		UserID:    userID,
		CommentID: commentID,
	}); err != nil {
		return err
	}

	count, err := s.commentRepo.CountLikes(commentID)
	if err != nil {
		return err
	}

	comment.LikesCount = int(count)
	if err := s.commentRepo.Update(comment); err != nil {
		return err
	}

	return s.commentRepo.UpdateLikesCount(commentID, int(count))
}

func (s *CommentService) UnlikeComment(userID, commentID uint) error {
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("user not found")
		}
		return err
	}

	_, err = s.commentRepo.GetByID(commentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("comment not found")
		}
		return err
	}

	hasLiked, err := s.commentRepo.HasUserLiked(userID, commentID)
	if err != nil {
		return err
	}
	if !hasLiked {
		return errors.New("user has not liked this comment")
	}

	if err := s.commentRepo.UnlikeComment(userID, commentID); err != nil {
		return err
	}

	count, err := s.commentRepo.CountLikes(commentID)
	if err != nil {
		return err
	}

	return s.commentRepo.UpdateLikesCount(commentID, int(count))
}

func (s *CommentService) ListLikesByCommentID(commentID uint) ([]models.CommentLike, error) {
	_, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}

	likes, err := s.commentRepo.ListCommentLikes(commentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("no likes found for this comment")
		}
		return nil, err
	}

	return likes, nil
}

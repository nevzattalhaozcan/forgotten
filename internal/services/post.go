package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/nevzattalhaozcan/forgotten/internal/config"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/internal/repository"
	"gorm.io/gorm"
)

type PostService struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
	clubRepo repository.ClubRepository
	bookRepo repository.BookRepository
	db       *gorm.DB
	config   *config.Config
}

func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository, clubRepo repository.ClubRepository, bookRepo repository.BookRepository, db *gorm.DB, config *config.Config) *PostService {
	return &PostService{
		postRepo: postRepo,
		userRepo: userRepo,
		clubRepo: clubRepo,
		bookRepo: bookRepo,
		db:       db,
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

	if req.TypeData != nil {
		switch req.Type {
		case "review":
			if reviewData, ok := req.TypeData.(map[string]interface{}); ok {
				if bookIDFloat, exists := reviewData["book_id"].(float64); exists {
					if bookID := uint(bookIDFloat); bookID > 0 {
						if book, err := s.bookRepo.GetByID(bookID); err == nil {
							reviewData["book_title"] = book.Title
							reviewData["book_author"] = book.Author
						}
					}
				}
				req.TypeData = reviewData
			}
		case "annotation":
			if annotationData, ok := req.TypeData.(map[string]interface{}); ok {
				if bookIDFloat, exists := annotationData["book_id"].(float64); exists {
					if bookID := uint(bookIDFloat); bookID > 0 {
						if book, err := s.bookRepo.GetByID(bookID); err == nil {
							annotationData["book_title"] = book.Title
							annotationData["book_author"] = book.Author
						}
					}
				}
				req.TypeData = annotationData
			}
		case "poll":
			if pollData, ok := req.TypeData.(map[string]interface{}); ok {
				if options, exists := pollData["options"].([]interface{}); exists {
					for i, option := range options {
						if optionMap, ok := option.(map[string]interface{}); ok {
							optionMap["id"] = fmt.Sprintf("opt_%d", i+1)
							optionMap["votes"] = 0
						}
					}
				}
				req.TypeData = pollData
			}
		}

		typeDataBytes, err := json.Marshal(req.TypeData)
        if err != nil {
            return nil, fmt.Errorf("invalid type data: %v", err)
        }
        post.TypeData = models.PostTypeData(typeDataBytes)
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

func (s *PostService) VoteOnPoll(postID, userID uint, req *models.PollVoteRequest) error {
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return err
	}

	if post.Type != "poll" {
		return errors.New("post is not a poll")
	}

	pollData, err := post.GetPollData()
	if err != nil || pollData == nil {
		return errors.New("invalid poll data")
	}

	if pollData.ExpiresAt != nil && pollData.ExpiresAt.Before(time.Now()) {
		return errors.New("poll has expired")
	}

	existingVotes, err := s.postRepo.GetUserPollVotes(postID, userID)
	if err != nil {
		return err
	}

	if len(existingVotes) > 0 && !pollData.AllowMultiple {
		return errors.New("multiple votes not allowed")
	}

	validOptions := make(map[string]bool)
	for _, option := range pollData.Options {
		validOptions[option.ID] = true
	}

	for _, optionID := range req.OptionIDs {
		if !validOptions[optionID] {
			return fmt.Errorf("invalid option ID: %s", optionID)
		}
	}

	if !pollData.AllowMultiple {
		for _, vote := range existingVotes {
			if err := s.postRepo.RemoveVoteFromPoll(postID, userID, vote.OptionID); err != nil {
				return err
			}
		}
	}

	for _, optionID := range req.OptionIDs {
		vote := &models.PollVote{
			PostID:   postID,
			UserID:   userID,
			OptionID: optionID,
		}
		if err := s.postRepo.VoteOnPoll(vote); err != nil {
			return err
		}
	}

	return s.updatePollCounts(postID)
}

func (s *PostService) updatePollCounts(postID uint) error {
	var votes []models.PollVote
	s.db.Where("post_id = ?", postID).Find(&votes)

	voteCounts := make(map[string]int)
	for _, vote := range votes {
		voteCounts[vote.OptionID]++
	}

	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return err
	}

	pollData, err := post.GetPollData()
	if err != nil {
		return err
	}

	for i := range pollData.Options {
		pollData.Options[i].Votes = voteCounts[pollData.Options[i].ID]
	}

	typeDataBytes, _ := json.Marshal(pollData)
	post.TypeData = models.PostTypeData(typeDataBytes)

	return s.postRepo.Update(post)
}

func (s *PostService) GetPostByIDForUser(postID, userID uint) (*models.PostResponse, error) {
    post, err := s.postRepo.GetByID(postID)
    if err != nil {
        return nil, err
    }
    
    response := post.ToResponse()
    
    if post.Type == "poll" {
        userVotes, err := s.postRepo.GetUserPollVotes(postID, userID)
        if err == nil {
            response.UserVoted = len(userVotes) > 0
            for _, vote := range userVotes {
                response.UserVotes = append(response.UserVotes, vote.OptionID)
            }
        }
    }
    
    return &response, nil
}

func (s *PostService) GetReviewsByBook(bookID uint) ([]models.PostResponse, error) {
	_, err := s.bookRepo.GetByID(bookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	posts, err := s.postRepo.GetReviewPostsByBookID(bookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no reviews found for this book")
		}
		return nil, err
	}

	var responses []models.PostResponse
	for _, post := range posts {
		responses = append(responses, post.ToResponse())
	}
	return responses, nil
}

func (s *PostService) GetPostsByType(postType string, limit, offset int) ([]models.PostResponse, error) {
	posts, err := s.postRepo.GetPostsByType(postType, limit, offset)
	if err != nil {
		return nil, err
	}

	var responses []models.PostResponse
	for _, post := range posts {
		responses = append(responses, post.ToResponse())
	}
	return responses, nil
}

func (s *PostService) GetPollPostsByClubID(clubID uint, includeExpired bool) ([]models.PostResponse, error) {
	_, err := s.clubRepo.GetByID(clubID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("club not found")
		}
		return nil, err
	}

	posts, err := s.postRepo.GetPollPostsByClubID(clubID, includeExpired)
	if err != nil {
		return nil, err
	}

	var responses []models.PostResponse
	for _, post := range posts {
		responses = append(responses, post.ToResponse())
	}
	return responses, nil
}

func (s *PostService) RemoveVoteFromPoll(postID, userID uint, optionID string) error {
	_, err := s.postRepo.GetByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return err
	}

	_, err = s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return s.postRepo.RemoveVoteFromPoll(postID, userID, optionID)
}

func (s *PostService) GetUserPollVotes(postID, userID uint) ([]models.PollVote, error) {
	_, err := s.postRepo.GetByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	_, err = s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return s.postRepo.GetUserPollVotes(postID, userID)
}
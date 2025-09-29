package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/internal/services"
)

type PostHandler struct {
	postService *services.PostService
	validator   *validator.Validate
}

func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
		validator:   validator.New(),
	}
}

// @Summary Create a new post
// @Description Create a new post associated with the authenticated user
// @Tags Posts
// @Accept json
// @Produce json
// @Param post body models.CreatePostRequest true "Post data"
// @Success 201 {object} map[string]interface{} "Post created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} models.ErrorResponse
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	useridRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, ok := useridRaw.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := h.postService.CreatePost(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "post created successfully",
		"post":    post,
	})
}

// @Summary Get a post by ID
// @Description Retrieve a post by its ID
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]interface{} "Post retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Post not found"
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id} [get]
func (h *PostHandler) GetPostByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	post, err := h.postService.GetPostByID(uint(id))
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

// @Summary Update a post
// @Description Update a post by its ID
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param post body models.UpdatePostRequest true "Updated post data"
// @Success 200 {object} models.Post "Post updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Post not found"
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := h.postService.UpdatePost(uint(id), &req)
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "post updated successfully",
		"post":    post,
	})
}

// @Summary Delete a post
// @Description Delete a post by its ID
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 204 {object} nil "Post deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Post not found"
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	if err := h.postService.DeletePost(uint(id)); err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary List posts by user ID
// @Description Retrieve all posts created by a specific user
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} models.Post "Posts retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id}/posts [get]
func (h *PostHandler) ListPostsByUserID(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	posts, err := h.postService.ListPostsByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// @Summary List posts by club ID
// @Description Retrieve all posts associated with a specific club
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Club ID"
// @Success 200 {array} models.Post "Posts retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} models.ErrorResponse
// @Router /clubs/{id}/posts [get]
func (h *PostHandler) ListPostsByClubID(c *gin.Context) {
	clubIDParam := c.Param("id")
	clubID, err := strconv.ParseUint(clubIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid club ID"})
		return
	}

	posts, err := h.postService.ListPostsByClubID(uint(clubID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// @Summary List all posts
// @Description Retrieve all posts in the system
// @Tags Posts
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{} "Posts retrieved successfully"
// @Failure 500 {object} models.ErrorResponse
// @Router /posts [get]
func (h *PostHandler) ListAllPosts(c *gin.Context) {
	posts, err := h.postService.ListAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// @Summary List public posts
// @Description Retrieve all posts from public clubs
// @Tags Posts
// @Accept json
// @Produce json
// @Success 200 {array} models.Post "Public posts retrieved successfully"
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/public [get]
func (h *PostHandler) ListPublicPosts(c *gin.Context) {
	posts, err := h.postService.ListPublicPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// @Summary List popular public posts
// @Description Retrieve popular posts from public clubs based on number of likes
// @Tags Posts
// @Accept json
// @Produce json
// @Success 200 {array} models.Post "Popular public posts retrieved successfully"
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/popular [get]
func (h *PostHandler) ListPopularPublicPosts(c *gin.Context) {
	posts, err := h.postService.ListPopularPublicPosts(20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// @Summary Like a post
// @Description Like a post by its ID
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]interface{} "Post liked successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Post not found"
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id}/like [post]
func (h *PostHandler) LikePost(c *gin.Context) {
	postIdParam := c.Param("id")
	postID, err := strconv.ParseUint(postIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	useridRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, ok := useridRaw.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.postService.LikePost(userID, uint(postID)); err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if strings.Contains(err.Error(), "already liked") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post liked successfully"})
}

// @Summary Unlike a post
// @Description Unlike a post by its ID
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]interface{} "Post unliked successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Post not found"
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id}/unlike [post]
func (h *PostHandler) UnlikePost(c *gin.Context) {
	postIdParam := c.Param("id")
	postID, err := strconv.ParseUint(postIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	useridRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, ok := useridRaw.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.postService.UnlikePost(userID, uint(postID)); err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if strings.Contains(err.Error(), "has not liked") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post unliked successfully"})
}

// @Summary List likes by post ID
// @Description Retrieve all likes associated with a specific post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {array} models.PostLikeResponse "Likes retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Post not found or no likes found"
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id}/likes [get]
func (h *PostHandler) ListLikesByPostID(c *gin.Context) {
	postIdParam := c.Param("id")
	postID, err := strconv.ParseUint(postIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	likes, err := h.postService.ListLikesByPostID(uint(postID))
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "no likes found for this post" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"likes": likes})
}

// @Summary Vote on a poll
// @Description Vote on a poll post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param vote body models.PollVoteRequest true "Vote data"
// @Success 200 {object} models.SuccessResponse "Vote recorded successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Post not found"
// @Router /posts/{id}/vote [post]
func (h *PostHandler) VoteOnPoll(c *gin.Context) {
    postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
        return
    }

    userID, _ := c.Get("user_id")
    
    var req models.PollVoteRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.postService.VoteOnPoll(uint(postID), userID.(uint), &req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "vote recorded successfully"})
}

// @Summary Get reviews by book
// @Description Get all review posts for a specific book
// @Tags Posts
// @Produce json
// @Param book_id query int true "Book ID"
// @Success 200 {array} models.PostResponse "Reviews retrieved successfully"
// @Router /posts/reviews [get]
func (h *PostHandler) GetReviewsByBook(c *gin.Context) {
    bookID, err := strconv.ParseUint(c.Query("book_id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "book_id is required"})
        return
    }

    reviews, err := h.postService.GetReviewsByBook(uint(bookID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

// @Summary Get posts by type
// @Description Get posts filtered by type with pagination
// @Tags Posts
// @Produce json
// @Param type query string true "Post type (e.g., 'announcement', 'discussion', 'poll', 'review')"
// @Param limit query int true "Number of posts to return"
// @Param offset query int true "Number of posts to skip"
// @Success 200 {array} models.PostResponse "Posts retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/by-type [get]
func (h *PostHandler) GetPostsByType(c *gin.Context) {
	postType := c.Query("type")
	if postType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type is required"})
		return
	}
	limitStr := c.Query("limit")
	if limitStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit is required"})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be an integer"})
		return
	}

	offsetStr := c.Query("offset")
	if offsetStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "offset is required"})
		return
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "offset must be an integer"})
		return
	}

	posts, err := h.postService.GetPostsByType(postType, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// @Summary Get poll posts by club ID
// @Description Retrieve all poll posts associated with a specific club
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Club ID"
// @Param include_expired query bool false "Include expired polls"
// @Success 200 {array} models.Post "Poll posts retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} models.ErrorResponse
// @Router /clubs/{id}/posts/polls [get]
func (h *PostHandler) GetPollPostsByClubID(c *gin.Context) {
	clubIDParam := c.Param("id")
	clubID, err := strconv.ParseUint(clubIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid club ID"})
		return
	}
	includeExpired := c.Query("include_expired") == "true"

	posts, err := h.postService.GetPollPostsByClubID(uint(clubID), includeExpired)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// @Summary Get user poll votes
// @Description Retrieve the poll votes made by the authenticated user for a specific poll post
// @Tags Posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]interface{} "User poll votes retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id}/poll/votes [get]
func (h *PostHandler) GetUserPollVotes(c *gin.Context) {
	postIDParam := c.Param("id")
	postID, err := strconv.ParseUint(postIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	useridRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, ok := useridRaw.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	votes, err := h.postService.GetUserPollVotes(uint(postID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"votes": votes})
}

// @Summary Remove vote from poll
// @Description Remove a user's vote from a poll post
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param vote body models.PollVoteRequest true "Vote data"
// @Success 200 {object} map[string]interface{} "Vote removed successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id}/unvote [post]
func (h *PostHandler) RemoveVoteFromPoll(c *gin.Context) {
	postIDParam := c.Param("id")
	postID, err := strconv.ParseUint(postIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	useridRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, ok := useridRaw.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	var req models.PollVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.postService.RemoveVoteFromPoll(uint(postID), userID, req.OptionIDs[0]); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vote removed successfully"})
}
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
// @Success 201 {object} models.Post "Post created successfully"
// @Failure 400 {object} gin.H "Bad request"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 500 {object} gin.H "Internal server error"
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
// @Success 200 {object} models.Post "Post retrieved successfully"
// @Failure 400 {object} gin.H "Bad request"
// @Failure 404 {object} gin.H "Post not found"
// @Failure 500 {object} gin.H "Internal server error"
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
// @Failure 400 {object} gin.H "Bad request"
// @Failure 404 {object} gin.H "Post not found"
// @Failure 500 {object} gin.H "Internal server error"
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
// @Failure 400 {object} gin.H "Bad request"
// @Failure 404 {object} gin.H "Post not found"
// @Failure 500 {object} gin.H "Internal server error"
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
// @Failure 400 {object} gin.H "Bad request"
// @Failure 500 {object} gin.H "Internal server error"
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
// @Failure 400 {object} gin.H "Bad request"
// @Failure 500 {object} gin.H "Internal server error"
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
// @Success 200 {array} models.Post "Posts retrieved successfully"
// @Failure 500 {object} gin.H "Internal server error"
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
// @Failure 500 {object} gin.H "Internal server error"
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
// @Failure 500 {object} gin.H "Internal server error"
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
// @Success 200 {object} gin.H "Post liked successfully"
// @Failure 400 {object} gin.H "Bad request"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "Post not found"
// @Failure 500 {object} gin.H "Internal server error"
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
// @Success 200 {object} gin.H "Post unliked successfully"
// @Failure 400 {object} gin.H "Bad request"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "Post not found"
// @Failure 500 {object} gin.H "Internal server error"
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
// @Success 200 {array} models.Like "Likes retrieved successfully"
// @Failure 400 {object} gin.H "Bad request"
// @Failure 404 {object} gin.H "Post not found or no likes found"
// @Failure 500 {object} gin.H "Internal server error"
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
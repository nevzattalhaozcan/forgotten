package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/internal/services"
)

type CommentHandler struct {
	CommentService *services.CommentService
	validator      *validator.Validate
}

func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{
		CommentService: commentService,
		validator:      validator.New(),
	}
}

func (c *CommentHandler) CreateComment(ctx *gin.Context) {
	uidRaw, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, ok := uidRaw.(uint)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	postidParam := ctx.Param("id")
	postID, err := strconv.ParseUint(postidParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	var req models.CreateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.validator.Struct(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := c.CommentService.CreateComment(uint(postID), userID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "comment created successfully",
		"comment": comment,
	})
}

func (c *CommentHandler) GetCommentByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
		return
	}

	comment, err := c.CommentService.GetCommentByID(uint(id))
	if err != nil {
		if err.Error() == "comment not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"comment": comment,
	})
}

func (c *CommentHandler) UpdateComment(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
		return
	}

	var req models.UpdateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.validator.Struct(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := c.CommentService.UpdateComment(uint(id), &req)
	if err != nil {
		if err.Error() == "comment not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "comment updated successfully",
		"comment": comment,
	})
}

func (c *CommentHandler) DeleteComment(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
		return
	}

	err = c.CommentService.DeleteComment(uint(id))
	if err != nil {
		if err.Error() == "comment not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *CommentHandler) ListCommentsByPostID(ctx *gin.Context) {
	postidParam := ctx.Param("id")
	postID, err := strconv.ParseUint(postidParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	comments, err := c.CommentService.ListCommentsByPostID(uint(postID))
	if err != nil {
		if err.Error() == "post not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}

func (c *CommentHandler) ListCommentsByUserID(ctx *gin.Context) {
	uidParam := ctx.Param("id")
	userID, err := strconv.ParseUint(uidParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	comments, err := c.CommentService.ListCommentsByUserID(uint(userID))
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}

func (c *CommentHandler) LikeComment(ctx *gin.Context) {
	uidRaw, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, ok := uidRaw.(uint)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	idParam := ctx.Param("id")
	commentID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
		return
	}

	err = c.CommentService.LikeComment(userID, uint(commentID))
	if err != nil {
        if err.Error() == "comment not found" {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
            return
        }
        if err.Error() == "user not found" {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
            return
        }
        if err.Error() == "user has already liked this comment" {
            ctx.JSON(http.StatusBadRequest, gin.H{"error": "user has already liked this comment"})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	ctx.JSON(http.StatusOK, gin.H{"message": "comment liked successfully"})
}

//TODO: fix unlike comment
func (c *CommentHandler) UnlikeComment(ctx *gin.Context) {
	uidRaw, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, ok := uidRaw.(uint)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	idParam := ctx.Param("id")
	commentID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
		return
	}

	err = c.CommentService.UnlikeComment(userID, uint(commentID))
	if err != nil {
        if err.Error() == "comment not found" {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
            return
        }
        if err.Error() == "user not found" {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
            return
        }
        if err.Error() == "user has not liked this comment" {
            ctx.JSON(http.StatusBadRequest, gin.H{"error": "user has not liked this comment"})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	ctx.JSON(http.StatusOK, gin.H{"message": "comment unliked successfully"})
}

func (c *CommentHandler) ListLikesByCommentID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	commentID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
		return
	}

	likes, err := c.CommentService.ListLikesByCommentID(uint(commentID))
	if err != nil {
		if err.Error() == "comment not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"likes": likes,
	})
}
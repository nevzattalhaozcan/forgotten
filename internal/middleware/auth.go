package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nevzattalhaozcan/forgotten/internal/config"
	"github.com/nevzattalhaozcan/forgotten/internal/repository"
	"github.com/nevzattalhaozcan/forgotten/pkg/utils"
	"gorm.io/gorm"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := utils.ValidateJWT(tokenString, cfg.JWT.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

/** Ensure that the user can only access their own resources
 * If the user ID in the path parameter does not match the user ID in the token, return 403 Forbidden
 * If there is no user ID in the path parameter, allow access (for routes that do not require a specific user ID)
 */
func AuthorizeSelf() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxUserIDRaw, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		ctxUserID, ok := ctxUserIDRaw.(uint)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID in context"})
			c.Abort()
			return
		}

		// Allow admin and superuser to access any resource
		userRoleRaw, exists := c.Get("user_role")
		if exists {
			if role, ok := userRoleRaw.(string); ok {
				if role == "admin" || role == "superuser" {
					c.Next()
					return
				}
			}
		}

		idParam := c.Param("id")
		if idParam == "" {
			c.Next()
			return
		}

		pathID, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			c.Abort()
			return
		}

		if uint(pathID) != ctxUserID {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}

/** Restrict access to certain roles
 * @param allowedRoles - list of roles that are allowed to access the route (admin user moderator support superuser)
 */
func RestrictToRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleRaw, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		userRole, ok := userRoleRaw.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user role in context"})
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		c.Abort()
	}
}

func RequireClubMembership(clubRepo repository.ClubRepository, postRepo repository.PostRepository, commentRepo repository.CommentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// allow admin and superuser to bypass club membership check
		userRoleRaw, exists := c.Get("user_role")
		if exists {
			if role, ok := userRoleRaw.(string); ok {
				if role == "admin" || role == "superuser" {
					c.Next()
					return
				}
			}
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		var clubID uint
		var err error

		path := c.FullPath()
		idParam := c.Param("id")

		if strings.Contains(path, "/posts/:id") && idParam != "" {
			// For post-related routes, get the post and extract club_id
			postID, parseErr := strconv.ParseUint(idParam, 10, 32)
			if parseErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
				c.Abort()
				return
			}

			post, err := postRepo.GetByID(uint(postID))
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch post"})
				}
				c.Abort()
				return
			}
			clubID = post.ClubID
		} else if strings.Contains(path, "/comments/:id") && idParam != "" {
			// For comment-related routes, get the comment, then post, then club_id
			commentID, parseErr := strconv.ParseUint(idParam, 10, 32)
			if parseErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
				c.Abort()
				return
			}

			comment, err := commentRepo.GetByID(uint(commentID))
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch comment"})
				}
				c.Abort()
				return
			}

			post, err := postRepo.GetByID(comment.PostID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch post"})
				}
				c.Abort()
				return
			}
			clubID = post.ClubID
		} else if strings.Contains(path, "/clubs/:id") && idParam != "" {
			// For direct club routes (e.g., /clubs/:id/posts)
			clubID64, parseErr := strconv.ParseUint(idParam, 10, 32)
			if parseErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid club ID"})
				c.Abort()
				return
			}
			clubID = uint(clubID64)
		} else {
			// Try to get club_id from request body
			bodyBytes, readErr := c.GetRawData()
			if readErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read request body"})
				c.Abort()
				return
			}

			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			var reqBody struct {
				ClubID uint `json:"club_id"`
			}

			if json.Unmarshal(bodyBytes, &reqBody) == nil && reqBody.ClubID > 0 {
				clubID = reqBody.ClubID
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "club_id is required"})
				c.Abort()
				return
			}
		}

		membership, err := clubRepo.GetClubMemberByUserID(uint(clubID), userID.(uint))
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{"error": "club membership required"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check membership"})
			}
			c.Abort()
			return
		}

		if !membership.IsApproved {
			c.JSON(http.StatusForbidden, gin.H{"error": "membership approval required"})
			c.Abort()
			return
		}

		c.Set("club_membership", membership)
		c.Next()
	}
}

func RequireClubMembershipWithRoles(clubRepo repository.ClubRepository, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// allow admin and superuser to bypass club membership check
		userRoleRaw, exists := c.Get("user_role")
		if exists {
			if role, ok := userRoleRaw.(string); ok {
				if role == "admin" || role == "superuser" {
					c.Next()
					return
				}
			}
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		clubIDParam := c.Param("id")
		if clubIDParam == "" {
			clubIDParam = c.Param("club_id")
		}

		clubID, err := strconv.ParseUint(clubIDParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid club ID"})
			c.Abort()
			return
		}

		membership, err := clubRepo.GetClubMemberByUserID(uint(clubID), userID.(uint))
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{"error": "club membership required"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check membership"})
			}
			c.Abort()
			return
		}

		if !membership.IsApproved {
			c.JSON(http.StatusForbidden, gin.H{"error": "membership approval required"})
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if membership.Role == role {
				c.Set("club_membership", membership)
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient club membership role"})
		c.Abort()
	}
}

func OptionalAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := utils.ValidateJWT(tokenString, cfg.JWT.Secret)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

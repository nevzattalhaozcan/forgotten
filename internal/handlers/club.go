package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/internal/services"
	"gorm.io/gorm"
)

type ClubHandler struct {
	clubService *services.ClubService
	validator   *validator.Validate
}

func NewClubHandler(clubService *services.ClubService) *ClubHandler {
	return &ClubHandler{
		clubService: clubService,
		validator:   validator.New(),
	}
}

// @Summary Create a new club
// @Description Create a new book club
// @Tags Clubs
// @Accept json
// @Produce json
// @Param request body models.CreateClubRequest true "Club creation data"
// @Success 201 {object} map[string]interface{} "Club created successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Router /api/v1/clubs [post]
func (h *ClubHandler) CreateClub(c *gin.Context) {
	var req models.CreateClubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uidRaw, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
	}

	ownerID, ok := uidRaw.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	club, err := h.clubService.CreateClub(ownerID, &req)
	if err != nil {
		if errors.Is(err, services.ErrClubNameExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "club created successfully",
		"club":    club,
	})
}

// @Summary Get club by ID
// @Description Retrieve a club by its ID
// @Tags Clubs
// @Produce json
// @Param id path int true "Club ID"
// @Success 200 {object} map[string]interface{} "Club retrieved successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Club not found"
// @Router /api/v1/clubs/{id} [get]
func (h *ClubHandler) GetClub(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid club ID"})
		return
	}

	club, err := h.clubService.GetClubByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "club not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"club": club})
}

// @Summary Get all clubs
// @Description Retrieve a list of all clubs with optional filters
// @Tags Clubs
// @Produce json
// @Param location query string false "Filter by location (partial match)"
// @Param genre query string false "Filter by genre (partial match)"
// @Param meeting_type query string false "Filter by meeting type" Enums(online, in-person, hybrid)
// @Param min_members query int false "Minimum member count"
// @Param max_members query int false "Maximum member count"
// @Param limit query int false "Number of results to return" default(20)
// @Param offset query int false "Number of results to skip" default(0)
// @Success 200 {object} map[string]interface{} "Clubs retrieved successfully"
// @Failure 400 {object} map[string]string "Bad request - invalid filter parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/clubs [get]
func (h *ClubHandler) GetAllClubs(c *gin.Context) {
	location := c.Query("location")
	genre := c.Query("genre")
	meetingType := c.Query("meeting_type")
	minMembers, _ := strconv.Atoi(c.DefaultQuery("min_members", "0"))
	maxMembers, _ := strconv.Atoi(c.DefaultQuery("max_members", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	hasFilters := location != "" || genre != "" || meetingType != "" || minMembers > 0 || maxMembers > 0

	if !hasFilters {
		clubs, err := h.clubService.GetAllClubs()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve clubs"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"clubs": clubs})
		return
	}

	if meetingType != "" && meetingType != "online" && meetingType != "in-person" && meetingType != "hybrid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "meeting_type must be one of: online, in-person, hybrid"})
		return
	}

	clubs, err := h.clubService.GetClubsWithFilters(location, genre, meetingType, minMembers, maxMembers, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve clubs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"clubs":  clubs,
		"count":  len(clubs),
		"limit":  limit,
		"offset": offset,
	})
}

// @Summary Update club
// @Description Update club information
// @Tags Clubs
// @Accept json
// @Produce json
// @Param id path int true "Club ID"
// @Param request body models.UpdateClubRequest true "Club update data"
// @Success 200 {object} map[string]interface{} "Club updated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Club not found"
// @Router /api/v1/clubs/{id} [put]
func (h *ClubHandler) UpdateClub(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid club ID"})
		return
	}

	var req models.UpdateClubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.clubService.GetClubByID(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "club not found"})
		return
	}

	uidRaw, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userID, ok := uidRaw.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}
	if !h.clubService.CanManageClub(uint(id), userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	club, err := h.clubService.UpdateClub(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "club updated successfully",
		"club":    club,
	})
}

// @Summary Delete club
// @Description Delete a club by its ID
// @Tags Clubs
// @Param id path int true "Club ID"
// @Success 204 {object} nil "Club deleted successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/clubs/{id} [delete]
func (h *ClubHandler) DeleteClub(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid club ID"})
		return
	}

	if _, err := h.clubService.GetClubByID(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "club not found"})
		return
	}

	uidRaw, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userID, ok := uidRaw.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}
	if !h.clubService.CanManageClub(uint(id), userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	err = h.clubService.DeleteClub(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Join a club
// @Description Join a club by its ID
// @Tags Clubs
// @Param id path int true "Club ID"
// @Success 200 {object} map[string]interface{} "Joined club successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/clubs/{id}/join [post]
func (h *ClubHandler) JoinClub(c *gin.Context) {
	uidRaw, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userID, ok := uidRaw.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	clubIDParam := c.Param("id")
	clubID, err := strconv.ParseUint(clubIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid club ID"})
		return
	}

	membership, err := h.clubService.JoinClub(uint(clubID), uint(userID))
	if err != nil {
		if err.Error() == "user is already a member of the club" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "joined successfully",
		"member":  membership,
	})
}

// @Summary Leave a club
// @Description Leave a club by its ID
// @Tags Clubs
// @Param id path int true "Club ID"
// @Success 200 {object} map[string]interface{} "Left club successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/clubs/{id}/leave [post]
func (h *ClubHandler) LeaveClub(c *gin.Context) {
	clubIDParam := c.Param("id")
	clubID, err := strconv.ParseUint(clubIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid club ID"})
		return
	}

	userIDParam, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, ok := userIDParam.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	var req models.OwnerLeaveRequest
	var ownerAction *models.OwnerLeaveRequest
	if err := c.ShouldBindJSON(&req); err == nil {
		ownerAction = &req
	}

	if err := h.clubService.LeaveClub(uint(clubID), userID, ownerAction); err != nil {
		switch err.Error() {
		case "club not found", "not a member":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		case "owner must choose action: transfer or close",
			"new_owner_id is required for transfer",
			"only the owner can transfer ownership",
			"new owner must be different from current owner",
			"new owner must be a member of the club",
			"new owner must be an approved member",
			"invalid action; must be one of: transfer, close":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "left club successfully",
	})
}

// @Summary List club members
// @Description List all members of a club by its ID
// @Tags Clubs
// @Produce json
// @Param id path int true "Club ID"
// @Success 200 {object} map[string]interface{} "List of club members"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/clubs/{id}/members [get]
func (h *ClubHandler) ListClubMembers(c *gin.Context) {
	clubIDParam := c.Param("id")
	clubID, err := strconv.ParseUint(clubIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid club ID"})
		return
	}

	members, err := h.clubService.ListClubMembers(uint(clubID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

// @Summary Update club member
// @Description Update a club member's information
// @Tags Clubs
// @Accept json
// @Produce json
// @Param id path int true "Club ID"
// @Param user_id path int true "User ID"
// @Param request body models.UpdateClubMembershipRequest true "Member update data"
// @Success 200 {object} map[string]interface{} "Member updated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Club not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/clubs/{id}/members/{user_id} [put]
func (h *ClubHandler) UpdateClubMember(c *gin.Context) {
	clubIDParam := c.Param("id")
	clubID, err := strconv.ParseUint(clubIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid club ID"})
		return
	}

	userIDParam := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req models.UpdateClubMembershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.clubService.GetClubByID(uint(clubID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "club not found"})
		return
	}

	uidRaw, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	requesterID, ok := uidRaw.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}
	if !h.clubService.CanManageClub(uint(clubID), requesterID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	updated, err := h.clubService.UpdateClubMemberFields(uint(clubID), uint(userID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "member updated successfully",
		"member":  updated,
	})
}

// @Summary Get club member by user ID
// @Description Retrieve a club member's information by user ID
// @Tags Clubs
// @Produce json
// @Param id path int true "Club ID"
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]interface{} "Member retrieved successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Member not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/clubs/{id}/members/{user_id} [get]
func (h *ClubHandler) GetClubMember(c *gin.Context) {
	clubIDParam := c.Param("id")
	clubID, err := strconv.ParseUint(clubIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid club ID"})
		return
	}

	userIDParam := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	member, err := h.clubService.GetClubMemberByUserID(uint(clubID), uint(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"member": member})
}

// @Summary Rate a club
// @Description Rate a club by its ID
// @Tags Clubs
// @Accept json
// @Produce json
// @Param id path int true "Club ID"
// @Param request body models.RateClubRequest true "Club rating data"
// @Success 200 {object} map[string]interface{} "Club rated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/clubs/{id}/rate [post]
func (h *ClubHandler) RateClub(c *gin.Context) {
	clubIDParam := c.Param("id")
	clubID64, err := strconv.ParseUint(clubIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid club ID"})
		return
	}
	clubID := uint(clubID64)

	userIDRaw, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userID, ok := userIDRaw.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	var req models.RateClubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.clubService.RateClub(userID, clubID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "club rated successfully",
		"rating":  resp,
	})
}

// @Summary List club ratings
// @Description List all ratings for a club by its ID
// @Tags Clubs
// @Produce json
// @Param id path int true "Club ID"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{} "List of club ratings"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/clubs/{id}/ratings [get]
func (h *ClubHandler) ListClubRatings(c *gin.Context) {
	clubIDParam := c.Param("id")
	clubID64, err := strconv.ParseUint(clubIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid club ID"})
		return
	}
	clubID := uint(clubID64)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	ratings, err := h.clubService.ListClubRatings(clubID, limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ratings": ratings,
	})
}

// @Summary Get user's clubs
// @Description Retrieve a list of clubs the authenticated user is a member of
// @Tags Users
// @Produce json
// @Success 200 {object} map[string]interface{} "List of user's clubs"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/users/my-clubs [get]
func (h *ClubHandler) GetMyClubs(c *gin.Context) {
	uidRaw, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userID, ok := uidRaw.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	clubs, err := h.clubService.ListUserClubs(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"clubs": clubs})
}

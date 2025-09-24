package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "github.com/nevzattalhaozcan/forgotten/internal/models"
    "github.com/nevzattalhaozcan/forgotten/internal/services"
)

type ReadingHandler struct {
    readingService *services.ReadingService
    validator      *validator.Validate
}

func NewReadingHandler(readingService *services.ReadingService) *ReadingHandler {
    return &ReadingHandler{
		readingService: readingService,
		validator:      validator.New(),
	}
}

// @Summary Start Reading a Book
// @Description Start reading a book by providing the book ID.
// @Tags Reading
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body models.StartReadingRequest true "Start Reading Request"
// @Success 201 {object} models.UserBookProgressResponse
// @Failure 400 {object} gin.H{"error": "error message"}
// @Router /users/{id}/readings [post]
// @Security Bearer
func (h *ReadingHandler) StartReading(c *gin.Context) {
    userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

    var req models.StartReadingRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"}); return
    }
    if err := h.validator.Struct(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
    }

    resp, err := h.readingService.StartReading(uint(userID), req.BookID)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    c.JSON(http.StatusCreated, resp)
}

// @Summary Update Reading Progress
// @Description Update the reading progress of a book.
// @Tags Reading
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param bookID path int true "Book ID"
// @Param request body models.UpdateReadingProgressRequest true "Update Reading Progress Request"
// @Success 200 {object} models.UserBookProgressResponse
// @Failure 400 {object} gin.H{"error": "error message"}
// @Router /users/{id}/readings/{bookID} [put]
// @Security Bearer
func (h *ReadingHandler) UpdateProgress(c *gin.Context) {
    userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    bookID, _ := strconv.ParseUint(c.Param("bookID"), 10, 64)

    var req models.UpdateReadingProgressRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"}); return
    }
    if err := h.validator.Struct(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
    }

    resp, err := h.readingService.UpdateProgress(uint(userID), uint(bookID), &req)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    c.JSON(http.StatusOK, resp)
}

// @Summary Complete Reading a Book
// @Description Mark a book as completed and optionally add a note.
// @Tags Reading
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param bookID path int true "Book ID"
// @Param request body models.CompleteReadingRequest true "Complete Reading Request"
// @Success 200 {object} models.UserBookProgressResponse
// @Failure 400 {object} gin.H{"error": "error message"}
// @Router /users/{id}/readings/{bookID}/complete [post]
// @Security Bearer
func (h *ReadingHandler) CompleteReading(c *gin.Context) {
    userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    bookID, _ := strconv.ParseUint(c.Param("bookID"), 10, 64)

    var req models.CompleteReadingRequest
    _ = c.ShouldBindJSON(&req)

    resp, err := h.readingService.CompleteReading(uint(userID), uint(bookID), req.Note)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    c.JSON(http.StatusOK, resp)
}

// @Summary List User Reading Progress
// @Description List all reading progress entries for a user.
// @Tags Reading
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} models.UserBookProgressResponse
// @Failure 400 {object} gin.H{"error": "error message"}
// @Router /users/{id}/reading [get]
// @Security Bearer
func (h *ReadingHandler) ListUserProgress(c *gin.Context) {
    userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    resp, err := h.readingService.ListUserProgress(uint(userID))
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    c.JSON(http.StatusOK, resp)
}

// @Summary Get User Reading History
// @Description Retrieve the reading history of a user.
// @Tags Reading
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} models.ReadingLogResponse
// @Failure 400 {object} gin.H{"error": "error message"}
// @Router /users/{id}/reading/history [get]
// @Security Bearer
func (h *ReadingHandler) UserReadingHistory(c *gin.Context) {
    userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    resp, err := h.readingService.UserReadingHistory(uint(userID))
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    c.JSON(http.StatusOK, resp)
}

// @Summary Assign Book to Club
// @Description Assign a book to a club for reading.
// @Tags Reading
// @Accept json
// @Produce json
// @Param id path int true "Club ID"
// @Param request body models.AssignBookRequest true "Assign Book Request"
// @Success 201 {object} models.ClubBookAssignmentResponse
// @Failure 400 {object} gin.H{"error": "error message"}
// @Router /clubs/{id}/reading/assign [post]
// @Security Bearer
func (h *ReadingHandler) AssignBookToClub(c *gin.Context) {
    clubID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

    var req models.AssignBookRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"}); return
    }
    if err := h.validator.Struct(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
    }

    resp, err := h.readingService.AssignBookToClub(uint(clubID), req.BookID, &req)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    c.JSON(http.StatusCreated, resp)
}

// @Summary Update Club Reading Checkpoint
// @Description Update the reading checkpoint for a club's current book assignment.
// @Tags Reading
// @Accept json
// @Produce json
// @Param id path int true "Club ID"
// @Param request body models.UpdateClubCheckpointRequest true "Update Club Checkpoint Request"
// @Success 200 {object} models.ClubBookAssignmentResponse
// @Failure 400 {object} gin.H{"error": "error message"}
// @Router /clubs/{id}/reading/checkpoint [patch]
// @Security Bearer
func (h *ReadingHandler) UpdateClubCheckpoint(c *gin.Context) {
    clubID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    var req models.UpdateClubCheckpointRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"}); return
    }
    resp, err := h.readingService.UpdateClubCheckpoint(uint(clubID), &req)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    c.JSON(http.StatusOK, resp)
}

// @Summary Complete Club Book Assignment
// @Description Mark the current book assignment for a club as completed.
// @Tags Reading
// @Produce json
// @Param id path int true "Club ID"
// @Success 200 {object} models.ClubBookAssignmentResponse
// @Failure 400 {object} gin.H{"error": "error message"}
// @Router /clubs/{id}/reading/complete [post]
// @Security Bearer
func (h *ReadingHandler) CompleteClubAssignment(c *gin.Context) {
    clubID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    resp, err := h.readingService.CompleteClubAssignment(uint(clubID))
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    c.JSON(http.StatusOK, resp)
}

// @Summary List Club Book Assignments
// @Description List all book assignments for a club.
// @Tags Reading
// @Produce json
// @Param id path int true "Club ID"
// @Success 200 {array} models.ClubBookAssignmentResponse
// @Failure 400 {object} gin.H{"error": "error message"}
// @Router /clubs/{id}/reading [get]
// @Security Bearer
func (h *ReadingHandler) ListClubAssignments(c *gin.Context) {
    clubID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    resp, err := h.readingService.ListClubAssignments(uint(clubID))
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    c.JSON(http.StatusOK, resp)
}

// @Summary Sync User Reading Stats
// @Description Synchronize the reading statistics for a user.
// @Tags Reading
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string{"message": "user stats synchronized successfully"}
// @Failure 400 {object} gin.H{"error": "error message"}
// @Router /users/{id}/reading/sync [post]
// @Security Bearer
func (h *ReadingHandler) SyncUserStats(c *gin.Context) {
    userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    if err := h.readingService.SyncUserStats(uint(userID)); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
    }
    c.JSON(http.StatusOK, gin.H{"message": "user stats synchronized successfully"})
}
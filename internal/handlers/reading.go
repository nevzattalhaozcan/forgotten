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

func (h *ReadingHandler) CompleteReading(c *gin.Context) {
    userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    bookID, _ := strconv.ParseUint(c.Param("bookID"), 10, 64)

    var req models.CompleteReadingRequest
    _ = c.ShouldBindJSON(&req)

    resp, err := h.readingService.CompleteReading(uint(userID), uint(bookID), req.Note)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    c.JSON(http.StatusOK, resp)
}

func (h *ReadingHandler) ListUserProgress(c *gin.Context) {
    userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    resp, err := h.readingService.ListUserProgress(uint(userID))
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    c.JSON(http.StatusOK, resp)
}

func (h *ReadingHandler) UserReadingHistory(c *gin.Context) {
    userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    resp, err := h.readingService.UserReadingHistory(uint(userID))
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    c.JSON(http.StatusOK, resp)
}

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

func (h *ReadingHandler) CompleteClubAssignment(c *gin.Context) {
    clubID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    resp, err := h.readingService.CompleteClubAssignment(uint(clubID))
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    c.JSON(http.StatusOK, resp)
}

func (h *ReadingHandler) ListClubAssignments(c *gin.Context) {
    clubID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
    resp, err := h.readingService.ListClubAssignments(uint(clubID))
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
    c.JSON(http.StatusOK, resp)
}
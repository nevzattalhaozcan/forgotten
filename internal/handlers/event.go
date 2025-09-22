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

type EventHandler struct {
	eventService *services.EventService
	validator    *validator.Validate
}

func NewEventHandler(eventService *services.EventService) *EventHandler {
	return &EventHandler{
		eventService: eventService,
		validator:    validator.New(),
	}
}

func (h *EventHandler) CreateEvent(c *gin.Context) {
	var req models.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		if strings.Contains(err.Error(), "gtfield") {
            c.JSON(http.StatusBadRequest, gin.H{"error": "end_time must be after start_time"})
            return
        }
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clubIDParam := c.Param("id")
	clubID, err := strconv.ParseUint(clubIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid club_id"})
		return
	}

	event, err := h.eventService.CreateEvent(uint(clubID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "event created successfully",
		"event":   event,
	})
}

func (h *EventHandler) GetClubEvents(c *gin.Context) {
	clubIDParam := c.Param("id")
	clubID, err := strconv.ParseUint(clubIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid club_id"})
		return
	}

	events, err := h.eventService.GetClubEvents(uint(clubID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

func (h *EventHandler) GetEvent(c *gin.Context) {
	eventIDParam := c.Param("id")
	eventID, err := strconv.ParseUint(eventIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event_id"})
		return
	}

	event, err := h.eventService.GetEventByID(uint(eventID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event})
}

func (h *EventHandler) UpdateEvent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event_id"})
		return
	}

	var req models.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := h.eventService.UpdateEvent(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "event updated successfully",
		"event":   event,
	})
}

func (h *EventHandler) DeleteEvent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event_id"})
		return
	}

	if err := h.eventService.DeleteEvent(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *EventHandler) RSVPToEvent(c *gin.Context) {
	eventParam := c.Param("id")
	eventID, err := strconv.ParseUint(eventParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event_id"})
		return
	}

	var req models.RSVPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdRaw, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found in context"})
		return
	}
	userID, ok := userIdRaw.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id in context"})
		return
	}

	err = h.eventService.RSVPToEvent(uint(eventID), &models.EventRSVP{
		UserID:  userID,
		EventID: uint(eventID),
		Status:  req.Status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "RSVP successful"})
}

func (h *EventHandler) GetEventAttendees(c *gin.Context) {
	eventParam := c.Param("id")
	eventID, err := strconv.ParseUint(eventParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event_id"})
		return
	}

	attendees, err := h.eventService.GetEventAttendees(uint(eventID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"attendees": attendees})
}
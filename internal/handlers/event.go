package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/internal/services"
	"gorm.io/gorm"
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

// @Summary Create a new event
// @Description Create a new event for a specific club
// @Tags Events
// @Accept json
// @Produce json
// @Param id path int true "Club ID"
// @Param request body models.CreateEventRequest true "Event data"
// @Success 201 {object} map[string]interface{} "Event created successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/clubs/{id}/events [post]
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

// @Summary List events for a club
// @Description Retrieve all events for a specific club
// @Tags Events
// @Produce json
// @Param id path int true "Club ID"
// @Success 200 {object} map[string]interface{} "List of events"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/clubs/{id}/events [get]
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

// @Summary Get event details
// @Description Retrieve details of a specific event by its ID
// @Tags Events
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} map[string]interface{} "Event details"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/events/{id} [get]
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

// @Summary Update an event
// @Description Update event information
// @Tags Events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param request body models.UpdateEventRequest true "Update event data"
// @Success 200 {object} map[string]interface{} "Event updated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/events/{id} [put]
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

// @Summary Delete an event
// @Description Delete an event by its ID
// @Tags Events
// @Param id path int true "Event ID"
// @Success 204 {object} nil "No Content"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/events/{id} [delete]
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

// @Summary RSVP to an event
// @Description RSVP to an event by its ID
// @Tags Events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param request body models.RSVPRequest true "RSVP data"
// @Success 200 {object} map[string]interface{} "RSVP successful"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/events/{id}/rsvp [post]
func (h *EventHandler) RSVPToEvent(c *gin.Context) {
	eventParam := c.Param("id")
	eventID, err := strconv.ParseUint(eventParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event_id"})
		return
	}

	event, err := h.eventService.GetEventByID(uint(eventID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
        return
    }

    clubID := event.ClubID

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

	membership, err := h.eventService.ClubRepo().GetClubMemberByUserID(clubID, userID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusForbidden, gin.H{"error": "club membership required"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check membership"})
        }
        return
    }
    if !membership.IsApproved {
        c.JSON(http.StatusForbidden, gin.H{"error": "membership approval required"})
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

// @Summary Get event attendees
// @Description Retrieve a list of attendees for a specific event
// @Tags Events
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} map[string]interface{} "List of attendees"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/events/{id}/attendees [get]
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

// @Summary Get public events
// @Description Retrieve a list of all public events
// @Tags Events
// @Produce json
// @Success 200 {object} map[string]interface{} "List of public events"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/events/public [get]
func (h *EventHandler) GetPublicEvents(c *gin.Context) {
	events, err := h.eventService.GetPublicEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}
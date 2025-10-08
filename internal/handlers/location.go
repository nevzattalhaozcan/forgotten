package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nevzattalhaozcan/forgotten/internal/cache"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/internal/services"
)

type LocationHandler struct {
	locationService *services.LocationService
	locationCache   *cache.LocationCache
	validator       *validator.Validate
}

func NewLocationHandler(locationService *services.LocationService, locationCache *cache.LocationCache) *LocationHandler {
	return &LocationHandler{
		locationService: locationService,
		locationCache:   locationCache,
		validator:       validator.New(),
	}
}

// @Summary Search locations
// @Description Search for Turkish cities and districts with autocomplete
// @Tags Locations
// @Produce json
// @Param q query string true "Search query (minimum 1 character)"
// @Param type query string false "Search type" Enums(city, district, all) default(all)
// @Param limit query int false "Maximum results to return" default(10)
// @Success 200 {object} map[string]interface{} "Search results"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/locations/search [get]
func (h *LocationHandler) SearchLocations(c *gin.Context) {
	var req models.LocationSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid query parameters"})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	var results []models.LocationSearchResponse
    var err error

	if h.locationCache != nil {
        results, err = h.locationCache.SearchLocations(context.Background(), req.Query, req.Type, req.Limit)
    } else {
        results = h.locationService.SearchLocations(req.Query, req.Type, req.Limit)
    }

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search locations"})
		return
	}

	c.Header("Cache-Control", "public, max-age=3600")
	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"count":   len(results),
		"query":   req.Query,
	})
}
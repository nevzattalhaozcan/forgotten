package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/internal/services"
	"github.com/redis/go-redis/v9"
)

type LocationCache struct {
	client          *redis.Client
	locationService *services.LocationService
}

func NewLocationCache(redisURL string, locationService *services.LocationService) *LocationCache {
	opts, _ := redis.ParseURL(redisURL)
	client := redis.NewClient(opts)
	
	return &LocationCache{
		client: client,
		locationService: locationService,
	}
}

func (c *LocationCache) SearchLocations(ctx context.Context, query string, searchType string, limit int) ([]models.LocationSearchResponse, error) {
	cacheKey := fmt.Sprintf("location_search:%s:%s:%d", query, searchType, limit)

	cached, err := c.client.Get(ctx, cacheKey).Result()
	if err == nil {
		var results []models.LocationSearchResponse
		if json.Unmarshal([]byte(cached), &results) == nil {
			return results, nil
		}
	}

	results := c.locationService.SearchLocations(query, searchType, limit)

	resultJSON, _ := json.Marshal(results)
	c.client.Set(ctx, cacheKey, resultJSON, time.Hour)

	return results, nil
}

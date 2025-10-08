package services

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/nevzattalhaozcan/forgotten/internal/models"
)

type LocationService struct {
	cities    []models.City
	districts []models.District
	cityMap   map[string]models.City
}

func NewLocationService() (*LocationService, error) {
	service := &LocationService{
		cityMap: make(map[string]models.City),
	}

	if err := service.loadData(); err != nil {
		return nil, err
	}
	return service, nil
}

func (s *LocationService) loadData() error {
	citiesData, err := os.ReadFile("data/sehirler.json")
	if err != nil {
		return fmt.Errorf("failed to read cities file: %w", err)
	}

	var rawCities []struct {
		SehirID  string `json:"sehir_id"`
		SehirAdi string `json:"sehir_adi"`
	}

	if err := json.Unmarshal(citiesData, &rawCities); err != nil {
		return fmt.Errorf("failed to unmarshal cities data: %w", err)
	}

	for _, city := range rawCities {
		c := models.City{
			ID:   city.SehirID,
			Name: city.SehirAdi,
		}
		s.cities = append(s.cities, c)
		s.cityMap[city.SehirID] = c
	}

	districtsData, err := os.ReadFile("data/ilceler.json")
	if err != nil {
		return fmt.Errorf("failed to read districts file: %w", err)
	}

	var rawDistricts []struct {
		IlceID   string `json:"ilce_id"`
		IlceAdi  string `json:"ilce_adi"`
		SehirID  string `json:"sehir_id"`
		SehirAdi string `json:"sehir_adi"`
	}

	if err := json.Unmarshal(districtsData, &rawDistricts); err != nil {
		return fmt.Errorf("failed to unmarshal districts data: %w", err)
	}

	for _, district := range rawDistricts {
		s.districts = append(s.districts, models.District{
			ID:     district.IlceID,
			Name:   district.IlceAdi,
			CityID: district.SehirID,
			City:   district.SehirAdi,
		})
	}

	return nil
}

func normalizeText(text string) string {
	text = strings.ToLower(text)

	replacer := strings.NewReplacer(
		"ç", "c", "ğ", "g", "ı", "i", "ö", "o", "ş", "s", "ü", "u",
        "Ç", "c", "Ğ", "g", "İ", "i", "Ö", "o", "Ş", "s", "Ü", "u",
	)

	return replacer.Replace(text)
}

func (s *LocationService) SearchLocations(query string, searchType string, limit int) []models.LocationSearchResponse {
	if limit <= 0 {
		limit = 10
	}

	normalizedQuery := normalizeText(query)
	var results []models.LocationSearchResponse

	if searchType == "city" || searchType == "all" || searchType == "" {
		for _, city := range s.cities {
			normalizedName := normalizeText(city.Name)
			if strings.Contains(normalizedName, normalizedQuery) {
				results = append(results, models.LocationSearchResponse{
					ID:  city.ID,
					Name: city.Name,
					Type: "city",
				})
			}
		}
	}

	if searchType == "district" || searchType == "all" || searchType == "" {
		for _, district := range s.districts {
			normalizedName := normalizeText(district.Name)
			if strings.Contains(normalizedName, normalizedQuery) {
				results = append(results, models.LocationSearchResponse{
					ID:   district.ID,
					Name: district.Name,
					Type: "district",
					CityName: district.City,
					CityID: district.CityID,
				})
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		iNorm := normalizeText(results[i].Name)
		jNorm := normalizeText(results[j].Name)
		
		iExact := strings.HasPrefix(iNorm, normalizedQuery)
		jExact := strings.HasPrefix(jNorm, normalizedQuery)

		if iExact && !jExact {
			return true
		}

		if !iExact && jExact {
			return false
		}

		return len(results[i].Name) < len(results[j].Name)
	})

	if len(results) > limit {
		results = results[:limit]
	}

	return results
}
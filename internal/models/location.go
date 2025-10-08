package models

type City struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type District struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	CityID string `json:"city_id"`
	City   string `json:"city"`
}

type LocationSearchRequest struct {
	Query string `form:"q" validate:"required,min=1"`
	Type  string `form:"type" validate:"omitempty,oneof=city district all"`
	Limit int    `form:"limit" validate:"omitempty,gte=1,lte=50"`
}

type LocationSearchResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	CityName string `json:"city_name,omitempty"`
	CityID   string `json:"city_id,omitempty"`
}

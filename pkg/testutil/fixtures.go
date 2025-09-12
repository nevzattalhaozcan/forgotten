package test_helpers

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nevzattalhaozcan/forgotten/internal/database"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
)

func CreateUser() (*models.User, error) {
	data, err := os.ReadFile("../../data/seed_users.json")
	if err != nil {
		log.Printf("failed to read seed file: %v", err)
		return nil, err
	}

	var seedData database.SeedData
	if err := json.Unmarshal(data, &seedData); err != nil {
		log.Printf("failed to unmarshal seed data: %v", err)
		return nil, err
	}

	seedUser := seedData.Users[0]
	user := &models.User{
		Username:     seedUser.Username,
		Email:        seedUser.Email,
		PasswordHash: seedUser.Password,
		FirstName:    seedUser.FirstName,
		LastName:     seedUser.LastName,
		IsActive:     seedUser.IsActive,
	}

	return user, nil
}

func CreateRegisterRequest() (*models.RegisterRequest, error) {
	data, err := os.ReadFile("../../data/seed_users.json")
	if err != nil {
		log.Printf("failed to read seed file: %v", err)
		return nil, err
	}

	var seedData database.SeedData
	if err := json.Unmarshal(data, &seedData); err != nil {
		log.Printf("failed to unmarshal seed data: %v", err)
		return nil, err
	}

	seedUser := seedData.Users[0]
	registerReq := &models.RegisterRequest{
		Username:     seedUser.Username,
		Email:        seedUser.Email,
		Password:     seedUser.Password,
		FirstName:    seedUser.FirstName,
		LastName:     seedUser.LastName,
	}

	return registerReq, nil
}

func CreateLoginRequest() (*models.LoginRequest, error) {
	data, err := os.ReadFile("../../data/seed_users.json")
	if err != nil {
		log.Printf("failed to read seed file: %v", err)
		return nil, err
	}

	var seedData database.SeedData
	if err := json.Unmarshal(data, &seedData); err != nil {
		log.Printf("failed to unmarshal seed data: %v", err)
		return nil, err
	}

	seedUser := seedData.Users[0]
	loginReq := &models.LoginRequest{
		Email:    seedUser.Email,
		Password: seedUser.Password,
	}

	return loginReq, nil
}
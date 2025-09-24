package database

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/pkg/utils"
	"gorm.io/gorm"
)

type SeedUser struct {
	Username       string   `json:"username"`
	Email          string   `json:"email"`
	Password       string   `json:"password"`
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	IsActive       bool     `json:"is_active"`
	Role           string   `json:"role"`
	AvatarURL      string   `json:"avatar_url"`
	Bio            string   `json:"bio"`
	Location       string   `json:"location"`
	FavoriteGenres []string `json:"favorite_genre"`
	ReadingGoal    int      `json:"reading_goal"`
}

type SeedData struct {
	Users []SeedUser `json:"users"`
}

func Seed(db *gorm.DB) error {
	log.Println("starting database seeding")

	if err := truncateAll(db); err != nil {
        log.Printf("truncate error: %v", err)
        return err
    }

	data, err := os.ReadFile("data/seed_users.json")
	if err != nil {
		log.Printf("failed to read seed file: %v", err)
		return err
	}

	var seedData SeedData
	if err := json.Unmarshal(data, &seedData); err != nil {
		log.Printf("failed to unmarshal seed data: %v", err)
		return err
	}

	for _, seedUser := range seedData.Users {
		hashedPassword, err := utils.HashPassword(seedUser.Password)
		if err != nil {
			log.Printf("failed to hash password for user %s: %v", seedUser.Username, err)
			continue
		}

		user := &models.User{
			Username:       seedUser.Username,
			Email:          seedUser.Email,
			PasswordHash:   hashedPassword,
			FirstName:      seedUser.FirstName,
			LastName:       seedUser.LastName,
			IsActive:       seedUser.IsActive,
			Role:           seedUser.Role,
			AvatarURL:      &seedUser.AvatarURL,
			Bio:            &seedUser.Bio,
			Location:       &seedUser.Location,
			FavoriteGenres: seedUser.FavoriteGenres,
			ReadingGoal:    seedUser.ReadingGoal,
		}

		if err := db.Create(user).Error; err != nil {
			log.Printf("failed to create user %s: %v", seedUser.Username, err)
			return err
		}
		log.Printf("created user: %s", seedUser.Username)
	}
	log.Printf("database seeding completed - %d users created", len(seedData.Users))
	return nil
}

func SeedForTest(db *gorm.DB) error {
	log.Println("starting test database seeding")

	if err := truncateAll(db); err != nil {
        log.Printf("truncate error: %v", err)
        return err
    }

	if err := db.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("could not reset user ID sequence: %v", err)
	}

	return Seed(db)
}

func truncateAll(db *gorm.DB) error {
    return db.Exec(`
        TRUNCATE TABLE 
            post_likes,
            comment_likes,
            comments,
            posts,
            event_rsvps,
            events,
            club_moderators,
            club_memberships,
            clubs,
            books,
            users
        RESTART IDENTITY CASCADE
    `).Error
}
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/redis/go-redis/v9"
)

type cachedUserRepository struct {
	base UserRepository
	rdb  *redis.Client
	ttl  time.Duration
}

func NewCachedUserRepository(base UserRepository, rdb *redis.Client, ttl time.Duration) UserRepository {
	return &cachedUserRepository{base: base, rdb: rdb, ttl: ttl}
}

func (r *cachedUserRepository) keyByID(id uint) string { return fmt.Sprintf("user:id:%d", id) }
func (r *cachedUserRepository) keyByEmail(email string) string {
	return fmt.Sprintf("user:email:%s", email)
}
func (r *cachedUserRepository) keyByUsername(username string) string {
	return fmt.Sprintf("user:username:%s", username)
}

type userCacheEntry struct {
	ID             uint           `json:"id"`
	Username       string         `json:"username"`
	Email          string         `json:"email"`
	PasswordHash   string         `json:"password_hash"`
	FirstName      string         `json:"first_name"`
	LastName       string         `json:"last_name"`
	IsActive       bool           `json:"is_active"`
	Role           string         `json:"role"`
	AvatarURL      *string        `json:"avatar_url"`
	Location       *string        `json:"location"`
	FavoriteGenres pq.StringArray `json:"favorite_genres"`
	Bio            *string        `json:"bio"`
	ReadingGoal    *int           `json:"reading_goal"`
	BooksRead      int            `json:"books_read"`
	Badges         pq.StringArray `json:"badges"`
	IsOnline       bool           `json:"is_online"`
	LastSeen       *time.Time     `json:"last_seen"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func toCache(user *models.User) *userCacheEntry {
	return &userCacheEntry{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		PasswordHash:   user.PasswordHash,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		IsActive:       user.IsActive,
		Role:           user.Role,
		AvatarURL:      user.AvatarURL,
		Location:       user.Location,
		FavoriteGenres: user.FavoriteGenres,
		Bio:            user.Bio,
		ReadingGoal:    &user.ReadingGoal,
		BooksRead:      user.BooksRead,
		Badges:         user.Badges,
		IsOnline:       user.IsOnline,
		LastSeen:       user.LastSeen,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

func toModel(entry *userCacheEntry) *models.User {
	return &models.User{
		ID:             entry.ID,
		Username:       entry.Username,
		Email:          entry.Email,
		PasswordHash:   entry.PasswordHash,
		FirstName:      entry.FirstName,
		LastName:       entry.LastName,
		IsActive:       entry.IsActive,
		Role:           entry.Role,
		AvatarURL:      entry.AvatarURL,
		Location:       entry.Location,
		FavoriteGenres: entry.FavoriteGenres,
		Bio:            entry.Bio,
		ReadingGoal:    func() int { if entry.ReadingGoal != nil { return *entry.ReadingGoal }; return 0 }(),
		BooksRead:      entry.BooksRead,
		Badges:         entry.Badges,
		IsOnline:       entry.IsOnline,
		LastSeen:       entry.LastSeen,
		CreatedAt:      entry.CreatedAt,
		UpdatedAt:      entry.UpdatedAt,
	}
}

func (r *cachedUserRepository) get(ctx context.Context, key string) (*models.User, error) {
	bytes, err := r.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var entry userCacheEntry
	if err := json.Unmarshal(bytes, &entry); err != nil {
		return nil, err
	}

	return toModel(&entry), nil
}

func (r *cachedUserRepository) set(ctx context.Context, key string, user *models.User) {
	if user == nil {
		return
	}

	entry := toCache(user)

	if bytes, err := json.Marshal(entry); err == nil {
		_ = r.rdb.Set(ctx, key, bytes, r.ttl)
	}
}

func (r *cachedUserRepository) del(ctx context.Context, keys ...string) {
	if len(keys) > 0 {
		_ = r.rdb.Del(ctx, keys...).Err()
	}
}

func (r *cachedUserRepository) Create(user *models.User) error {
	if err := r.base.Create(user); err != nil {
		return err
	}

	ctx := context.Background()

	if saved, err := r.base.GetByID(user.ID); err == nil && saved != nil {
		user = saved
	}

	r.set(ctx, r.keyByID(user.ID), user)
	r.set(ctx, r.keyByEmail(user.Email), user)
	r.set(ctx, r.keyByUsername(user.Username), user)

	return nil
}

func (r *cachedUserRepository) GetByID(id uint) (*models.User, error) {
	ctx := context.Background()

	if user, err := r.get(ctx, r.keyByID(id)); err == nil {
		return user, nil
	}

	user, err := r.base.GetByID(id)
	if err != nil {
		return nil, err
	}

	r.set(ctx, r.keyByID(id), user)
	return user, nil
}

func (r *cachedUserRepository) GetByEmail(email string) (*models.User, error) {
	ctx := context.Background()

	if user, err := r.get(ctx, r.keyByEmail(email)); err == nil {
		return user, nil
	}

	user, err := r.base.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	r.set(ctx, r.keyByEmail(email), user)
	r.set(ctx, r.keyByID(user.ID), user)
	return user, nil
}

func (r *cachedUserRepository) GetByEmailIncludingDeleted(email string) (*models.User, error) {
	return r.base.GetByEmailIncludingDeleted(email)
}

func (r *cachedUserRepository) GetByUsernameIncludingDeleted(username string) (*models.User, error) {
	return r.base.GetByUsernameIncludingDeleted(username)
}

func (r *cachedUserRepository) GetByUsername(username string) (*models.User, error) {
	ctx := context.Background()

	if user, err := r.get(ctx, r.keyByUsername(username)); err == nil {
		return user, nil
	}

	user, err := r.base.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	r.set(ctx, r.keyByUsername(username), user)
	r.set(ctx, r.keyByID(user.ID), user)
	return user, nil
}

func (r *cachedUserRepository) Update(user *models.User) error {
	if err := r.base.Update(user); err != nil {
		return err
	}

	ctx := context.Background()

	r.set(ctx, r.keyByID(user.ID), user)
	r.set(ctx, r.keyByEmail(user.Email), user)
	r.set(ctx, r.keyByUsername(user.Username), user)

	return nil
}

func (r *cachedUserRepository) Delete(id uint) error {
	user, _ := r.base.GetByID(id)
	if err := r.base.Delete(id); err != nil {
		return err
	}

	ctx := context.Background()
	keys := []string{r.keyByID(id)}
	if user != nil {
		keys = append(keys, r.keyByEmail(user.Email), r.keyByUsername(user.Username))
	}

	r.del(ctx, keys...)
	return nil
}

func (r *cachedUserRepository) List(limit, offset int) ([]*models.User, error) {
	return r.base.List(limit, offset)
}

package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nevzattalhaozcan/forgotten/internal/clients"
	"github.com/nevzattalhaozcan/forgotten/internal/config"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/internal/repository"
	"gorm.io/gorm"
)

type BookService struct {
	bookRepo   repository.BookRepository
	bookClient clients.BookAPIClient
	config     *config.Config
}

func NewBookService(bookRepo repository.BookRepository, bookClient clients.BookAPIClient, config *config.Config) *BookService {
	return &BookService{
		bookRepo:   bookRepo,
		bookClient: bookClient,
		config:     config,
	}
}

func (s *BookService) CreateBook(req *models.CreateBookRequest) (*models.BookResponse, error) {
	book := &models.Book{
		Title:         req.Title,
		Author:        req.Author,
		CoverURL:      req.CoverURL,
		Genre:         req.Genre,
		Pages:         req.Pages,
		PublishedYear: req.PublishedYear,
		ISBN:          req.ISBN,
		Description:   req.Description,
	}

	if err := s.bookRepo.Create(book); err != nil {
		return nil, err
	}

	response := book.ToResponse()
	return &response, nil
}

func (s *BookService) GetBookByID(id uint) (*models.BookResponse, error) {
	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	response := book.ToResponse()
	return &response, nil
}

func (s *BookService) UpdateBook(id uint, req *models.UpdateBookRequest) (*models.BookResponse, error) {
	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	if req.Title != nil {
		book.Title = *req.Title
	}
	if req.Author != nil {
		book.Author = req.Author
	}
	if req.CoverURL != nil {
		book.CoverURL = req.CoverURL
	}
	if req.Genre != nil {
		book.Genre = req.Genre
	}
	if req.Pages != nil {
		book.Pages = req.Pages
	}
	if req.PublishedYear != nil {
		book.PublishedYear = req.PublishedYear
	}
	if req.ISBN != nil {
		book.ISBN = req.ISBN
	}
	if req.Description != nil {
		book.Description = req.Description
	}
	if req.Rating != nil {
		book.Rating = req.Rating
	}

	if err := s.bookRepo.Update(book); err != nil {
		return nil, err
	}

	response := book.ToResponse()
	return &response, nil
}

func (s *BookService) DeleteBook(id uint) error {
	_, err := s.bookRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("book not found")
		}
		return err
	}

	return s.bookRepo.Delete(id)
}

func (s *BookService) ListBooks() ([]*models.BookResponse, error) {
	books, err := s.bookRepo.List(50, 0)
	if err != nil {
		return nil, err
	}

	var responses []*models.BookResponse
	for _, book := range books {
		response := book.ToResponse()
		responses = append(responses, &response)
	}

	return responses, nil
}

// multi-source search with caching
func (s *BookService) SearchBooks(query string, limit int, source string) ([]models.BookResponse, error) {
	limitArg := limit
	if limitArg <= 0 {
		limitArg = 20
	}

	var responses []models.BookResponse
	seen := make(map[string]bool)

	if source == "local" || source == "all" {
		local, err := s.bookRepo.SearchLocal(query, limitArg)
		if err == nil {
			for _, b := range local {
				key := getBookKey(b)
				seen[key] = true
				responses = append(responses, b.ToResponse())
			}
			if source == "local" {
				return responses, nil
			}
		}
	}

	if source == "external" || source == "all" {
		exts, err := s.bookClient.SearchMerged(query, limitArg)
		if err != nil {
			exts, err = s.bookClient.Search(query, limitArg)
		}

		if err != nil {
			if source == "external" {
				return nil, fmt.Errorf("external book search error: %w", err)
			}
			log.Printf("external book search error: %v", err)
		} else {
			for _, eb := range exts {
				if len(responses) >= limitArg {
					break
				}

				key := getExternalBookKey(eb)
				if seen[key] {
					continue
				}

				var local *models.Book
				if eb.ISBN != nil {
					if b, err := s.bookRepo.GetByISBN(*eb.ISBN); err == nil {
						local = b
					}
				}
				if local == nil {
					if b, err := s.bookRepo.GetByExternalID(eb.Source, eb.ExternalID); err == nil {
						local = b
					}
				}

				if local != nil {
					if !seen[getBookKey(local)] {
						now := time.Now()
						local.LastAccessed = &now
						if err := s.bookRepo.Update(local); err != nil {
							log.Printf("warning: failed to update last_accessed for book %d: %v", local.ID, err)
						}
						responses = append(responses, local.ToResponse())
						seen[getBookKey(local)] = true
					}
					continue
				}

				book := eb.ToBook()
				if err := s.bookRepo.UpsertByExternalID(book); err != nil {
					log.Printf("failed to cache external book %s as id=%d: %v", book.Title, book.ID, err)
				}
				responses = append(responses, book.ToResponse())
				seen[key] = true
			}
		}
	}

	return responses, nil
}

func getBookKey(b *models.Book) string {
	if b.ISBN != nil && *b.ISBN != "" {
		return "isbn:" + *b.ISBN
	}

	if b.ExternalID != nil && *b.ExternalID != "" {
		return "ext:" + *b.ExternalID
	}
	return fmt.Sprintf("id:%d", b.ID)
}

func getExternalBookKey(eb *models.ExternalBook) string {
	if eb.ISBN != nil && *eb.ISBN != "" {
		return "isbn:" + *eb.ISBN
	}
	return "ext:" + eb.ExternalID
}

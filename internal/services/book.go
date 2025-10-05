package services

import (
	"errors"
	"time"

	"github.com/nevzattalhaozcan/forgotten/internal/clients"
	"github.com/nevzattalhaozcan/forgotten/internal/config"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/internal/repository"
	"gorm.io/gorm"
)

type BookService struct {
	bookRepo repository.BookRepository
	olClient *clients.OpenLibraryClient
	config   *config.Config
}

func NewBookService(bookRepo repository.BookRepository, olClient *clients.OpenLibraryClient, config *config.Config) *BookService {
	return &BookService{
		bookRepo: bookRepo,
		olClient: olClient,
		config:   config,
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

// source: "local", "external", "all"
func (s *BookService) SearchBooks(query string, limit int, source string) ([]models.BookResponse, error) {
	limitArg := limit
	if limitArg <= 0 {
		limitArg = 20
	}

	var responses []models.BookResponse
	if source == "local" || source == "all" {
		local, err := s.bookRepo.SearchLocal(query, limitArg)
		if err == nil {
			for _, b := range local {
				responses = append(responses, b.ToResponse())
			}
			if source == "local" {
				return responses, nil
			}
		}
	}

	if source == "external" || source == "all" {
		exts, err := s.olClient.Search(query, limitArg)
		if err == nil {
			for _, eb := range exts {
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
					now := time.Now()
					local.LastAccessed = &now
					_ = s.bookRepo.Update(local)
					responses = append(responses, local.ToResponse())
					continue
				}

				book := eb.ToBook()
				_ = s.bookRepo.UpsertByExternalID(book)
				responses = append(responses, book.ToResponse())
			}
		}
	}

	return responses, nil
}
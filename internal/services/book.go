package services

import (
	"errors"

	"github.com/nevzattalhaozcan/forgotten/internal/config"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/internal/repository"
	"gorm.io/gorm"
)

type BookService struct {
	bookRepo repository.BookRepository
	config   *config.Config
}

func NewBookService(bookRepo repository.BookRepository, config *config.Config) *BookService {
	return &BookService{
		bookRepo: bookRepo,
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
package repository

import (
	"strings"

	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"gorm.io/gorm"
)

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *bookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(book *models.Book) error {
	return r.db.Create(book).Error
}

func (r *bookRepository) GetByID(id uint) (*models.Book, error) {
	var book models.Book
	err := r.db.First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) Update(book *models.Book) error {
	return r.db.Save(book).Error
}

func (r *bookRepository) Delete(id uint) error {
	return r.db.Delete(&models.Book{}, id).Error
}

func (r *bookRepository) List(limit, offset int) ([]*models.Book, error) {
	var books []*models.Book
	err := r.db.Limit(limit).Offset(offset).Find(&books).Error
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (r *bookRepository) GetByExternalID(source, externalID string) (*models.Book, error) {
	var book models.Book
	if err := r.db.Where("source = ? AND external_id = ?", source, externalID).First(&book).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) GetByISBN(isbn string) (*models.Book, error) {
	var book models.Book
	if err := r.db.Where("isbn = ?", isbn).First(&book).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) UpsertByExternalID(book *models.Book) error {
	if book.ExternalID == nil || book.Source == nil {
		return r.db.Save(book).Error
	}

	var existing models.Book
	err := r.db.Where("source = ? AND external_id = ?", *book.Source, *book.ExternalID).First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return r.db.Create(book).Error
		}
		return err
	}

	existing.Title = book.Title
	existing.Author = book.Author
	existing.CoverURL = book.CoverURL
	existing.Genre = book.Genre
	existing.Pages = book.Pages
	existing.PublishedYear = book.PublishedYear
	existing.ISBN = book.ISBN
	existing.Description = book.Description
	existing.Rating = book.Rating
	now := book.LastAccessed
	if now != nil {
		existing.LastAccessed = book.LastAccessed
	}
	return r.db.Save(&existing).Error
}

func (r *bookRepository) SearchLocal(query string, limit int) ([]*models.Book, error) {
	q := "%" + strings.ToLower(query) + "%"
	var books []*models.Book
	if err := r.db.Where("LOWER(title) LIKE ? OR LOWER(author) LIKE ?", q, q).Limit(limit).Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}
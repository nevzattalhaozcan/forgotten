package repository

import (
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) Create(book *models.Book) error {
	return r.db.Create(book).Error
}

func (r *BookRepository) GetByID(id uint) (*models.Book, error) {
	var book models.Book
	err := r.db.First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) Update(book *models.Book) error {
	return r.db.Save(book).Error
}

func (r *BookRepository) Delete(id uint) error {
	return r.db.Delete(&models.Book{}, id).Error
}

func (r *BookRepository) List(limit, offset int) ([]*models.Book, error) {
	var books []*models.Book
	err := r.db.Limit(limit).Offset(offset).Find(&books).Error
	if err != nil {
		return nil, err
	}
	return books, nil
}
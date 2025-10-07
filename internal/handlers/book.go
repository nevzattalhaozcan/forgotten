package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nevzattalhaozcan/forgotten/internal/models"
	"github.com/nevzattalhaozcan/forgotten/internal/services"
)

type BookHandler struct {
	bookService *services.BookService
	validator   *validator.Validate
}

func NewBookHandler(bookService *services.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
		validator:   validator.New(),
	}
}

// @Summary Create a new book
// @Description Add a new book to the collection
// @Tags Books
// @Accept json
// @Produce json
// @Param request body models.CreateBookRequest true "Book data"
// @Success 201 {object} map[string]interface{} "Book created successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Router /api/v1/books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	var req models.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := h.bookService.CreateBook(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "book created successfully",
		"book":    book,
	})
}

// @Summary Get a book by ID
// @Description Retrieve a book's details by its ID
// @Tags Books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} map[string]interface{} "Book details"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Book not found"
// @Router /api/v1/books/{id} [get]
func (h *BookHandler) GetBookByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	book, err := h.bookService.GetBookByID(uint(id))
	if err != nil {
		if err.Error() == "book not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{ "book": book })
}

// @Summary Update a book
// @Description Update book information
// @Tags Books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param request body models.UpdateBookRequest true "Update book data"
// @Success 200 {object} map[string]interface{} "Book updated successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Book not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/books/{id} [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	var req models.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := h.bookService.UpdateBook(uint(id), &req)
	if err != nil {
		if err.Error() == "book not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "book updated successfully",
		"book":    book,
	})
}

// @Summary Delete a book
// @Description Delete a book by ID
// @Tags Books
// @Produce json
// @Param id path int true "Book ID"
// @Success 204 {object} nil "Book deleted successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Book not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/books/{id} [delete]
func (h *BookHandler) DeleteBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book ID"})
		return
	}

	if err := h.bookService.DeleteBook(uint(id)); err != nil {
		if err.Error() == "book not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary List all books
// @Description Retrieve a list of all books
// @Tags Books
// @Produce json
// @Success 200 {object} map[string]interface{} "List of books"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/books [get]
func (h *BookHandler) ListBooks(c *gin.Context) {
	books, err := h.bookService.ListBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"books": books})
}

// @Summary Search for books
// @Description Search for books using an external API
// @Tags Books
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Number of results to return" default(20)
// @Param source query string false "Source to search (google, isbndb, openlibrary)" default(all)
// @Success 200 {object} map[string]interface{} "Search results"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/books/search [get]
func (h *BookHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	source := c.DefaultQuery("source", "all")

	results, err := h.bookService.SearchBooks(query, limit, source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"books": results})
}
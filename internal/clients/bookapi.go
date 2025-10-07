package clients

import (
	"fmt"
	"log"

	"github.com/nevzattalhaozcan/forgotten/internal/models"
)

type BookAPIClient interface {
	Search(query string, limit int) ([]*models.ExternalBook, error)
	SearchMerged(query string, limit int) ([]*models.ExternalBook, error)
}

type MultiSourceClient struct {
	googleBooks     *GoogleBooksClient
	openLibrary     *OpenLibraryClient
	preferredSource string
}

func NewMultiSourceClient(googleAPIKey string, preferredSource string) *MultiSourceClient {
	var googleClient *GoogleBooksClient
	if googleAPIKey != "" {
		googleClient = NewGoogleBooksClient(googleAPIKey)
	}
	return &MultiSourceClient{
		googleBooks:     googleClient,
		openLibrary:     NewOpenLibraryClient(),
		preferredSource: preferredSource,
	}
}

func (c *MultiSourceClient) Search(query string, limit int) ([]*models.ExternalBook, error) {
	var books []*models.ExternalBook
	var err error

	if c.preferredSource == "google" || c.preferredSource == "all" {
		if c.googleBooks != nil {
			books, err = c.googleBooks.Search(query, limit, true)
			if err == nil && len(books) > 0 {
				return books, nil
			}
			if err != nil {
				fmt.Printf("Google Books search error: %v, falling back to Open Library\n", err)
			}
		}
	}

	books, err = c.openLibrary.Search(query, limit)
	if err != nil {
		return nil, fmt.Errorf("all sources failed: %w", err)
	}
	log.Printf("Found %d books from Open Library\n", len(books))
	return books, nil
}

func (c *MultiSourceClient) SearchMerged(query string, limit int) ([]*models.ExternalBook, error) {
	var allBooks []*models.ExternalBook
	seen := make(map[string]bool)

	if c.googleBooks != nil {
		gbBooks, err := c.googleBooks.SearchMultiLang(query, limit)
		if err == nil {
			for _, b := range gbBooks {
				key := getExternalBookKey(b)
				if !seen[key] {
					seen[key] = true
					allBooks = append(allBooks, b)
				}
			}
			log.Printf("Found %d books from Google Books\n", len(gbBooks))
		} else {
			fmt.Printf("Google Books search error: %v\n", err)
		}
	}

	if len(allBooks) < limit {
		remaining := limit - len(allBooks)
		olBooks, err := c.openLibrary.Search(query, remaining*2)
		if err == nil {
			for _, b := range olBooks {
				if len(allBooks) >= limit {
					break
				}
				key := getExternalBookKey(b)
				if !seen[key] {
					seen[key] = true
					allBooks = append(allBooks, b)
				}
			}
			log.Printf("Found %d books from Open Library\n", len(olBooks))
		}
	}

	if len(allBooks) == 0 {
		return nil, fmt.Errorf("no books found from any source")
	}

	return allBooks, nil
}

func getExternalBookKey(eb *models.ExternalBook) string {
	if eb.ISBN != nil && *eb.ISBN != "" {
		return "isbn:" + *eb.ISBN
	}
	return "ext:" + eb.ExternalID
}
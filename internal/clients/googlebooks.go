package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/nevzattalhaozcan/forgotten/internal/models"
)

type GoogleBooksClient struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

func NewGoogleBooksClient(apiKey string) *GoogleBooksClient {
	return &GoogleBooksClient{
		APIKey:  apiKey,
		BaseURL: "https://www.googleapis.com/books/v1",
		Client:  &http.Client{},
	}
}

type gbSearchResponse struct {
	Items []struct {
		ID         string
		VolumeInfo struct {
			Title               string   `json:"title"`
			Authors             []string `json:"authors"`
			Publisher           string   `json:"publisher"`
			PublishedDate       string   `json:"publishedDate"`
			Description         string   `json:"description"`
			PageCount           int      `json:"pageCount"`
			Categories          []string `json:"categories"`
			AverageRating       float32  `json:"averageRating"`
			RatingsCount        int      `json:"ratingsCount"`
			Language            string   `json:"language"`
			IndustryIdentifiers []struct {
				Type       string `json:"type"`
				Identifier string `json:"identifier"`
			} `json:"industryIdentifiers"`
			ImageLinks struct {
				Thumbnail      string `json:"thumbnail"`
				SmallThumbnail string `json:"smallThumbnail"`
			} `json:"imageLinks"`
		}
	}
}

func (c *GoogleBooksClient) Search(query string, limit int, tr bool) ([]*models.ExternalBook, error) {
	params := url.Values{}
	params.Set("q", query)
	if tr {
		params.Set("langRestrict", "tr")
	}
	params.Set("maxResults", fmt.Sprintf("%d", limit))
	params.Set("key", c.APIKey)

	u := fmt.Sprintf("%s/volumes?%s", c.BaseURL, params.Encode())
	resp, err := c.Client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data gbSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var out []*models.ExternalBook
	for _, item := range data.Items {
		vi := item.VolumeInfo

		var author *string
		if len(vi.Authors) > 0 {
			a := vi.Authors[0]
			author = &a
		}

		var isbn *string
		for _, id := range vi.IndustryIdentifiers {
			if id.Type == "ISBN_13" {
				isbn = &id.Identifier
				break
			}
		}

		var genre *string
		if len(vi.Categories) > 0 {
			g := vi.Categories[0]
			genre = &g
		}

		var coverURL *string
		if vi.ImageLinks.Thumbnail != "" {
			cover := vi.ImageLinks.Thumbnail
			coverURL = &cover
		}

		var year *int
		if len(vi.PublishedDate) >= 4 {
			y := 0
			_, err := fmt.Sscanf(vi.PublishedDate[:4], "%d", &y)
			if err == nil && y > 0 {
				year = &y
			}
		}

		var desc *string
		if vi.Description != "" {
			desc = &vi.Description
		}

		var rating *float32
		if vi.AverageRating > 0 {
			rating = &vi.AverageRating
		}

		eb := &models.ExternalBook{
			ExternalID:    fmt.Sprintf("GB_%s", item.ID),
			Source:        "googlebooks",
			Title:         vi.Title,
			Author:        author,
			CoverURL:      coverURL,
			Genre:         genre,
			Pages:         &vi.PageCount,
			PublishedYear: year,
			ISBN:          isbn,
			Description:   desc,
			Rating:        rating,
		}
		out = append(out, eb)
	}
	return out, nil
}

func (c *GoogleBooksClient) SearchMultiLang(query string, limit int) ([]*models.ExternalBook, error) {
	books, err := c.Search(query, limit, true)
	if err == nil && len(books) > limit/2 {
		return books, nil
	}

	return c.Search(query, limit, false)
}

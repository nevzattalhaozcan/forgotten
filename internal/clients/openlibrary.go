package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/nevzattalhaozcan/forgotten/internal/models"
)

type OpenLibraryClient struct {
	BaseURL string
	Client  *http.Client
}

func NewOpenLibraryClient() *OpenLibraryClient {
	return &OpenLibraryClient{
		BaseURL: "https://openlibrary.org",
		Client:  &http.Client{Timeout: 10 * time.Second},
	}
}

type olSearchResponse struct {
	Docs []struct {
		Key              string   `json:"key"`
		Title            string   `json:"title"`
		AuthorName       []string `json:"author_name"`
		FirstPublishYear *int     `json:"first_publish_year"`
		ISBN             []string `json:"isbn"`
		CoverI           *int     `json:"cover_i"`
	} `json:"docs"`
}

func (c *OpenLibraryClient) Search(query string, limit int) ([]*models.ExternalBook, error) {
	u := fmt.Sprintf("%s/search.json?q=%s&language=tur&limit=%d",
		c.BaseURL, url.QueryEscape(query),
		limit)
	resp, err := c.Client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data olSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var out []*models.ExternalBook
	for _, d := range data.Docs {
		var author *string
		if len(d.AuthorName) > 0 {
			a := d.AuthorName[0]
			author = &a
		}
		var isbn *string
		if len(d.ISBN) > 0 {
			i := d.ISBN[0]
			isbn = &i
		}
		var pubYear *int
		if d.FirstPublishYear != nil {
			pubYear = d.FirstPublishYear
		}
		var coverURL *string
		if d.CoverI != nil {
			u := fmt.Sprintf("https://covers.openlibrary.org/b/id/%d-L.jpg", *d.CoverI)
			coverURL = &u
		}

		eb := &models.ExternalBook{
			ExternalID: fmt.Sprintf("OL%s", d.Key),
			Source:     "openlibrary",
			Title:      d.Title,
			Author:    author,
			CoverURL:  coverURL,
			PublishedYear: pubYear,
			ISBN:      isbn,
		}
		out = append(out, eb)
	}
	return out, nil
}
package search

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/agnivade/levenshtein"
)

type BlogAuthorPhoto struct {
	URL string `json:"url"`
}

type BlogAuthor struct {
	Username string          `json:"username"`
	Photo    BlogAuthorPhoto `json:"photo"`
}

type BlogPhoto struct {
	URL string `json:"url"`
}

type BlogContentBlock struct {
	Type     string          `json:"type"`
	Level    int             `json:"level,omitempty"`
	Format   string          `json:"format,omitempty"`
	Children []BlogTextBlock `json:"children"`
}

type BlogTextBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
	Bold bool   `json:"bold,omitempty"`
}

type VisaStatus string

const (
	VisaFree      VisaStatus = "visa-free"
	VisaRequired  VisaStatus = "visa"
	VisaOnArrival VisaStatus = "visa-on-arrival"
	ETA           VisaStatus = "eta"
)

type BlogData struct {
	ID          int                `json:"id"`
	Title       string             `json:"title"`
	Slug        string             `json:"slug"`
	Content     []BlogContentBlock `json:"content"`
	DocumentId  string             `json:"documentId"`
	VisaStatus  VisaStatus         `json:"visaStatus"`
	MainPhoto   BlogPhoto          `json:"mainPhoto"`
	Author      BlogAuthor         `json:"author"`
	Description string             `json:"description"`
}

type BlogResponse struct {
	Data []BlogData `json:"data"`
}

type BlogSearchResult struct {
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	VisaStatus  VisaStatus `json:"visaStatus"`
	MainPhoto   BlogPhoto  `json:"mainPhoto"`
	Description string     `json:"description"`
}

func betterMatch(query, target string) bool {
	query = strings.ToLower(query)
	target = strings.ToLower(target)

	if target == query || strings.HasPrefix(target, query) {
		return true
	}

	distance := levenshtein.ComputeDistance(query, target)
	threshold := len(target) / 4

	return distance > 0 && distance <= threshold
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	limitParam := r.URL.Query().Get("limit")
	var limit int
	fmt.Sscanf(limitParam, "%d", &limit)

	if query == "" {
		http.Error(w, "Query param is required", http.StatusBadRequest)
		return
	}

	resp, err := http.Get("https://dashboard.sorgulat.com/api/passport-blogs?populate[author][populate]=photo&populate=mainPhoto")
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "Failed to fetch blog data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var blogResp BlogResponse
	if err := json.Unmarshal(body, &blogResp); err != nil {
		http.Error(w, "Failed to parse blog data", http.StatusInternalServerError)
		return
	}

	var results []BlogSearchResult

	for _, blog := range blogResp.Data {
		matched := false

		if betterMatch(query, blog.Title) ||
			betterMatch(query, blog.Description) ||
			betterMatch(query, blog.Slug) {
			matched = true
		}

		if !matched {
			fieldsToCheck := []string{blog.Title, blog.Description, blog.Slug}
			for _, field := range fieldsToCheck {
				words := strings.Fields(field)
				for _, word := range words {
					if betterMatch(query, word) {
						matched = true
						break
					}
				}
				if matched {
					break
				}
			}
		}

		if matched {
			results = append(results, BlogSearchResult{
				Title:       blog.Title,
				Slug:        blog.Slug,
				VisaStatus:  blog.VisaStatus,
				MainPhoto:   blog.MainPhoto,
				Description: blog.Description,
			})
		}

		if limit > 0 && len(results) >= limit {
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if len(results) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Sonuç bulunamadı"})
		return
	}

	json.NewEncoder(w).Encode(results)
}

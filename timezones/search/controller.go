package search

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sorgulat-api/timezones/models"
	"sorgulat-api/timezones/utils"
	"strings"
	"time"

	"github.com/agnivade/levenshtein"
)

var (
	cities    = utils.LoadData[models.City]("cities")
	countries = utils.LoadData[models.Country]("countries")
)

type ResponseItem struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Timezone string `json:"timezone"`
	Time     string `json:"time"`
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

func getCurrentTime(timezone string) string {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Now().Format(time.RFC3339)
	}
	return time.Now().In(loc).Format(time.RFC3339)
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

	var results []ResponseItem

	for _, city := range cities {
		if betterMatch(query, city.Name) {
			results = append(results, ResponseItem{
				Name:     fmt.Sprintf("%s, %s", city.Name, city.Country),
				Slug:     city.Slug,
				Timezone: city.Timezone,
				Time:     getCurrentTime(city.Timezone),
			})
			if limit > 0 && len(results) >= limit {
				break
			}
		}
	}

	for _, country := range countries {
		if betterMatch(query, country.Name) {
			results = append(results, ResponseItem{
				Name:     country.Name,
				Slug:     country.Slug,
				Timezone: country.Timezone,
				Time:     getCurrentTime(country.Timezone),
			})
			if limit > 0 && len(results) >= limit {
				break
			}
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

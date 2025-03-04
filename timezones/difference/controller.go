package difference

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sorgulat-api/timezones/models"
	"sorgulat-api/timezones/utils"
	"sort"
	"strings"
	"time"
)

var (
	cities    = utils.LoadData[models.City]("cities")
	countries = utils.LoadData[models.Country]("countries")
)

func getUTCOffset(timezone string) float64 {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		log.Printf("Zaman dilimi yüklenemedi: %s, hata: %v", timezone, err)
		return 0
	}
	utcNow := time.Now().UTC()
	localTime := utcNow.In(location)
	_, offsetSeconds := localTime.Zone()
	return float64(offsetSeconds) / 3600.0
}

func formatCityName(city string) string {
	words := strings.Split(city, "-")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

func findReferenceTimezone(name string) string {
	for _, city := range cities {
		if city.Slug == name {
			return city.Timezone
		}
	}
	for _, country := range countries {
		if country.Slug == name {
			return country.Timezone
		}
	}
	return ""
}

func getRandomCities(exclude string) []models.City {
	var filteredCities []models.City
	for _, city := range cities {
		if city.Slug != exclude {
			filteredCities = append(filteredCities, city)
		}
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(filteredCities), func(i, j int) {
		filteredCities[i], filteredCities[j] = filteredCities[j], filteredCities[i]
	})

	return filteredCities[:20]
}

func GetDifferenceBySlug(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	slug := parts[len(parts)-1]
	name := formatCityName(slug)
	refTimezone := findReferenceTimezone(slug)
	if refTimezone == "" {
		http.Error(w, "Şehir veya ülke bulunamadı", http.StatusNotFound)
		return
	}

	refOffset := getUTCOffset(refTimezone)

	randomCities := getRandomCities(slug)

	response := struct {
		From         string        `json:"from"`
		LocationText string        `json:"locationText"`
		Destinations []models.City `json:"destinations"`
	}{From: name, LocationText: utils.GetLocationSuffix(name), Destinations: []models.City{}}

	for _, city := range randomCities {
		cityOffset := getUTCOffset(city.Timezone)
		offsetDiff := cityOffset - refOffset

		response.Destinations = append(response.Destinations, models.City{
			Name:   city.Name,
			Slug:   city.Slug,
			Offset: offsetDiff,
		})
	}
	sort.Slice(response.Destinations, func(i, j int) bool {
		return response.Destinations[i].Offset < response.Destinations[j].Offset
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

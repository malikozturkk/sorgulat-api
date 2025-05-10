package compare

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sorgulat-api/timezones/models"
	"sorgulat-api/timezones/utils"
)

var (
	allCities    = utils.LoadData[models.City]("cities")
	allCountries = utils.LoadData[models.Country]("countries")
)

func SitemapHandler(w http.ResponseWriter, r *http.Request) {
	var urls []string
	var allItems []models.City
	allItems = append(allItems, allCities...)
	for _, country := range allCountries {
		allItems = append(allItems, models.City{
			Name:     country.Name,
			Slug:     country.Slug,
			Timezone: country.Timezone,
			Country:  country.Name,
			Latitude: country.Latitude,
			Longitude: country.Longitude,
		})
	}

	for _, from := range allItems {
		for _, to := range allItems {
			if from.Slug != to.Slug {
				url := fmt.Sprintf("from-%s-to-%s", from.Slug, to.Slug)
				urls = append(urls, url)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urls)
}

package city

import (
	"encoding/json"
	"net/http"
	"sorgulat-api/timezones/models"
	"sorgulat-api/timezones/utils"
	"strconv"
)

var cities = utils.LoadData[models.City]("cities")

func GetCityTimeZone(w http.ResponseWriter, r *http.Request) {
	limitParam := r.URL.Query().Get("limit")
	limit := len(cities)

	if limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			if parsedLimit < limit {
				limit = parsedLimit
			}
		}
	}

	selectedIndexes := utils.GetRandomIndexes(len(cities))

	updatedCities := make([]models.City, len(cities))
	copy(updatedCities, cities)

	for i := range updatedCities {
		if _, exists := selectedIndexes[i]; exists {
			randomType := utils.GetRandomType()
			updatedCities[i].Type = &randomType
		}
	}

	updatedCities = updatedCities[:limit]
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCities)
}

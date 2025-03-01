package city

import (
	"encoding/json"
	"net/http"
	"sorgulat-api/timezones/models"
	"sorgulat-api/timezones/utils"
)

var cities = utils.LoadData[models.City]("cities")

func GetCityTimeZone(w http.ResponseWriter, r *http.Request) {
	selectedIndexes := utils.GetRandomIndexes(len(cities))

	updatedCities := make([]models.City, len(cities))
	copy(updatedCities, cities)

	for i := range updatedCities {
		if _, exists := selectedIndexes[i]; exists {
			randomType := utils.GetRandomType()
			updatedCities[i].Type = &randomType
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCities)
}

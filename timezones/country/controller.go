package country

import (
	"encoding/json"
	"net/http"
	"sorgulat-api/timezones/models"
	"sorgulat-api/timezones/utils"
	"strconv"
)

var countries = utils.LoadData[models.Country]("countries")

func GetCountryTimeZone(w http.ResponseWriter, r *http.Request) {
	limitParam := r.URL.Query().Get("limit")
	limit := len(countries)

	if limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			if parsedLimit < limit {
				limit = parsedLimit
			}
		}
	}

	selectedIndexes := utils.GetRandomIndexes(len(countries))

	updatedCountries := make([]models.Country, len(countries))
	copy(updatedCountries, countries)

	for i := range updatedCountries {
		if _, exists := selectedIndexes[i]; exists {
			randomType := utils.GetRandomType()
			updatedCountries[i].Type = &randomType
		}
	}

	updatedCountries = updatedCountries[:limit]
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCountries)
}

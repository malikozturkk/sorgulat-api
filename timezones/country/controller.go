package country

import (
	"encoding/json"
	"net/http"
	"sorgulat-api/timezones/models"
	"sorgulat-api/timezones/utils"
)

var countries = utils.LoadData[models.Country]("countries")

func GetCountryTimeZone(w http.ResponseWriter, r *http.Request) {
	selectedIndexes := utils.GetRandomIndexes(len(countries))

	updatedCountries := make([]models.Country, len(countries))
	copy(updatedCountries, countries)

	for i := range updatedCountries {
		if _, exists := selectedIndexes[i]; exists {
			randomType := utils.GetRandomType()
			updatedCountries[i].Type = &randomType
		}
	}

	updatedCountries = updatedCountries[:45]
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCountries)
}

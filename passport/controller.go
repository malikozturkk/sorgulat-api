package passport

import (
	"encoding/json"
	"net/http"
	"sorgulat-api/passport/models"
	"sorgulat-api/timezones/utils"
)

var (
	countries = utils.LoadData[models.PassportCountries]("countries", "passport/")
)

func GetCountriesPassport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	counts := make(map[string]int)
	for _, item := range countries {
		counts[item.Value]++
	}

	response := struct {
		Countries []models.PassportCountries `json:"countries"`
		Counts    map[string]int             `json:"counts"`
	}{
		Countries: countries,
		Counts:    counts,
	}

	json.NewEncoder(w).Encode(response)
}

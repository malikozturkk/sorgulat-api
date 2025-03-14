package passport

import (
	"encoding/json"
	"net/http"
	"sorgulat-api/passport"
	"sorgulat-api/passport/models"
)

func GetCountriesPassport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	counts := make(map[string]int)
	for _, item := range passport.Countries {
		counts[item.Value]++
	}

	response := struct {
		Countries []models.PassportCountries `json:"countries"`
		Counts    map[string]int             `json:"counts"`
	}{
		Countries: passport.Countries,
		Counts:    counts,
	}

	json.NewEncoder(w).Encode(response)
}

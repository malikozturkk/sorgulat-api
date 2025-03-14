package passport

import (
	"encoding/json"
	"net/http"
	"sorgulat-api/passport"
	"sorgulat-api/passport/models"
	"strings"
)

func GetFilteredCountriesPassport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slug := strings.TrimPrefix(r.URL.Path, "/passport/")

	filterMap := map[string]string{
		"vizesiz-seyahat":     "Vizesiz",
		"vizeli-seyahat":      "Vizeli",
		"kapida-vize-seyahat": "Kapıda Vize",
		"eta-seyahat":         "eTA",
	}

	filterTypeMap := map[string]string{
		"Vizesiz":     "visa-free",
		"Vizeli":      "visa",
		"Kapıda Vize": "visa-on-arrival",
		"eTA":         "eta",
	}

	filter, exists := filterMap[slug]
	if !exists {
		http.Error(w, `{"error": "Geçersiz slug"}`, http.StatusBadRequest)
		return
	}

	filterType, typeExists := filterTypeMap[filter]
	if !typeExists {
		http.Error(w, `{"error": "Bilinmeyen filtre tipi"}`, http.StatusInternalServerError)
		return
	}

	var filteredCountries []models.PassportCountries
	for _, country := range passport.Countries {
		if country.Value == filter {
			filteredCountries = append(filteredCountries, country)
		}
	}
	filteredCountries = append(filteredCountries, models.PassportCountries{
		Country: "TR",
		Value:   "main",
	})

	count := len(filteredCountries)

	response := struct {
		Countries  []models.PassportCountries `json:"countries"`
		Count      int                        `json:"count"`
		Filter     string                     `json:"filter"`
		FilterType string                     `json:"filter_type"`
		Slug       string                     `json:"slug"`
	}{
		Countries:  filteredCountries,
		Count:      count,
		Filter:     filter,
		FilterType: filterType,
		Slug:       slug,
	}

	json.NewEncoder(w).Encode(response)
}

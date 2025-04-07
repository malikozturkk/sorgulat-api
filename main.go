package main

import (
	"log"
	"net/http"
	passport "sorgulat-api/passport/controllers"
	"sorgulat-api/timezones"
	"sorgulat-api/timezones/city"
	"sorgulat-api/timezones/compare"
	"sorgulat-api/timezones/country"
	"sorgulat-api/timezones/difference"
	"sorgulat-api/timezones/search"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/timezones/city", city.GetCityTimeZone)
	mux.HandleFunc("/timezones/country", country.GetCountryTimeZone)
	mux.HandleFunc("/timezones/", timezones.GetTimezoneBySlug)
	mux.HandleFunc("/timezones/difference/", difference.GetDifferenceBySlug)
	mux.HandleFunc("/timezones/search", search.SearchHandler)
	mux.HandleFunc("/compare", compare.CompareTimezones)
	mux.HandleFunc("/passport", passport.GetCountriesPassport)
	mux.HandleFunc("/passport/", passport.GetFilteredCountriesPassport)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(mux)))
}

package main

import (
	"log"
	"net/http"
	"sorgulat-api/timezones"
	"sorgulat-api/timezones/city"
	"sorgulat-api/timezones/country"
	"sorgulat-api/timezones/difference"
	"sorgulat-api/timezones/search"
)

func main() {
	http.HandleFunc("/timezones/city", city.GetCityTimeZone)
	http.HandleFunc("/timezones/country", country.GetCountryTimeZone)
	http.HandleFunc("/timezones/", timezones.GetTimezoneBySlug)
	http.HandleFunc("/timezones/difference/", difference.GetDifferenceBySlug)
	http.HandleFunc("/timezones/search", search.SearchHandler)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

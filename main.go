package main

import (
	"log"
	"net/http"
	"sorgulat-api/timezones"
	"sorgulat-api/timezones/city"
	"sorgulat-api/timezones/country"
)

func main() {
	http.HandleFunc("/timezones/city", city.GetCityTimeZone)
	http.HandleFunc("/timezones/country", country.GetCountryTimeZone)
	http.HandleFunc("/timezones/", timezones.GetTimezoneBySlug)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

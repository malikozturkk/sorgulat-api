package compare

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sorgulat-api/timezones/models"
	"sorgulat-api/timezones/utils"
	"strings"
	"time"
)

var (
	cities    = utils.LoadData[models.City]("cities")
	countries = utils.LoadData[models.Country]("countries")
)

func findCityOrCountry(slug string) (models.City, string, error) {
	for _, city := range cities {
		if city.Slug == slug {
			return city, city.Timezone, nil
		}
	}
	for _, country := range countries {
		if country.Slug == slug {
			return models.City{
				Name:     country.Name,
				Slug:     country.Slug,
				Timezone: country.Timezone,
				Country:  country.Name,
			}, country.Timezone, nil
		}
	}
	return models.City{}, "", http.ErrMissingFile
}

func getTimeInLocation(tz string) (time.Time, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().In(loc), nil
}

func formatLocationName(place models.City) string {
	if place.Name == place.Country {
		return place.Name
	}
	return place.Name + ", " + place.Country
}

func getDiff(fromTime, toTime time.Time) (time.Duration, string) {
	_, fromOffset := fromTime.Zone()
	_, toOffset := toTime.Zone()

	diffSeconds := toOffset - fromOffset
	diff := time.Duration(diffSeconds) * time.Second

	absHours := int(math.Abs(diff.Hours()))
	absMins := int(math.Abs(diff.Minutes())) % 60

	return diff, formatDiff(absHours, absMins)
}

func getDifferenceText(fromTime, toTime time.Time, from, to models.City) string {
	fromFormatted := formatLocationName(from)
	toFormatted := formatLocationName(to)
	diff, formattedDiff := getDiff(fromTime, toTime)

	if diff == 0 {
		return fmt.Sprintf("%s ile %s arasında saat farkı yok.", fromFormatted, toFormatted)
	}

	if diff > 0 {
		return fmt.Sprintf("%s konumu %s konumundan %s ileridedir.", toFormatted, fromFormatted, formattedDiff)
	}
	return fmt.Sprintf("%s konumu %s konumundan %s geridedir.", toFormatted, fromFormatted, formattedDiff)
}

func formatDiff(hours, mins int) string {
	if hours == 0 {
		return formatUnit(mins, "dakika")
	}
	if mins == 0 {
		return formatUnit(hours, "saat")
	}
	return formatUnit(hours, "saat") + " " + formatUnit(mins, "dakika")
}

func formatUnit(val int, unit string) string {
	return fmt.Sprintf("%d %s", val, unit)
}

func buildHourTable(fromTime, toTime time.Time) []models.HourPair {
	table := make([]models.HourPair, 0, 24)
	for i := 0; i < 24; i++ {
		ft := fromTime.Add(time.Duration(i) * time.Hour).Format("15:04")
		tt := toTime.Add(time.Duration(i) * time.Hour).Format("15:04")
		table = append(table, models.HourPair{
			From: ft,
			To:   tt,
		})
	}
	return table
}

func CompareTimezones(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	fromSlug := strings.ToLower(query.Get("from"))
	toSlug := strings.ToLower(query.Get("to"))

	if fromSlug == "" || toSlug == "" {
		http.Error(w, "Eksik parametre. from ve to gereklidir.", http.StatusBadRequest)
		return
	}

	fromCity, fromTZ, err := findCityOrCountry(fromSlug)
	if err != nil {
		http.Error(w, "From şehri/ülkesi bulunamadı", http.StatusNotFound)
		return
	}
	toCity, toTZ, err := findCityOrCountry(toSlug)
	if err != nil {
		http.Error(w, "To şehri/ülkesi bulunamadı", http.StatusNotFound)
		return
	}

	fromTime, err := getTimeInLocation(fromTZ)
	if err != nil {
		http.Error(w, "From timezone yüklenemedi", http.StatusInternalServerError)
		return
	}
	toTime, err := getTimeInLocation(toTZ)
	if err != nil {
		http.Error(w, "To timezone yüklenemedi", http.StatusInternalServerError)
		return
	}

	diffText := getDifferenceText(fromTime, toTime, fromCity, toCity)
	hourTable := buildHourTable(fromTime, toTime)
	diff, formattedDiff := getDiff(fromTime, toTime)

	resp := models.CompareResponse{
		From: models.CityWithTime{
			Name:     fromCity.Name,
			Slug:     fromCity.Slug,
			Timezone: fromCity.Timezone,
			Country:  fromCity.Country,
			DateTime: fromTime.Format(time.RFC3339),
		},
		To: models.CityWithTime{
			Name:     toCity.Name,
			Slug:     toCity.Slug,
			Timezone: toCity.Timezone,
			Country:  toCity.Country,
			DateTime: toTime.Format(time.RFC3339),
		},
		DifferenceText: diffText,
		HourTable:      hourTable,
		Diff:           diff,
		FormattedDiff:  formattedDiff,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

package timezones

import (
	"encoding/json"
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

func getTimeForLocation(timezone string, locationName string, slug string, country *string, typeOpt *string) (models.Response, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return models.Response{}, err
	}
	now := time.Now().In(loc)
	suffix := utils.GetLocationSuffix(locationName)
	populerCities := make([]models.CityTime, 0, 6)
	for i, city := range cities {
		if i >= 6 {
			break
		}
		cityLoc, err := time.LoadLocation(city.Timezone)
		if err != nil {
			continue
		}
		cityNow := time.Now().In(cityLoc)
		populerCities = append(populerCities, models.CityTime{
			Slug:     city.Slug,
			Name:     city.Name,
			Hour:     cityNow.Hour(),
			Minute:   cityNow.Minute(),
			DateTime: cityNow.Format(time.RFC3339),
			Selected: city.Slug == timezone,
		})
	}

	return models.Response{
		Year:         now.Year(),
		Month:        int(now.Month()),
		Day:          now.Day(),
		Hour:         now.Hour(),
		Minute:       now.Minute(),
		Seconds:      now.Second(),
		MilliSeconds: now.Nanosecond() / 1e6,
		DateTime:     now.Format(time.RFC3339),
		Date:         now.Format("01/02/2006"),
		Time:         now.Format("15:04"),
		DayOfWeek:    now.Weekday().String(),
		DstActive:    now.IsDST(),
		LocationText: suffix,
		Timezone: models.City{
			Name:     locationName,
			Slug:     slug,
			Timezone: timezone,
			Country: func() string {
				if country != nil {
					return *country
				}
				return ""
			}(),
			Type: typeOpt,
		},
		PopulerCities: populerCities,
	}, nil
}

func GetTimezoneBySlug(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	name := parts[len(parts)-1]

	var (
		found        models.City
		foundCountry models.Country
		locationName string
		country      *string
		typeOpt      *string
	)

	for _, city := range cities {
		if city.Slug == name {
			found = city
			locationName = city.Name
			country = &city.Country
			break
		}
	}

	if found.Slug == "" {
		for _, countryData := range countries {
			if countryData.Slug == name {
				foundCountry = countryData
				locationName = countryData.Name
				country = &foundCountry.Name
				break
			}
		}
	}

	if found.Slug == "" && foundCountry.Slug == "" {
		http.Error(w, "Böyle bir şehir veya ülke bulunamadı", http.StatusNotFound)
		return
	}

	timezone := found.Timezone
	if found.Slug == "" {
		timezone = foundCountry.Timezone
	}

	response, err := getTimeForLocation(timezone, locationName, name, country, typeOpt)
	if err != nil {
		http.Error(w, "Saat hesaplanırken hata oluştu", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

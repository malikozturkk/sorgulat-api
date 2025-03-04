package timezones

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sorgulat-api/timezones/models"
	"sorgulat-api/timezones/utils"
	"sort"
	"strings"
	"time"

	"github.com/nathan-osman/go-sunrise"
)

var (
	cities    = utils.LoadData[models.City]("cities")
	countries = utils.LoadData[models.Country]("countries")
)

func getTimeForLocation(timezone string, locationName string, slug string, country *string, typeOpt *string, latitude float64, longitude float64, allCities []models.City) (models.Response, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return models.Response{}, err
	}
	now := time.Now().In(loc)

	year, month, day := now.Date()
	sunriseTime, sunsetTime := sunrise.SunriseSunset(latitude, longitude, year, month, day)

	sunriseLocal := sunriseTime.In(loc).Format("15:04")
	sunsetLocal := sunsetTime.In(loc).Format("15:04")
	sunsetTimeDifference := sunsetTime.Sub(sunriseTime)

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
		LocationText: utils.GetLocationSuffix(locationName),
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
		Sunrise:       sunriseLocal,
		Sunset:        sunsetLocal,
		SunsetDifference: fmt.Sprintf("%02ds %02dd",
			int(sunsetTimeDifference.Hours()),
			int(sunsetTimeDifference.Minutes())%60),
		AllCities: allCities,
	}, nil
}

func GetTimezoneBySlug(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	slug := parts[len(parts)-1]

	var (
		found        models.City
		foundCountry models.Country
		locationName string
		country      *string
		typeOpt      *string
		latitude     float64
		longitude    float64
		allCities    []models.City
	)

	for _, city := range cities {
		if city.Slug == slug {
			found = city
			locationName = city.Name
			country = &city.Country
			latitude = city.Latitude
			longitude = city.Longitude
			for _, c := range cities {
				if c.Country == city.Country {
					randomType := utils.GetRandomType()
					c.Type = &randomType
					allCities = append(allCities, c)
				}
			}
			break
		}
	}

	if found.Slug == "" {
		for _, countryData := range countries {
			if countryData.Slug == slug {
				foundCountry = countryData
				locationName = countryData.Name
				country = &foundCountry.Name
				latitude = foundCountry.Latitude
				longitude = foundCountry.Longitude
				for _, c := range cities {
					if c.Country == foundCountry.Name {
						randomType := utils.GetRandomType()
						c.Type = &randomType
						allCities = append(allCities, c)
					}
				}
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

	sort.Slice(allCities, func(i, j int) bool {
		return allCities[i].Name < allCities[j].Name
	})

	response, err := getTimeForLocation(timezone, locationName, slug, country, typeOpt, latitude, longitude, allCities)
	if err != nil {
		http.Error(w, "Saat hesaplanırken hata oluştu", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

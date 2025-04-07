package models

import "time"

type Country struct {
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	Timezone  string  `json:"timezone"`
	Type      *string `json:"type,omitempty"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type City struct {
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	Timezone  string  `json:"timezone"`
	Country   string  `json:"country"`
	Type      *string `json:"type,omitempty"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Offset    float64 `json:"offset"`
}

type Response struct {
	Year              int        `json:"year"`
	Month             int        `json:"month"`
	Day               int        `json:"day"`
	Hour              int        `json:"hour"`
	Minute            int        `json:"minute"`
	Seconds           int        `json:"seconds"`
	MilliSeconds      int        `json:"milliSeconds"`
	DateTime          string     `json:"dateTime"`
	Date              string     `json:"date"`
	Time              string     `json:"time"`
	DayOfWeek         string     `json:"dayOfWeek"`
	DstActive         bool       `json:"dstActive"`
	LocationText      string     `json:"locationText"`
	Timezone          City       `json:"timezone"`
	PopulerCities     []CityTime `json:"populerCities"`
	Sunrise           string     `json:"sunrise"`
	Sunset            string     `json:"sunset"`
	SunsetDifference  string     `json:"sunsetDifference"`
	AllCities         []City     `json:"allCities"`
	NoonTime          string     `json:"noonTime"`
	NoonDifferenceMin int        `json:"noonDifferenceMin"`
}

type CityTime struct {
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Hour     int    `json:"hour"`
	Minute   int    `json:"minute"`
	DateTime string `json:"dateTime"`
	Selected bool   `json:"selected"`
}

type HourPair struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type CityWithTime struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Timezone string `json:"timezone"`
	Country  string `json:"country"`
	DateTime string `json:"dateTime"`
}

type CompareResponse struct {
	From           CityWithTime  `json:"from"`
	To             CityWithTime  `json:"to"`
	DifferenceText string        `json:"differenceText"`
	HourTable      []HourPair    `json:"hourTable"`
	Diff           time.Duration `json:"diff"`
	FormattedDiff  string        `json:"formattedDiff"`
}

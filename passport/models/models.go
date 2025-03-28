package models

type PassportCountries struct {
	Country string `json:"country"`
	Value   string `json:"value"`
	BlogUrl string `json:"blogUrl,omitempty"`
}

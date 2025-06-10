package schoolscores

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"sorgulat-api/db"
	"sorgulat-api/school-scores/models"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func fetchUniversitiesFromMongo() []models.University {
	var universities []models.University
	cursor, err := db.UniversityCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Println("MongoDB'den veri çekilirken hata:", err)
		return universities
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var uni models.University
		if err := cursor.Decode(&uni); err == nil {
			universities = append(universities, uni)
		}
	}
	return universities
}

func generateUniversityID(name, city, district string) string {
	combined := normalizeTurkish(name + "|" + city + "|" + district)
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])[:16]
}

func normalizeTurkish(s string) string {
	replacer := strings.NewReplacer(
		"ç", "c", "Ç", "c",
		"ğ", "g", "Ğ", "g",
		"ı", "i", "I", "i",
		"İ", "i",
		"ö", "o", "Ö", "o",
		"ş", "s", "Ş", "s",
		"ü", "u", "Ü", "u",
	)
	s = strings.ToLower(s)
	s = replacer.Replace(s)
	return s
}

func containsCI(slice []string, val string) bool {
	normVal := normalizeTurkish(val)
	for _, item := range slice {
		if normalizeTurkish(item) == normVal {
			return true
		}
	}
	return false
}

func slugify(s string) string {
	s = normalizeTurkish(s)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "&", "ve")
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, "--", "-")
	return s
}

func slugifyList(list []string) []string {
	var result []string
	for _, item := range list {
		if item != "" {
			result = append(result, slugify(item))
		}
	}
	return result
}

func GetUniversities(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	universities := fetchUniversitiesFromMongo()

	filterCities := strings.Split(query.Get("city"), ",")
	filterDistricts := strings.Split(query.Get("district"), ",")
	filterDegreeLevels := strings.Split(query.Get("degree_level"), ",")
	filterScoreTypes := strings.Split(query.Get("score_type"), ",")
	filterUniversityTypes := strings.Split(query.Get("university_type"), ",")

	yearStrings := strings.Split(query.Get("year"), ",")
	var filterYears []int
	for _, y := range yearStrings {
		if y == "" {
			continue
		}
		if parsedYear, err := strconv.Atoi(y); err == nil {
			filterYears = append(filterYears, parsedYear)
		}
	}

	baseMin, _ := strconv.ParseFloat(query.Get("base_score_min"), 64)
	baseMax, _ := strconv.ParseFloat(query.Get("base_score_max"), 64)

	quotaMin, _ := strconv.Atoi(query.Get("quota_min"))
	quotaMax, _ := strconv.Atoi(query.Get("quota_max"))

	filterDepartmentSlugs := slugifyList(strings.Split(query.Get("department"), ","))
	filterUniversitySlugs := slugifyList(strings.Split(query.Get("university"), ","))
	filterLanguages := strings.Split(query.Get("language"), ",")
	filterEducationTypes := strings.Split(query.Get("education_type"), ",")

	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	applyPagination := false
	var page, limit int

	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
		if limit < 1 {
			limit = 10
		}
		applyPagination = true
	} else if pageStr != "" {
		limit = 10
		applyPagination = true
	}

	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}
	} else {
		page = 1
	}

	var filtered []models.University

	for _, uni := range universities {
		if len(filterCities) > 0 && filterCities[0] != "" && !containsCI(filterCities, uni.City) {
			continue
		}
		if len(filterDistricts) > 0 && filterDistricts[0] != "" && !containsCI(filterDistricts, uni.District) {
			continue
		}
		if len(filterUniversitySlugs) > 0 && filterUniversitySlugs[0] != "" && !containsCI(filterUniversitySlugs, slugify(uni.Name)) {
			continue
		}
		if len(filterUniversityTypes) > 0 && filterUniversityTypes[0] != "" && !containsCI(filterUniversityTypes, string(uni.UniversityType)) {
			continue
		}

		var matchingDepartments []models.Department
		for _, dept := range uni.Departments {
			if len(filterDepartmentSlugs) > 0 && filterDepartmentSlugs[0] != "" && !containsCI(filterDepartmentSlugs, slugify(dept.Name)) {
				continue
			}
			if len(filterLanguages) > 0 && filterLanguages[0] != "" && !containsCI(filterLanguages, dept.Language) {
				continue
			}
			if len(filterEducationTypes) > 0 && filterEducationTypes[0] != "" && !containsCI(filterEducationTypes, string(dept.EducationType)) {
				continue
			}
			if len(filterDegreeLevels) > 0 && filterDegreeLevels[0] != "" && !containsCI(filterDegreeLevels, string(dept.DegreeLevel)) {
				continue
			}
			if len(filterScoreTypes) > 0 && filterScoreTypes[0] != "" && !containsCI(filterScoreTypes, string(dept.ScoreType)) {
				continue
			}

			var matchingYearlyData []models.YearlyData
			for _, yd := range dept.YearlyData {
				if len(filterYears) > 0 {
					matched := false
					for _, y := range filterYears {
						if yd.Year == y {
							matched = true
							break
						}
					}
					if !matched {
						continue
					}
				}
				if baseMin != 0 && yd.BaseScore < baseMin {
					continue
				}
				if baseMax != 0 && yd.BaseScore > baseMax {
					continue
				}
				if quotaMin != 0 && yd.Quota < quotaMin {
					continue
				}
				if quotaMax != 0 && yd.Quota > quotaMax {
					continue
				}
				matchingYearlyData = append(matchingYearlyData, yd)
			}

			if len(matchingYearlyData) > 0 {
				dept.YearlyData = matchingYearlyData
				matchingDepartments = append(matchingDepartments, dept)
			}
		}

		if len(matchingDepartments) > 0 {
			newUni := uni
			newUni.Slug = slugify(newUni.Name)
			newUni.ID = generateUniversityID(newUni.Name, newUni.City, newUni.District)
			var updatedDepartments []models.Department
			for _, dept := range matchingDepartments {
				dept.Slug = slugify(dept.Name)
				updatedDepartments = append(updatedDepartments, dept)
			}
			newUni.Departments = updatedDepartments
			filtered = append(filtered, newUni)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if applyPagination {
		start := (page - 1) * limit
		end := start + limit
		if start > len(filtered) {
			start = len(filtered)
		}
		if end > len(filtered) {
			end = len(filtered)
		}
		json.NewEncoder(w).Encode(filtered[start:end])
	} else {
		json.NewEncoder(w).Encode(filtered)
	}
}

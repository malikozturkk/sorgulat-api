package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

func GetRandomIndexes(length int) map[int]bool {
	rand.Seed(time.Now().UnixNano())
	count := max(2, length/1)
	indexes := make(map[int]bool)

	for len(indexes) < count {
		indexes[rand.Intn(length)] = true
	}
	return indexes
}

func GetRandomType() string {
	weightedTypes := []string{
		"5xl", "5xl", "5xl",
		"3xl", "3xl", "3xl", "3xl", "3xl", "3xl",
		"xl", "xl", "xl", "xl",
		"base",
	}
	rand.Seed(time.Now().UnixNano())
	return weightedTypes[rand.Intn(len(weightedTypes))]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func LoadData[T any](dataType string) []T {
	filePath := "timezones/" + dataType + ".json"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("%s dosyası açılamadı: %v", filePath, err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("%s dosyası okunamadı: %v", filePath, err)
	}

	var data []T
	if err := json.Unmarshal(byteValue, &data); err != nil {
		log.Fatalf("JSON parse hatası: %v", err)
	}
	return data
}

func GetLocationSuffix(name string) string {
	vowels := "aeıioöuü"
	deVowels := "eiöü"

	for i := len(name) - 1; i >= 0; i-- {
		char := unicode.ToLower(rune(name[i]))
		if strings.ContainsRune(vowels, char) {
			if strings.ContainsRune(deVowels, char) {
				return "de"
			}
			return "da"
		}
	}
	return "da"
}

package utils

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

const (
	earthRadiusKM = 6371
)

// TrimSpace trim space
func TrimSpace(i interface{}, depth int) {
	e := reflect.ValueOf(i).Elem()
	for i := 0; i < e.NumField(); i++ {
		if depth < 3 && e.Type().Field(i).Type.Kind() == reflect.Struct {
			depth++
			TrimSpace(e.Field(i).Addr().Interface(), depth)
		}

		if e.Type().Field(i).Type.Kind() != reflect.String {
			continue
		}

		value := e.Field(i).Interface().(string)
		e.Field(i).SetString(strings.TrimSpace(value))
	}
}

// isValidCitizenID is valid citizenId
func isValidCitizenID(citizenID string) bool {
	if !regexCitizenID.MatchString(citizenID) {
		return false
	}

	sum, row := 0, 13
	citizenIDRune := []rune(citizenID)
	for _, n := range string(citizenIDRune) {
		number, _ := strconv.Atoi(string(n))
		sum += number * row
		row--

		if row == 1 {
			break
		}
	}

	citizenIDInt, _ := strconv.Atoi(citizenID)
	result := (11 - (int(sum) % 11)) % 10

	return (citizenIDInt % 10) == result
}

// DeleteMapping delete map
func DeleteMapping(s map[string]interface{}, fields []string) {
	for _, field := range fields {
		delete(s, field)
	}

}

// ReplaceMapping replace map
func ReplaceMapping(s map[string]interface{}, replaces, olds []string) {
	for i, replace := range replaces {
		s[replace] = s[olds[i]]
	}

	DeleteMapping(s, olds)
}

// DistanceKiloMeters distance kilo meter between 2 poin coordinate
// use harversine fomula ref: https://en.wikipedia.org/wiki/Haversine_formula
func DistanceKiloMeters(latFirst float64, lngFirst float64, latSecond float64, lngSecond float64) float64 {
	radianLatFirst := latFirst * math.Pi / 180
	radianLngFirst := lngFirst * math.Pi / 180
	radianLatSecond := latSecond * math.Pi / 180
	radianLngSecond := lngSecond * math.Pi / 180

	diffLat := radianLatSecond - radianLatFirst
	diffLon := radianLngSecond - radianLngFirst
	h := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(radianLatFirst)*math.Cos(radianLatSecond)*math.Pow(math.Sin(diffLon/2), 2)
	distance := 2 * math.Asin(math.Sqrt(h))
	return earthRadiusKM * distance
}

// GetBindDataWithoutLast4Digit get binding data without last 4 digit
func GetBindDataWithoutLast4Digit(identificationNumber string) (binding string) {
	for i := 0; i < len(identificationNumber)-4; i++ {
		binding += "X"
	}

	return fmt.Sprint(binding, identificationNumber[len(identificationNumber)-4:])
}

// IntersectionString intersection from string array
func IntersectionString(a, b []string) []string {
	r := []string{}
	for _, i := range a {
		for _, j := range b {
			if i == j {
				r = append(r, i)
			}
		}
	}

	return r
}

package model

import (
	"fmt"
	"reflect"
	"sort"
	"time"
)

func (s *StravaOauthResponse) IsExpired() bool {
	return time.Now().Unix() > int64(s.ExpiresAt)
}

func (apiActivity *StravaActivity) CompareStravaData(dbActivity *StravaActivity) bool {
	return apiActivity.Distance == dbActivity.Distance &&
		apiActivity.MovingTime == dbActivity.MovingTime &&
		apiActivity.TotalElevationGain == dbActivity.TotalElevationGain &&
		apiActivity.Type == dbActivity.Type &&
		apiActivity.SportType == dbActivity.SportType &&
		apiActivity.WorkoutType == dbActivity.WorkoutType
}
func (athleteData *AthleteData) FormatDuration() string {
	hours := int(athleteData.AverageTime) / 3600
	minutes := (int(athleteData.AverageTime) % 3600) / 60
	seconds := int(athleteData.AverageTime) % 60
	if hours > 0 {
		return fmt.Sprintf("%dh:%02dm:%02ds", hours, minutes, seconds)
	}
	return fmt.Sprintf("%02dm:%02ds", minutes, seconds)
}

func SortAthletesData(athletesData []AthleteData, sortField string) {
	sort.Slice(athletesData, func(i, j int) bool {
		valI := reflect.ValueOf((athletesData)[i]).FieldByName(sortField)
		valJ := reflect.ValueOf((athletesData)[j]).FieldByName(sortField)
		if !valI.IsValid() || !valJ.IsValid() {
			return false
		}
		switch valI.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return valI.Int() > valJ.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return valI.Uint() > valJ.Uint()
		case reflect.Float32, reflect.Float64:
			return valI.Float() > valJ.Float()
		case reflect.String:
			return valI.String() > valJ.String()
		default:
			return false
		}
	})
}

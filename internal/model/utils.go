package model

import (
	"fmt"
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

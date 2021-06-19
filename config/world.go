package config

import (
	"time"
)

var (
TimeStart = "2021-06-01T00:00:00+00:00"
GodsWarStart = 2705
ImperialYearStart = 2228
Genocide = 1028
Awakening = 2
DayTime = false

Months = map[int]map[string]interface{}{
	0:{"name": "phoenix", "season": "summer", "weather": "summer", "sunrise": 6, "sunset": 20},
	1:{"name": "dragon", "season": "summer", "weather": "summer", "sunrise": 5, "sunset": 21},
	2:{"name": "chimera", "season": "summer", "weather": "summer","sunrise": 6, "sunset": 20},
	3:{"name": "twilight", "season": "autumn", "weather": "fall", "sunrise": 6, "sunset": 20},
	4:{"name": "prairiefire", "season": "autumn", "weather": "fall", "sunrise": 6, "sunset": 20},
	5:{"name": "wildfire", "season": "autumn", "weather": "fall", "sunrise": 8, "sunset": 18},
	6:{"name": "midnight", "season": "winter", "weather": "fall", "sunrise": 9, "sunset": 18},
	7:{"name": "icedrake", "season": "winter", "weather": "fall", "sunrise": 9, "sunset": 18},
	8:{"name": "chrysalis", "season": "winter", "weather": "fall", "sunrise": 8, "sunset": 18},
	9:{"name": "dawn", "season": "spring", "weather": "spring", "sunrise": 7, "sunset": 19},
	10:{"name": "torrents", "season": "spring", "weather": "spring", "sunrise": 7, "sunset": 19},
	11:{"name": "blossoms", "season": "spring", "weather": "spring", "sunrise": 6, "sunset": 20},
}

Days = []string{"panur", "maaur", "ruvur", "dilur", "malkur", "arsur", "andur", "aalur", "tilur"}

)

func SyncCurrentTime() (diffHours int, diffDays int, diffMonths int, diffYears int){
	startTime, e := time.Parse(
		time.RFC3339,
		TimeStart)
	if e != nil {
		return
	}

	// Multiply by 4 because time moves at 4X the rate
	diff := time.Now().Sub(startTime)*4
	diffHours = int(diff.Minutes()/60)
	diffDays = diffHours/24
	diffMonths = diffDays/30
	diffYears = diffMonths/12
	return
}
package cmd

import (
	"fmt"
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
	"strconv"
	"strings"
)

func init() {
	addHandler(timeStat{},
		"Display time and world information",
		permissions.Player,
		"time")
}

type timeStat cmd

func (timeStat) process(s *state) {

	//TODO: Add some weather patterns later

	dayNight, sunAction, difference := CalculateTimeToNext()
	returnTime := fmt.Sprintf(`It is %s, the %s of the month of %s, season is %s
%s years since the Godswar
%s Tholmic Imperial Year, %s years after the Goblin Genocide Event.
%s years since the awakening began.
It is now %s. The suns will %s in %s hours.`,
		strings.Title(config.Days[objects.CurrentDay]),
		config.TextTiers[objects.DayOfMonth],
		strings.Title(config.Months[objects.CurrentMonth]["name"].(string)),
		strings.Title(config.Months[objects.CurrentMonth]["season"].(string)),
		strconv.Itoa(config.GodsWarStart+objects.YearPlus),
		strconv.Itoa(config.ImperialYearStart+objects.YearPlus),
		strconv.Itoa(config.Genocide+objects.YearPlus),
		strconv.Itoa(config.Awakening+objects.YearPlus),
		dayNight,
		sunAction,
		config.TextNumbers[difference],
	)
	s.msg.Actor.Send(text.White + returnTime + text.Reset)

	s.ok = true
}

func CalculateTimeToNext() (string, string, int) {
	//Calculate sunset
	if objects.CurrentHour >= config.Months[objects.CurrentMonth]["sunrise"].(int) &&
	objects.CurrentHour < config.Months[objects.CurrentMonth]["sunset"].(int){
		return "daytime", "set", config.Months[objects.CurrentMonth]["sunset"].(int) - objects.CurrentHour
	}else{
		if objects.CurrentHour < config.Months[objects.CurrentMonth]["sunrise"].(int) {
			return "nighttime", "rise", config.Months[objects.CurrentMonth]["sunrise"].(int) - objects.CurrentHour
		}else{
			return "nighttime", "rise", (24 - objects.CurrentHour) + config.Months[objects.CurrentMonth]["sunrise"].(int)
		}
	}
}
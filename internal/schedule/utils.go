package schedule

import (
	"log"
	"time"
)

func GetWeekDay(date string) string {
	weekDays := map[string]string{
		"Sunday":    "воскресенье",
		"Monday":    "понедельник",
		"Tuesday":   "вторник",
		"Wednesday": "среда",
		"Thursday":  "четверг",
		"Friday":    "пятница",
		"Saturday":  "суббота",
	}
	timeObject, err := time.Parse("02.01.2006", date)
	if err != nil {
		log.Println(err)
	}
	enDay := timeObject.Format("Monday")
	return weekDays[enDay]
}

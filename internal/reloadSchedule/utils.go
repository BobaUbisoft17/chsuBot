package reload

import (
	"sort"
	"time"

	"github.com/BobaUbisoft17/chsuBot/internal/schedule"
)

func collectLecture(sched []schedule.Lecture) (map[int]map[int][]schedule.Lecture, error) {
	var scheds = map[int]map[int][]schedule.Lecture{}
	for _, lesson := range sched {
		dateEvent := lesson.DateEvent
		for _, group := range lesson.Groups {
			timestamp, err := stringToTimestamp(dateEvent)
			if err != nil {
				return scheds, err
			}
			if _, ok := scheds[group.ID]; !ok {
				scheds[group.ID] = make(map[int][]schedule.Lecture)
			}
			if _, ok := scheds[group.ID][timestamp]; !ok {
				scheds[group.ID][timestamp] = []schedule.Lecture{}
			}
			scheds[group.ID][timestamp] = append(scheds[group.ID][timestamp], lesson)
		}
	}
	return scheds, nil
}

func GetKeys[T comparable, N any](v map[T]N) []T {
	keys := make([]T, 0, len(v))
	for key := range v {
		keys = append(keys, key)
	}
	return keys
}

func splitSchedule(schedules map[int][]schedule.Lecture) (string, string) {
	if len(schedules) == 2 {
		return schedule.New(schedules[0]).Render(), schedule.New(schedules[1]).Render()
	} else {
		timestamps := GetKeys(schedules)
		sort.Ints(timestamps)
		//if now time - time of 00:00 today < hours * minutes * seconds => time is today
		if int(time.Now().Unix())-timestamps[0] < 24*60*60 {
			return schedule.New(schedules[0]).Render(), schedule.New([]schedule.Lecture{}).Render()
		} else {
			return schedule.New([]schedule.Lecture{}).Render(), schedule.New(schedules[0]).Render()
		}
	}
}

func stringToTimestamp(date string) (int, error) {
	timeObject, err := time.Parse("02.01.2006", date)
	if err != nil {
		return 0, err
	}
	return int(timeObject.Unix()), nil
}

package reload

import (
	"sort"
	"time"

	"github.com/BobaUbisoft17/chsuBot/internal/schedule"
	"github.com/BobaUbisoft17/chsuBot/pkg"
)

func collectLecture(sched []schedule.Lecture) (map[int]map[int][]schedule.Lecture, error) {
	var scheds = map[int]map[int][]schedule.Lecture{}
	for _, lesson := range sched {
		dateEvent := lesson.DateEvent
		for _, group := range lesson.Groups {
			timestamp, err := pkg.StringToTimestamp(dateEvent)
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

func splitSchedule(schedules map[int][]schedule.Lecture) ([]schedule.Lecture, []schedule.Lecture) {
	datesInTimestamp := pkg.GetKeys(schedules)
	if len(schedules) == 2 {
		sort.Ints(datesInTimestamp)
		return schedules[datesInTimestamp[0]], schedules[datesInTimestamp[1]]
	}
	today := time.Unix(int64(datesInTimestamp[0]), 0)
	durationFromToNow := time.Since(today)
	if durationFromToNow < time.Hour*24 && durationFromToNow >= 0 {
		return schedules[datesInTimestamp[0]], []schedule.Lecture{}
	}
	return []schedule.Lecture{}, schedules[datesInTimestamp[0]]
}

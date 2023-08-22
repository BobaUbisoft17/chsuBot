package reload

import (
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
	if len(schedules) == 2 {
		return schedules[0], schedules[1]
	}
	datesInTimestamp := pkg.GetKeys(schedules)
	todayTimestamp := min(datesInTimestamp)
	today := time.Unix(int64(todayTimestamp), 0)

	if time.Since(today) < time.Hour*24 {
		return schedules[0], []schedule.Lecture{}
	}
	return []schedule.Lecture{}, schedules[0]
}

func min(num []int) int {
	var minNum int
	for i := 0; i < len(num); i++ {
		if minNum > num[i] {
			minNum = num[i]
		}
	}
	return minNum
}

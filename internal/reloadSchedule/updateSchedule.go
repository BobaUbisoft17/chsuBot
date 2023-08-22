package reload

import (
	"sync"
	"time"

	"github.com/BobaUbisoft17/chsuBot/internal/schedule"
)

type chsuAPI interface {
	All() ([]schedule.Lecture, error)
}

type groupStorage interface {
	UpdateSchedule(todaySchedule, tomorrowSchedule string, groupID int)
	UnusedID(ID []int) []int
}

type logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

type Reloader struct {
	api     chsuAPI
	groupDb groupStorage
	logger  logger
}

func NewReloader(api chsuAPI, db groupStorage, logger logger) *Reloader {
	return &Reloader{
		api:     api,
		groupDb: db,
		logger:  logger,
	}
}

func (r *Reloader) ReloadSchedule(waitingTimeSeconds int) {
	time.Sleep(time.Duration(waitingTimeSeconds) * time.Second)
	r.logger.Info("Schedule update process started")
	var wg sync.WaitGroup
	unSortedSchedule, err := r.api.All()
	if err != nil {
		r.logger.Errorf("%w", err)
		return
	}
	sortedScheduleByIDs, err := collectLecture(unSortedSchedule)
	if err != nil {
		r.logger.Errorf("%w", err)
	}
	for key := range sortedScheduleByIDs {
		wg.Add(1)
		go func(id int, schedules map[int][]schedule.Lecture) {
			defer wg.Done()
			today, tomorrow := splitSchedule(schedules)

			r.groupDb.UpdateSchedule(today, tomorrow, id)
		}(key, sortedScheduleByIDs[key])
	}
	wg.Wait()
	r.addScheduleMissingGroups(GetKeys(sortedScheduleByIDs))
	r.logger.Info("Schedule updated succesfully")
}

func (r *Reloader) addScheduleMissingGroups(keys []int) {
	var wg sync.WaitGroup
	undefindTimetable := "Расписание не найдено"
	for _, id := range r.groupDb.UnusedID(keys) {
		wg.Add(1)
		go func(id int, timetable string) {
			defer wg.Done()
			r.groupDb.UpdateSchedule(timetable, timetable, id)
		}(id, undefindTimetable)
	}
	wg.Wait()
}

func (r *Reloader) Start() {
	for {
		nowTime := time.Now()
		nextUpdateTimeMinutes := (nowTime.Hour()/6+1)*6*60 + 1
		nowTimeMinutes := nowTime.Hour()*60 + nowTime.Minute()
		waitingTimeSeconds := (nextUpdateTimeMinutes - nowTimeMinutes) * 60
		r.ReloadSchedule(waitingTimeSeconds)
	}
}

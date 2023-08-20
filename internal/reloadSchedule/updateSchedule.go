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
	api    chsuAPI
	db     groupStorage
	logger logger
}

func NewReloader(api chsuAPI, db groupStorage, logger logger) *Reloader {
	return &Reloader{
		api:    api,
		db:     db,
		logger: logger,
	}
}

func (r *Reloader) ParseSchedule(waitTime int) {
	time.Sleep(time.Duration(waitTime) * time.Second)
	r.logger.Info("Начался процесс обновления расписания")
	var wg sync.WaitGroup
	unSortedSchedule, err := r.api.All()
	if err != nil {
		r.logger.Errorf("При запросе на получение расписания произошла ошибка: %s", err)
		return
	}
	sortedScheduleByIDs, err := collectLecture(unSortedSchedule)
	if err != nil {
		r.logger.Errorf("Ошибка перевода даты в timestamp:%s", err)
	}
	for key := range sortedScheduleByIDs {
		wg.Add(1)
		go func(ID int, schedules map[int][]schedule.Lecture) {
			defer wg.Done()
			today, tomorrow := splitSchedule(schedules)

			r.db.UpdateSchedule(today, tomorrow, ID)
		}(key, sortedScheduleByIDs[key])
	}
	wg.Wait()
	r.addScheduleMissingGroups(GetKeys(sortedScheduleByIDs))
	r.logger.Info("Расписание успешно обновлено")
}

func (r *Reloader) addScheduleMissingGroups(keys []int) {
	var wg sync.WaitGroup
	undefindTimetable := "Расписание не найдено"
	for _, ID := range r.db.UnusedID(keys) {
		wg.Add(1)
		go func(ID int, timetable string) {
			defer wg.Done()
			r.db.UpdateSchedule(timetable, timetable, ID)
		}(ID, undefindTimetable)
	}
	wg.Wait()
}

func (r *Reloader) UpdateSchedule() {
	for {
		nowTime := time.Now()

		waitingTime := ((nowTime.Hour()/6+1)*6*60 + 1 - (nowTime.Hour()*60 + nowTime.Minute())) * 60
		r.ParseSchedule(waitingTime)
	}
}

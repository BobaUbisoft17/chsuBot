package reload

import (
	"sync"
	"time"

	"github.com/BobaUbisoft17/chsuBot/internal/schedule"
	"github.com/BobaUbisoft17/chsuBot/pkg"
	"github.com/BobaUbisoft17/chsuBot/pkg/logging"
)

type chsuAPI interface {
	All() ([]schedule.Lecture, error)
}

type groupStorage interface {
	UpdateSchedule(todaySchedule, tomorrowSchedule []schedule.Lecture, groupID int) error
	UnusedID(ID []int) ([]int, error)
}

type Reloader struct {
	api     chsuAPI
	groupDb groupStorage
	logger  *logging.Logger
}

func NewReloader(api chsuAPI, db groupStorage, logger *logging.Logger) *Reloader {
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
		r.logger.Errorf("%v", err)
		return
	}
	sortedScheduleByIDs, err := collectLecture(unSortedSchedule)
	if err != nil {
		r.logger.Errorf("%v", err)
		return
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
	r.addScheduleMissingGroups(pkg.GetKeys(sortedScheduleByIDs))
	r.logger.Info("Schedule updated succesfully")
}

func (r *Reloader) addScheduleMissingGroups(keys []int) {
	var wg sync.WaitGroup
	ids, err := r.groupDb.UnusedID(keys)
	if err != nil {
		r.logger.Errorf("%v", err)
		return
	}
	for _, id := range ids {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			r.groupDb.UpdateSchedule([]schedule.Lecture{}, []schedule.Lecture{}, id)
		}(id)
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

package reload

import (
	"sort"
	"time"

	"github.com/BobaUbisoft17/chsuBot/internal/schedule"
	"github.com/BobaUbisoft17/chsuBot/pkg"
	"github.com/BobaUbisoft17/chsuBot/pkg/logging"
)

func NewReloader(api chsuAPI, db groupStorage, logger *logging.Logger) *Reloader {
	return &Reloader{
		api:     api,
		groupDb: db,
		logger:  logger,
	}
}

func (r *Reloader) Start() {
	for {
		waitingTime := r.getWaitingTime()
		time.Sleep(time.Duration(waitingTime) * time.Second)
		r.ReloadSchedule()
	}
}

func (r *Reloader) getWaitingTime() int {
	nowTime := time.Now()
	nextUpdateTimeMinutes := (nowTime.Hour()/6+1)*6*60 + 1
	nowTimeMinutes := nowTime.Hour()*60 + nowTime.Minute()
	waitingTimeSeconds := (nextUpdateTimeMinutes - nowTimeMinutes) * 60
	return waitingTimeSeconds
}

func (r *Reloader) ReloadSchedule() error {
	r.logger.Info("Schedule update process started")
	schedule, err := r.getSchedule()
	if err != nil {
		r.logger.Errorf("%v", err)
		return err
	}

	r.groupDb.UpdateSchedule(schedule)

	r.logger.Info("Schedule updated succesfully")
	return nil
}

func (r *Reloader) getSchedule() (map[int][2]string, error) {
	unsortedSchedule, err := r.api.All()
	if err != nil {
		r.logger.Errorf("%v", err)
		return nil, err
	}
	sortedSchedule, err := r.processingSchedule(unsortedSchedule)
	if err != nil {
		r.logger.Errorf("%v", err)
		return nil, err
	}
	return sortedSchedule, nil
}

func (r *Reloader) processingSchedule(lectures []schedule.Lecture) (map[int][2]string, error) {
	groupsScedule, err := r.collectLectureByGroupsId(lectures)
	if err != nil {
		r.logger.Errorf("%v", err)
		return nil, err
	}

	groupsSceduleByDays := map[int][2][]schedule.Lecture{}
	for group := range groupsScedule {
		todaySchedule, tomorrowSchedule := r.splitScheduleByDays(groupsScedule[group])
		groupsSceduleByDays[group] = [2][]schedule.Lecture{todaySchedule, tomorrowSchedule}
	}

	if err := r.addNotStudyingGroupsSchedule(groupsSceduleByDays); err != nil {
		r.logger.Errorf("%v", err)
		return nil, err
	}

	renderedSchedule := r.renderSchedule(groupsSceduleByDays)

	return renderedSchedule, nil
}

func (r *Reloader) collectLectureByGroupsId(lectures []schedule.Lecture) (map[int]map[int][]schedule.Lecture, error) {
	var groupLectures = map[int]map[int][]schedule.Lecture{}
	for _, lecture := range lectures {
		if err := r.addLectureToGroups(groupLectures, &lecture); err != nil {
			r.logger.Errorf("%v", err)
			return nil, err
		}
	}
	return groupLectures, nil
}

func (r *Reloader) addLectureToGroups(
	groupsLectures map[int]map[int][]schedule.Lecture,
	lecture *schedule.Lecture,
) error {
	lectureTimestamp, err := pkg.StringToTimestamp(lecture.DateEvent)
	if err != nil {
		r.logger.Errorf("%v", err)
		return err
	}
	for _, group := range lecture.Groups {
		if _, ok := (groupsLectures)[group.ID]; !ok {
			groupsLectures[group.ID] = make(map[int][]schedule.Lecture)
		}
		if _, ok := groupsLectures[group.ID][lectureTimestamp]; !ok {
			groupsLectures[group.ID][lectureTimestamp] = nil
		}
		groupsLectures[group.ID][lectureTimestamp] = append(
			groupsLectures[group.ID][lectureTimestamp],
			*lecture,
		)
	}
	return nil
}

func (r *Reloader) splitScheduleByDays(schedule map[int][]schedule.Lecture) ([]schedule.Lecture, []schedule.Lecture) {
	dates := pkg.GetKeys(schedule)
	if len(dates) == 2 {
		sort.Ints(dates)
		return schedule[dates[0]], schedule[dates[1]]
	}
	today := time.Unix(int64(dates[0]), 0).Add(-3 * time.Hour)
	durationFromToNow := time.Since(today)
	if durationFromToNow < time.Hour*24 && durationFromToNow >= 0 {
		return schedule[dates[0]], nil
	}
	return nil, schedule[dates[0]]
}

func (r *Reloader) addNotStudyingGroupsSchedule(groupsSchedule map[int][2][]schedule.Lecture) error {
	keys := pkg.GetKeys(groupsSchedule)
	ids, err := r.groupDb.UnusedID(keys)
	if err != nil {
		r.logger.Errorf("%v", err)
		return err
	}

	for _, id := range ids {
		groupsSchedule[id] = [2][]schedule.Lecture{{}, {}}
	}
	return nil
}

func (r *Reloader) renderSchedule(unrenderedSchedule map[int][2][]schedule.Lecture) map[int][2]string {
	renderedSchedule := make(map[int][2]string, len(unrenderedSchedule))
	for groupId := range unrenderedSchedule {
		todaySchedule := unrenderedSchedule[groupId][0]
		tomorrowSchedule := unrenderedSchedule[groupId][1]
		renderedSchedule[groupId] = [2]string{
			schedule.New(todaySchedule).Render(),
			schedule.New(tomorrowSchedule).Render(),
		}
	}
	return renderedSchedule
}

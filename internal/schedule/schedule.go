package schedule

import (
	"fmt"
	"strings"
)

func New(data []Lecture) *Schedule {
	return &Schedule{
		Data: data,
	}
}

func (s *Schedule) Render() string {
	if len(s.Data) != 0 {
		for i := range s.Data {
			nowLecture := s.Data[i]
			if i == 0 {
				s.addTitle()
			}
			s.timeDuration(nowLecture)
			s.lessonName(nowLecture)
			s.lecturerName(nowLecture)
			s.location(nowLecture)
		}
		return s.schedule
	}
	return "Расписание не найдено"
}

func (s *Schedule) addTitle() {
	eventDate := s.Data[0].DateEvent
	s.schedule += fmt.Sprintf(
		"*Расписание на %s - %s*\n\n",
		eventDate,
		GetWeekDay(eventDate),
	)
}

func (s *Schedule) timeDuration(lesson Lecture) {
	s.schedule += fmt.Sprintf("⌚ %s-%s\n", lesson.StartTime, lesson.EndTime)
}

func (s *Schedule) lessonName(lesson Lecture) {
	var lecture string
	if lesson.Abbrlessontype == "" {
		lecture = fmt.Sprintf("%s\n", lesson.Discipline.Title)
	} else {
		lecture = fmt.Sprintf("%s. %s\n", lesson.Abbrlessontype, lesson.Discipline.Title)
	}
	s.schedule += fmt.Sprintf("🏫 %s", lecture)
}

func (s *Schedule) lecturerName(lesson Lecture) {
	var lecturers string
	for _, lecturer := range lesson.Lecturers {
		lecturers += fmt.Sprintf("%s ", lecturer.ShortName)
	}
	s.schedule += fmt.Sprintf("🧑 %s\n", lecturers)
}

func (s *Schedule) location(lesson Lecture) {
	var adress string
	if lesson.OnlineEvent == 0 {
		adress = "Онлайн"
	} else if lesson.Auditory.Title == "" {
		adress = "-/-"
	} else {
		adress = fmt.Sprintf("%s, %s",
			lesson.Auditory.Title,
			strings.ToLower(lesson.Build.Title),
		)
	}
	s.schedule += fmt.Sprintf("🏢 %s\n\n", adress)
}

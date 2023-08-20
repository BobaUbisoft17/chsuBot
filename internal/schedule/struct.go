package schedule

type Schedule struct {
	Data     []Lecture
	schedule string
}

type Lecture struct {
	ID         int    `json:"id"`
	DateEvent  string `json:"dateEvent"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	Discipline struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
	} `json:"discipline"`
	Groups []Groups `json:"groups"`
	Build  struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
	} `json:"build"`
	Auditory struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
	} `json:"auditory"`
	Lecturers []struct {
		ID         int64  `json:"id"`
		LastName   string `json:"lastName"`
		FirstName  string `json:"firstName"`
		MiddleName string `json:"middleName"`
		ShortName  string `json:"shortName"`
		Fio        string `json:"fio"`
	} `json:"lecturers"`
	Abbrlessontype string      `json:"abbrlessontype"`
	Lessontype     string      `json:"lessontype"`
	Week           int         `json:"week"`
	Weekday        int         `json:"weekday"`
	WeekType       string      `json:"weekType"`
	OnlineEvent    interface{} `json:"onlineEvent"`
	Online         int         `json:"online"`
}

type Groups struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

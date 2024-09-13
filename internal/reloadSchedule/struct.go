package reload

import (
	"github.com/BobaUbisoft17/chsuBot/internal/schedule"
	"github.com/BobaUbisoft17/chsuBot/pkg/logging"
)

type chsuAPI interface {
	All() ([]schedule.Lecture, error)
}

type groupStorage interface {
	UpdateSchedule(map[int][2]string) error
	UnusedID(ID []int) ([]int, error)
}

type Reloader struct {
	api     chsuAPI
	groupDb groupStorage
	logger  *logging.Logger
}

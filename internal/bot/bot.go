package bot

import (
	"github.com/BobaUbisoft17/chsuBot/internal/schedule"
	"github.com/NicoNex/echotron/v3"
)

type api interface {
	One(startDate, endDate string, groupId int) ([]schedule.Lecture, error)
}

type groupStorage interface {
	GetGroupNames() []string
	GetTodaySchedule(groupID int) string
	GetTomorrowSchedule(groupID int) string
	GroupId(string) int
	GroupNameIsCorrect(groupName string) bool
}

type userStorage interface {
	AddUser(userID int64)
	ChangeUserGroup(userID int64, groupName string)
	DeleteGroup(userID int64)
	GetUserGroup(userID int64) int
	IsUserHasGroup(userID int64) bool
	IsUserInDB(userID int64) bool
}

type logger interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
}

type stateFn func(*echotron.Update) stateFn

type nextFn func()

type bot struct {
	chatID    int64
	state     stateFn
	nextState nextFn
	group     int
	startDate string
	endDate   string
	echotron.API
	chsuAPI api
	groupDb groupStorage
	logger  logger
	token   string
	usersDb userStorage
}

func New(api api, groupDb groupStorage, userDb userStorage, logger logger, token string) *bot {
	return &bot{
		chsuAPI: api,
		groupDb: groupDb,
		usersDb: userDb,
		logger:  logger,
		token:   token,
	}
}

func (b *bot) newBot(chatID int64) echotron.Bot {
	b.chatID = chatID
	b.API = echotron.NewAPI(b.token)
	b.state = b.HandleMessage
	return b
}

func (b *bot) Update(update *echotron.Update) {
	if update.Message != nil || update.CallbackQuery != nil {
		b.state = b.state(update)
	}
}

func (b *bot) StartBot() {
	dsp := echotron.NewDispatcher(b.token, b.newBot)
	b.logger.Info(dsp.Poll())
}

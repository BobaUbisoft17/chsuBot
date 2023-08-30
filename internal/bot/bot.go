package bot

import (
	"time"

	"github.com/BobaUbisoft17/chsuBot/internal/database"
	"github.com/BobaUbisoft17/chsuBot/internal/schedule"
	"github.com/BobaUbisoft17/chsuBot/pkg/logging"
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
	GroupsStartsWith(firstSymbol string) []database.GroupInfo
}

type userStorage interface {
	AddUser(userID int64)
	ChangeUserGroup(userID int64, groupID int)
	DeleteGroup(userID int64)
	DeleteUser(userID int64)
	GetUserGroup(userID int64) int
	GetUsersId() []int
	IsUserHasGroup(userID int64) bool
	IsUserInDB(userID int64) bool
}

type stateFn func(*echotron.Update) stateFn

type nextFn func()

type bot struct {
	chatID     int64
	state      stateFn
	nextFn     nextFn
	previousFn nextFn
	group      int
	startDate  time.Time
	endDate    time.Time
	postText   string
	echotron.API
	chsuAPI api
	groupDb groupStorage
	logger  *logging.Logger
	token   string
	usersDb userStorage
}

func New(api api, groupDb groupStorage, userDb userStorage, logger *logging.Logger, token string) *bot {
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

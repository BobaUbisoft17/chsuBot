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

type usePackages struct {
	adminId int
	chsuAPI api
	groupDb groupStorage
	logger  *logging.Logger
	token   string
	usersDb userStorage
}

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
	usePackages *usePackages
}

func New(api api, groupDb groupStorage, userDb userStorage, logger *logging.Logger, token string, adminId int) *usePackages {
	return &usePackages{
		chsuAPI: api,
		groupDb: groupDb,
		usersDb: userDb,
		logger:  logger,
		token:   token,
		adminId: adminId,
	}
}

func (u *usePackages) newBot(chatID int64) echotron.Bot {
	bot := &bot{
		chatID:      chatID,
		API:         echotron.NewAPI(u.token),
		usePackages: u,
	}
	bot.state = bot.HandleMessage
	return bot
}

func (b *bot) Update(update *echotron.Update) {
	if update.Message != nil || update.CallbackQuery != nil {
		b.state = b.state(update)
	}
}

func (u *usePackages) StartBot() {
	dsp := echotron.NewDispatcher(u.token, u.newBot)
	u.logger.Info(dsp.Poll())
}

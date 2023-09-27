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
	GetGroupNames() ([]string, error)
	GetTodaySchedule(groupID int) (string, error)
	GetTomorrowSchedule(groupID int) (string, error)
	GroupNameIsCorrect(groupName string) (bool, error)
	GroupsStartsWith(firstSymbol string) ([]database.GroupInfo, error)
}

type userStorage interface {
	AddUser(userID int64) error
	ChangeUserGroup(userID int64, groupID int) error
	DeleteGroup(userID int64) error
	DeleteUser(userID int64) error
	GetUserGroup(userID int64) (int, error)
	GetUsersId() ([]int, error)
	IsUserHasGroup(userID int64) (bool, error)
	IsUserInDB(userID int64) (bool, error)
}

type stateFn func(*echotron.Update) stateFn

type nextFn func()

type usePackages struct {
	adminId   int
	chsuAPI   api
	groupDb   groupStorage
	logger    *logging.Logger
	token     string
	typeStart string
	usersDb   userStorage
	webhook   string
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

func New(api api, groupDb groupStorage, userDb userStorage, logger *logging.Logger, adminId int, token, typeStart, webhook string) *usePackages {
	return &usePackages{
		chsuAPI:   api,
		groupDb:   groupDb,
		usersDb:   userDb,
		logger:    logger,
		token:     token,
		adminId:   adminId,
		typeStart: typeStart,
		webhook:   webhook,
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
	if u.typeStart == "webhook" {
		u.logger.Info("Start on webhook")
		u.logger.Info(dsp.ListenWebhook(u.webhook))
	} else {
		u.logger.Info("Start on long polling")
		for {
			u.logger.Info(dsp.PollOptions(false, echotron.UpdateOptions{Timeout: 120}))
		}
	}
}

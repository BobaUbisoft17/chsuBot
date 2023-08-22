package bot

import (
	"sort"
	"strconv"
	"strings"
	"time"

	calendar "github.com/BobaUbisoft17/chsuBot/internal/bot/keyboard/inlineKeyboard"
	kb "github.com/BobaUbisoft17/chsuBot/internal/bot/keyboard/replyKeyboard"
	reload "github.com/BobaUbisoft17/chsuBot/internal/reloadSchedule"
	"github.com/BobaUbisoft17/chsuBot/internal/schedule"
	"github.com/NicoNex/echotron/v3"
)

func (b *bot) answer(answer string, keyboard echotron.ReplyMarkup) {
	messageOptions := getReplyMarkupMessageOptions(keyboard)
	b.SendMessage(answer, b.chatID, &messageOptions)
}

func getReplyMarkupMessageOptions(replyMarkup echotron.ReplyMarkup) echotron.MessageOptions {
	return echotron.MessageOptions{
		ReplyMarkup: replyMarkup,
		ParseMode:   "Markdown",
	}
}

func (b *bot) sendSchedule() {
	if b.endDate == "" {
		b.endDate = b.startDate
	}

	unParseSchedule, err := b.chsuAPI.One(b.startDate, b.endDate, b.group)
	if err != nil {
		b.logger.Info(err)
	}

	schedule, err := buildSchedule(unParseSchedule)
	if err != nil {
		b.logger.Errorf("%v", err)
		b.group = 0
		return
	}
	for i := 0; i < len(schedule); i++ {
		if i == len(schedule)-1 {
			b.answer(schedule[i], kb.ChooseDateMarkup())
		} else {
			b.answer(schedule[i], nil)
		}
	}
	b.startDate, b.endDate, b.group = "", "", 0
}

func buildSchedule(schedules []schedule.Lecture) ([]string, error) {
	var messages []string
	var message string
	if len(schedules) == 0 {
		return append(
			messages,
			schedule.New(schedules).Render(),
		), nil
	}

	sortedSchedule, err := sortScheduleByDate(schedules)
	if err != nil {
		return []string{}, err
	}
	keys := reload.GetKeys(sortedSchedule)
	sort.Ints(keys)
	for _, key := range keys {
		daySchedule := schedule.New(sortedSchedule[key]).Render()
		if len(message+daySchedule) < 4096 {
			message += daySchedule
		} else {
			messages, message = append(messages, message), daySchedule
		}
	}
	return append(messages, message), nil
}

func (b *bot) changeMonth(callback *echotron.CallbackQuery) {
	month, year, err := getDate(callback.Data)
	if err != nil {
		b.logger.Errorf("Ошибка получение даты: %v", err)
	}
	var markup echotron.InlineKeyboardMarkup
	if strings.Contains(callback.Data, "next") {
		markup = calendar.New(month, year).NextMonth()
	} else {
		markup = calendar.New(month, year).NextMonth()
	}
	message := echotron.NewMessageID(b.chatID, callback.Message.ID)
	opts := echotron.MessageReplyMarkup{ReplyMarkup: markup}
	b.EditMessageReplyMarkup(message, &opts)
}

func (b *bot) closeCalendarMarkup(callback *echotron.CallbackQuery) {
	message := echotron.NewMessageID(b.chatID, callback.Message.ID)
	b.EditMessageText("Вложение удалено", message, nil)
}

func (b *bot) getGroupKeyboard() {
	replyMarkup := kb.FirstPartGroups(b.groupDb.GetGroupNames())
	messageOptions := getReplyMarkupMessageOptions(replyMarkup)
	b.SendMessage("Введите назвние вашей группы", b.chatID, &messageOptions)
}

func (b *bot) manageCalendarKeyboard(callback *echotron.CallbackQuery) {
	switch {
	case strings.Contains(callback.Data, "next") || strings.Contains(callback.Data, "back"):
		b.changeMonth(callback)
	case callback.Data == "menu":
		b.closeCalendarMarkup(callback)
		b.state = b.HandleMessage
	}
}

func (b *bot) editGroupKeyboard(message string) {
	if message == "Назад" {
		b.state = b.HandleMessage
		b.answer(
			"Выберите дату",
			kb.ChooseDateMarkup(),
		)
	} else if message == "Дальше »" || message == "« Обратно" {
		b.answer(
			"Меняем клавиатуру",
			kb.GetKeyboardPart(message, b.groupDb.GetGroupNames()),
		)
	}
}

func sortScheduleByDate(timetable []schedule.Lecture) (map[int][]schedule.Lecture, error) {
	scheduleByDays := make(map[int][]schedule.Lecture)
	for _, lecture := range timetable {
		timestamp, err := stringToTimestamp(lecture.DateEvent)
		if err != nil {
			return map[int][]schedule.Lecture{}, err
		}
		scheduleByDays[timestamp] = append(scheduleByDays[timestamp], lecture)
	}
	return scheduleByDays, nil
}

func getDate(data string) (int, int, error) {
	date := strings.Split(strings.Split(data, " ")[1], ".")
	month, err := strconv.Atoi(date[0])
	if err != nil {
		return 0, 0, err
	}
	year, err := strconv.Atoi(date[1])
	if err != nil {
		return 0, 0, err
	}
	return month, year, nil
}

func parseDate(date string) (time.Time, error) {
	dateTime, err := time.Parse("02.01.2006", date)
	return dateTime, err
}

func (b *bot) validDuration() bool {
	startDate, _ := parseDate(b.startDate)
	endDate, _ := parseDate(b.endDate)
	duration := endDate.Sub(startDate).Hours() / 24
	return duration <= 31
}

func (b *bot) orderDateCheck() error {
	startDate, err := parseDate(b.startDate)
	if err != nil {
		return err
	}

	endDate, err := parseDate(b.endDate)
	if err != nil {
		return err
	}

	if startDate.After(endDate) {
		b.startDate, b.endDate = b.endDate, b.startDate
	}
	return nil
}

func stringToTimestamp(date string) (int, error) {
	timeObject, err := time.Parse("02.01.2006", date)
	return int(timeObject.Unix()), err
}

package bot

import (
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	ikb "github.com/BobaUbisoft17/chsuBot/internal/bot/keyboard/inlineKeyboard"
	kb "github.com/BobaUbisoft17/chsuBot/internal/bot/keyboard/replyKeyboard"
	"github.com/BobaUbisoft17/chsuBot/internal/schedule"
	"github.com/BobaUbisoft17/chsuBot/pkg"
	"github.com/NicoNex/echotron/v3"
)

func (b *bot) answer(answer string, keyboard echotron.ReplyMarkup) {
	messageOptions := getReplyMarkupMessageOptions(keyboard)
	_, err := b.SendMessage(answer, b.chatID, &messageOptions)
	if err != nil {
		if err.Error() == "API error: 403 Forbidden: bot was blocked by the user" {
			b.usersDb.DeleteUser(int64(b.chatID))
		} else {
			b.logger.Errorf("%v", err)
		}
	}
}

func getReplyMarkupMessageOptions(replyMarkup echotron.ReplyMarkup) echotron.MessageOptions {
	return echotron.MessageOptions{
		ReplyMarkup: replyMarkup,
		ParseMode:   "Markdown",
	}
}

func (b *bot) sendTextPost() {
	var wg sync.WaitGroup
	userIDs := b.usersDb.GetUsersId()
	for _, userID := range userIDs {
		wg.Add(1)
		go func(userID int, text string) {
			defer wg.Done()
			_, err := b.SendMessage(text, int64(userID), nil)
			if err != nil {
				if err.Error() == "API error: 403 Forbidden: bot was blocked by the user" {
					b.usersDb.DeleteUser(int64(userID))
				} else {
					b.logger.Errorf("%v", err)
				}
			}
		}(userID, b.postText)
	}
	wg.Wait()
	b.state = b.HandleMessage
	b.postText = ""
	b.answer("Все пользователи оповещены", kb.GreetingKeyboard())
}

func (b *bot) sendPostWithImage(postPhoto echotron.InputFile) {
	var wg sync.WaitGroup
	userIDs := b.usersDb.GetUsersId()
	photoOpts := echotron.PhotoOptions{
		Caption: b.postText,
	}
	for _, userID := range userIDs {
		wg.Add(1)
		go func(userID int64, photo echotron.InputFile, photoOpts echotron.PhotoOptions) {
			defer wg.Done()
			_, err := b.SendPhoto(photo, int64(userID), &photoOpts)
			if err.Error() == "API error: 403 Forbidden: bot was blocked by the user" {
				b.usersDb.DeleteUser(int64(userID))
			} else {
				b.logger.Errorf("%v", err)
			}
		}(int64(userID), postPhoto, photoOpts)
	}
	wg.Wait()
	b.postText = ""
	b.answer("Все пользователи оповещены", kb.GreetingKeyboard())
}

func (b *bot) sendSchedule() {
	b.state = b.HandleMessage
	if b.endDate == "" {
		b.endDate = b.startDate
	}

	unParseSchedule, err := b.chsuAPI.One(b.startDate, b.endDate, b.group)
	if err != nil {
		b.logger.Errorf("%v", err)
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
	keys := pkg.GetKeys(sortedSchedule)
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
		b.logger.Errorf("Error getting date: %v", err)
	}
	var markup echotron.InlineKeyboardMarkup
	if strings.Contains(callback.Data, "next") {
		markup = ikb.New(month, year).NextMonth()
	} else {
		markup = ikb.New(month, year).PreviousMonth()
	}
	message := echotron.NewMessageID(b.chatID, callback.Message.ID)
	opts := echotron.MessageReplyMarkup{ReplyMarkup: markup}
	b.EditMessageReplyMarkup(message, &opts)
}

func (b *bot) closeCalendarMarkup(callback *echotron.CallbackQuery) {
	message := echotron.NewMessageID(b.chatID, callback.Message.ID)
	b.EditMessageText("Вложение удалено", message, nil)
}

func (b *bot) getFirstSymbolKeyboard() {
	replyMarkup := ikb.FirstSymbolKeyboard()
	messageOptions := getReplyMarkupMessageOptions(replyMarkup)
	b.SendMessage("Выберите первую цифру номера вашей группы", b.chatID, &messageOptions)
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

//func (b *bot) editGroupKeyboard(message string) {
//	if message == "Назад" {
//		b.state = b.HandleMessage
//		b.answer(
//			"Возвращаемся в главное меню",
//			kb.GreetingKeyboard(),
//		)
//	} else if message == "Дальше »" || message == "« Обратно" {
//		b.answer(
//			"Меняем клавиатуру",
//			kb.GetKeyboardPart(message, b.groupDb.GetGroupNames()),
//		)
//	}
//}

func sortScheduleByDate(timetable []schedule.Lecture) (map[int][]schedule.Lecture, error) {
	scheduleByDays := make(map[int][]schedule.Lecture)
	for _, lecture := range timetable {
		timestamp, err := pkg.StringToTimestamp(lecture.DateEvent)
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
	return time.Parse("02.01.2006", date)
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

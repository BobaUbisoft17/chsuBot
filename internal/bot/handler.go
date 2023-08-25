package bot

import (
	"fmt"
	"strings"
	"time"

	calendar "github.com/BobaUbisoft17/chsuBot/internal/bot/keyboard/inlineKeyboard"
	kb "github.com/BobaUbisoft17/chsuBot/internal/bot/keyboard/replyKeyboard"
	"github.com/NicoNex/echotron/v3"
)

var inlineKeyboardCallbacks = map[string]bool{
	"next": true,
	"back": true,
	"menu": true,
}

var manageGroupKeyboard = map[string]bool{
	"Назад":     true,
	"Дальше »":  true,
	"« Обратно": true,
}

func (b *bot) HandleMessage(update *echotron.Update) stateFn {
	switch update.Message.Text {
	case "/start":
		b.sendWelcome()
	case "Сделать запись":
		b.IsAdmin(b.createPost)
	case "Узнать расписание":
		b.chooseDate()
	case "Назад":
		b.back()
	case "На сегодня":
		b.getTodaySchedule()
	case "На завтра":
		b.getTomorrowSchedule()
	case "Выбрать другой день":
		b.getAnotherDateSchedule()
	case "Выбрать диапазон":
		b.getDurationSchedule()
	case "Настройки":
		b.getSettings()
	case "Запомнить группу":
		b.memoryGroup()
	case "Изменить группу":
		b.changeGroup()
	case "Удалить данные о группе":
		b.deleteGroupInfo()
	}
	return b.state
}

func (b *bot) sendWelcome() {
	replyMarkup := kb.GreetingKeyboard()
	res, _ := b.GetChat(b.chatID)
	b.answer(
		fmt.Sprintf(
			"Здравствуйте, %s!!!\nЯ бот, упрощающий получение расписания занятий ЧГУ",
			res.Result.FirstName,
		),
		replyMarkup,
	)
}

func (b *bot) chooseDate() {
	replyMarkup := kb.ChooseDateMarkup()
	b.answer("Выберите дату", replyMarkup)
}

func (b *bot) back() {
	replyMarkup := kb.GreetingKeyboard()
	b.answer("Возвращаемся в главное меню", replyMarkup)
}

func (b *bot) getTodaySchedule() {
	if b.usersDb.IsUserHasGroup(b.chatID) || b.group != 0 {
		var groupID int
		if b.group != 0 {
			groupID, b.group = b.group, 0
		} else {
			groupID = b.usersDb.GetUserGroup(b.chatID)
		}
		b.answer(
			b.groupDb.GetTodaySchedule(groupID),
			kb.ChooseDateMarkup(),
		)
	} else {
		b.state = b.getGroup
		b.nextState = b.getTodaySchedule
		b.answer("Введите название вашей группы", kb.FirstPartGroups(b.groupDb.GetGroupNames()))
	}
}

func (b *bot) getTomorrowSchedule() {
	if b.usersDb.IsUserHasGroup(b.chatID) || b.group != 0 {
		var groupID int
		if b.group != 0 {
			groupID, b.group = b.group, 0
		} else {
			groupID = b.usersDb.GetUserGroup(b.chatID)
		}
		b.answer(
			b.groupDb.GetTomorrowSchedule(groupID),
			kb.ChooseDateMarkup(),
		)
	} else {
		b.nextState = b.getTomorrowSchedule
		b.state = b.getGroup
		b.answer("Введите название вашей группы", kb.FirstPartGroups(b.groupDb.GetGroupNames()))
	}
}

func (b *bot) getAnotherDateSchedule() {
	b.state = b.getDate
	timeNow := time.Now()
	b.answer(
		"Выберите день:",
		calendar.New(
			int(timeNow.Month()),
			timeNow.Year(),
		).BuildMarkup(),
	)
}

func (b *bot) getDate(update *echotron.Update) stateFn {
	callback := update.CallbackQuery
	switch {
	case callback != nil && inlineKeyboardCallbacks[strings.Split(callback.Data, " ")[0]]:
		b.manageCalendarKeyboard(callback)
	case callback != nil && callback.Data != "nil":
		b.startDate = callback.Data
		b.closeCalendarMarkup(callback)
		if b.usersDb.IsUserHasGroup(b.chatID) {
			b.group = b.usersDb.GetUserGroup(b.chatID)
			b.sendSchedule()
			b.state = b.HandleMessage
		} else {
			b.state = b.getGroup
			b.nextState = b.sendSchedule
			b.getGroupKeyboard()
		}
	}
	return b.state
}

func (b *bot) getGroup(update *echotron.Update) stateFn {
	message := update.Message.Text
	if manageGroupKeyboard[message] {
		b.editGroupKeyboard(message)
	} else if b.groupDb.GroupNameIsCorrect(message) {
		b.group = b.groupDb.GroupId(message)
		b.nextState()
		b.state = b.HandleMessage
	} else {
		b.answer(
			"Вы ввели некорректные данные, попробуйте ещё раз, пожалуйста",
			nil,
		)
	}
	return b.state
}

func (b *bot) getDurationSchedule() {
	b.state = b.getStartDate
	timeNow := time.Now()
	b.answer(
		"Выберите первый день диапазона:",
		calendar.New(
			int(timeNow.Month()),
			timeNow.Year(),
		).BuildMarkup(),
	)
}

func (b *bot) getStartDate(update *echotron.Update) stateFn {
	callback := update.CallbackQuery
	switch {
	case callback != nil && inlineKeyboardCallbacks[strings.Split(callback.Data, " ")[0]]:
		b.manageCalendarKeyboard(callback)
	case callback != nil && callback.Data != "nil":
		b.startDate = callback.Data
		b.state = b.getSecondDate
		b.answer("Выберите последний день диапазона (выберите день на клавиатуре сверху)", nil)
	}
	return b.state
}

func (b *bot) getSecondDate(update *echotron.Update) stateFn {
	callback := update.CallbackQuery
	switch {
	case callback != nil && inlineKeyboardCallbacks[strings.Split(callback.Data, " ")[0]]:
		b.manageCalendarKeyboard(callback)
	case callback != nil && callback.Data != "nil":
		b.endDate = callback.Data
		_ = b.orderDateCheck()
		if b.validDuration() {
			b.closeCalendarMarkup(callback)
			if b.usersDb.IsUserHasGroup(b.chatID) {
				b.group = b.usersDb.GetUserGroup(b.chatID)
				b.sendSchedule()
				b.state = b.HandleMessage
			} else {
				b.state = b.getGroup
				b.nextState = b.sendSchedule
				b.getGroupKeyboard()
			}
		} else {
			b.answer(
				"Вы ввели слишком большой диапазон. Максимальная длина диапазона не должна превышать 31 дня. (Выберите другой день на клавиатуре)",
				nil,
			)
		}
	}
	return b.state
}

func (b *bot) getSettings() {
	var replyMarkup echotron.ReplyKeyboardMarkup
	if !b.usersDb.IsUserInDB(b.chatID) {
		b.usersDb.AddUser(b.chatID)
	}
	if b.usersDb.IsUserHasGroup(b.chatID) {
		replyMarkup = kb.ChangeGroupKeyboard()
	} else {
		replyMarkup = kb.MemoryGroupKeyboard()
	}
	b.answer("Переходим в меню настроек", replyMarkup)
}

func (b *bot) memoryGroup() {
	if !b.usersDb.IsUserHasGroup(b.chatID) {
		b.state = b.addUserGroup
		b.getGroupKeyboard()
	} else {
		b.answer("Не ломайте меня, пожалуйста🙏", nil)
	}
}

func (b *bot) addUserGroup(update *echotron.Update) stateFn {
	message := update.Message.Text
	if manageGroupKeyboard[message] {
		b.editGroupKeyboard(message)
	} else if b.groupDb.GroupNameIsCorrect(message) {
		b.usersDb.ChangeUserGroup(b.chatID, message)
		b.state = b.HandleMessage
		b.answer(
			"Я вас запомнил, теперь вам не нужно выбирать группу",
			kb.GreetingKeyboard(),
		)
	} else {
		b.answer(
			"Вы ввели некорректные данные, попробуйте ещё раз, пожалуйста",
			nil,
		)
	}
	return b.state
}

func (b *bot) changeGroup() {
	if b.usersDb.IsUserHasGroup(b.chatID) {
		b.state = b.updateUserGroup
		b.getGroupKeyboard()
	} else {
		b.answer("Не ломайте меня, пожалуйста🙏", nil)
	}
}

func (b *bot) updateUserGroup(update *echotron.Update) stateFn {
	message := update.Message.Text
	if manageGroupKeyboard[message] {
		b.editGroupKeyboard(message)
	} else if b.groupDb.GroupNameIsCorrect(message) {
		if b.usersDb.GetUserGroup(b.chatID) != b.groupDb.GroupId(message) {
			b.state = b.HandleMessage
			b.usersDb.ChangeUserGroup(b.chatID, message)
			b.answer(
				"Вы успешно изменили группу",
				kb.GreetingKeyboard(),
			)
		} else {
			b.answer(
				"Эта группа уже выбрана вами",
				nil,
			)
		}
	} else {
		b.answer(
			"Вы ввели некорректные данные, попробуйте ещё раз, пожалуйста",
			nil,
		)
	}
	return b.state
}

func (b *bot) deleteGroupInfo() {
	if b.usersDb.IsUserHasGroup(b.chatID) {
		b.usersDb.DeleteGroup(b.chatID)
		replyMarkup := kb.GreetingKeyboard()
		b.answer("Данные о вашей группе успешно удалены", replyMarkup)
	} else {
		b.answer("Не ломайте меня, пожалуйста🙏", nil)
	}
}

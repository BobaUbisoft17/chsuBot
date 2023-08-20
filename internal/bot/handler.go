package bot

import (
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
	b.answer("Hello world!", replyMarkup)
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
	if b.db.IsUserHasGroup(b.chatID) || b.group != 0 {
		var groupID int
		if b.group != 0 {
			groupID, b.group = b.group, 0
		} else {
			groupID = b.db.GetUserGroup(b.chatID)
		}
		b.answer(
			b.db.GetTodaySchedule(groupID),
			kb.ChooseDateMarkup(),
		)
	} else {
		b.state = b.getGroup
		b.nextState = b.getTodaySchedule
		b.answer("Введите название вашей группы", kb.FirstPartGroups(b.db.GetGroupNames()))
	}
}

func (b *bot) getTomorrowSchedule() {
	if b.db.IsUserHasGroup(b.chatID) || b.group != 0 {
		var groupID int
		if b.group != 0 {
			groupID, b.group = b.group, 0
		} else {
			groupID = b.db.GetUserGroup(b.chatID)
		}
		b.answer(
			b.db.GetTomorrowSchedule(groupID),
			kb.ChooseDateMarkup(),
		)
	} else {
		b.nextState = b.getTomorrowSchedule
		b.state = b.getGroup
		b.answer("Введите название вашей группы", kb.FirstPartGroups(b.db.GetGroupNames()))
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
		if b.db.IsUserHasGroup(b.chatID) {
			b.group = b.db.GetUserGroup(b.chatID)
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
	} else if b.db.GroupNameIsCorrect(message) {
		b.group = b.db.GroupId(message)
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
			if b.db.IsUserHasGroup(b.chatID) {
				b.group = b.db.GetUserGroup(b.chatID)
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
	if !b.db.IsUserInDB(b.chatID) {
		b.db.AddUser(b.chatID)
	}
	if b.db.IsUserHasGroup(b.chatID) {
		replyMarkup = kb.ChangeGroupKeyboard()
	} else {
		replyMarkup = kb.MemoryGroupKeyboard()
	}
	b.answer("Переходим в меню настроек", replyMarkup)
}

func (b *bot) memoryGroup() {
	if !b.db.IsUserHasGroup(b.chatID) {
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
	} else if b.db.GroupNameIsCorrect(message) {
		b.db.ChangeUserGroup(b.chatID, message)
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
	if b.db.IsUserHasGroup(b.chatID) {
		b.state = b.updateUserGroup
		b.getGroupKeyboard()
	} else {
		b.answer("Не ломайте меня, пожалуйста🙏", nil)
	}
}

func (b *bot) updateUserGroup(update *echotron.Update) stateFn {
	groupName := update.Message.Text
	if b.db.GroupNameIsCorrect(groupName) {
		if b.db.GetUserGroup(b.chatID) != b.db.GroupId(groupName) {
			b.state = b.HandleMessage
			b.db.ChangeUserGroup(b.chatID, groupName)
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
	if b.db.IsUserHasGroup(b.chatID) {
		b.db.DeleteGroup(b.chatID)
		replyMarkup := kb.GreetingKeyboard()
		b.answer("Данные о вашей группе успешно удалены", replyMarkup)
	} else {
		b.answer("Не ломайте меня, пожалуйста🙏", nil)
	}
}

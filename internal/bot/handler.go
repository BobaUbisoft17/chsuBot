package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"slices"

	ikb "github.com/BobaUbisoft17/chsuBot/internal/bot/keyboard/inlineKeyboard"
	kb "github.com/BobaUbisoft17/chsuBot/internal/bot/keyboard/replyKeyboard"
	"github.com/NicoNex/echotron/v3"
)

var calendarKeyboardCallbacks = []string{"next", "back", "menu"}

var groupKeyboardCallbacks = []string{"next", "previous", "back"}

func (b *bot) HandleMessage(update *echotron.Update) stateFn {
	if update.Message != nil {
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
			b.rememberGroup()
		case "Изменить группу":
			b.changeGroup()
		case "Удалить данные о группе":
			b.deleteGroupInfo()
		case "Помощь":
			b.help()
		}
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
	hasGroup, err := b.usePackages.usersDb.IsUserHasGroup(b.chatID)
	if err != nil {
		b.botError(err)
		return
	}
	if hasGroup || b.group != 0 {
		var groupID int
		if b.group != 0 {
			groupID, b.group = b.group, 0
		} else {
			groupID, err = b.usePackages.usersDb.GetUserGroup(b.chatID)
			if err != nil {
				b.botError(err)
				return
			}
		}
		b.state = b.HandleMessage
		todaySchedule, err := b.usePackages.groupDb.GetTodaySchedule(groupID)
		if err != nil {
			b.botError(err)
			return
		}
		b.answer(
			todaySchedule,
			kb.ChooseDateMarkup(),
		)
	} else {
		b.state = b.chooseUniversity
		b.nextFn = b.getTodaySchedule
		b.previousFn = b.chooseDate
		b.answer("Введите название вашей группы", ikb.FirstSymbolKeyboard())
	}
}

func (b *bot) getTomorrowSchedule() {
	hasGroup, err := b.usePackages.usersDb.IsUserHasGroup(b.chatID)
	if err != nil {
		b.botError(err)
		return
	}
	if hasGroup || b.group != 0 {
		var groupID int
		if b.group != 0 {
			groupID, b.group = b.group, 0
		} else {
			groupID, err = b.usePackages.usersDb.GetUserGroup(b.chatID)
			if err != nil {
				b.botError(err)
				return
			}
		}
		b.state = b.HandleMessage
		tomorrowScedule, err := b.usePackages.groupDb.GetTomorrowSchedule(groupID)
		if err != nil {
			b.botError(err)
			return
		}
		b.answer(
			tomorrowScedule,
			kb.ChooseDateMarkup(),
		)
	} else {
		b.nextFn = b.getTomorrowSchedule
		b.previousFn = b.chooseDate
		b.state = b.chooseUniversity
		b.answer("Введите название вашей группы", ikb.FirstSymbolKeyboard())
	}
}

func (b *bot) getAnotherDateSchedule() {
	b.state = b.getDate
	b.previousFn = b.chooseDate
	timeNow := time.Now()
	b.endDate = time.Time{}
	b.answer(
		"Выберите день:",
		ikb.New(
			int(timeNow.Month()),
			timeNow.Year(),
		).BuildMarkup(),
	)
}

func (b *bot) getDate(update *echotron.Update) stateFn {
	callback := update.CallbackQuery
	if callback != nil {
		command := strings.Split(callback.Data, " ")[0]
		switch {
		case slices.Contains(calendarKeyboardCallbacks, command):
			b.manageCalendarKeyboard(callback)
		case callback.Data != "nil":
			b.startDate, _ = parseDate(callback.Data)
			hasGroup, err := b.usePackages.usersDb.IsUserHasGroup(b.chatID)
			if err != nil {
				b.botError(err)
				return b.state
			}
			if hasGroup {
				b.closeCalendarMarkup(callback)
				b.group, err = b.usePackages.usersDb.GetUserGroup(b.chatID)
				if err != nil {
					b.botError(err)
					return b.state
				}
				b.sendSchedule()
				b.state = b.HandleMessage
			} else {
				b.state = b.chooseUniversity
				b.nextFn = b.sendSchedule
				b.editMessage(
					callback.Message.ID,
					"Выберите первую цифру вашей группы",
					ikb.FirstSymbolKeyboard(),
				)
			}
		}
	}
	return b.state
}

func (b *bot) chooseUniversity(update *echotron.Update) stateFn {
	callback := update.CallbackQuery
	if callback != nil {
		switch {
		case callback.Data == "back":
			b.editMessage(
				callback.Message.ID,
				"Вложение удалено",
				echotron.InlineKeyboardMarkup{},
			)
			b.previousFn()
			b.state = b.HandleMessage
		default:
			b.state = b.getGroup
			groupsInfo, err := b.usePackages.groupDb.GroupsStartsWith(callback.Data)
			if err != nil {
				b.botError(err)
				return b.state
			}
			keyboard := ikb.CreateGroupKeyboard(
				groupsInfo,
				callback.Data,
				1,
			)
			b.editMessage(callback.Message.ID, "Выберите вашу группу", keyboard)
		}
	}
	return b.state
}

func (b *bot) getGroup(update *echotron.Update) stateFn {
	callback := update.CallbackQuery
	if callback != nil {
		command := strings.Split(callback.Data, " ")[0]
		switch {
		case command == "back":
			b.editMessage(
				callback.Message.ID,
				"Выберите первую цифру вашей группы",
				ikb.FirstSymbolKeyboard(),
			)
			b.state = b.chooseUniversity
		case slices.Contains(groupKeyboardCallbacks, command):
			splitData := strings.Split(callback.Data, " ")
			university, stringPart := splitData[1], splitData[2]
			part, _ := strconv.Atoi(stringPart)
			groups, err := b.usePackages.groupDb.GroupsStartsWith(university)
			if err != nil {
				b.botError(err)
				return b.state
			}
			b.editKeyboard(
				callback.Message.ID,
				ikb.CreateGroupKeyboard(groups, university, part),
			)
		default:
			b.editMessage(
				callback.Message.ID,
				"Вложение удалено",
				echotron.InlineKeyboardMarkup{},
			)
			b.group, _ = strconv.Atoi(callback.Data)
			b.nextFn()
		}
	}
	return b.state
}

func (b *bot) getDurationSchedule() {
	b.state = b.getStartDate
	b.previousFn = b.chooseDate
	timeNow := time.Now()
	b.answer(
		"Выберите первый день диапазона:",
		ikb.New(
			int(timeNow.Month()),
			timeNow.Year(),
		).BuildMarkup(),
	)
}

func (b *bot) getStartDate(update *echotron.Update) stateFn {
	callback := update.CallbackQuery
	if callback != nil {
		command := strings.Split(callback.Data, " ")[0]
		switch {
		case slices.Contains(calendarKeyboardCallbacks, command):
			b.manageCalendarKeyboard(callback)
		case callback.Data != "nil":
			b.startDate, _ = parseDate(callback.Data)
			b.state = b.getSecondDate
			b.answer(
				"Выберите последний день диапазона (выберите день на клавиатуре сверху)",
				nil,
			)
		}
	}
	return b.state
}

func (b *bot) getSecondDate(update *echotron.Update) stateFn {
	callback := update.CallbackQuery
	if callback != nil {
		switch {
		case slices.Contains(calendarKeyboardCallbacks, strings.Split(callback.Data, " ")[0]):
			b.manageCalendarKeyboard(callback)
		case callback.Data != "nil":
			b.endDate, _ = parseDate(callback.Data)
			b.dateSequenceCorrection()
			if b.validDuration() {
				hasGroup, err := b.usePackages.usersDb.IsUserHasGroup(b.chatID)
				if err != nil {
					b.botError(err)
					return b.state
				}
				if hasGroup {
					b.closeCalendarMarkup(callback)
					b.group, err = b.usePackages.usersDb.GetUserGroup(b.chatID)
					if err != nil {
						b.botError(err)
						return b.state
					}
					b.sendSchedule()
					b.state = b.HandleMessage
				} else {
					b.state = b.chooseUniversity
					b.nextFn = b.sendSchedule
					b.editMessage(
						callback.Message.ID,
						"Выберите первую цифру вашей группы",
						ikb.FirstSymbolKeyboard(),
					)
				}
			} else {
				b.answer(
					"Вы ввели слишком большой диапазон. "+
						"Максимальная длина диапазона не должна превышать 31 дня. "+
						"(Выберите другой день на клавиатуре)",
					nil,
				)
			}
		}
	}
	return b.state
}

func (b *bot) getSettings() {
	var replyMarkup echotron.ReplyKeyboardMarkup
	inDB, err := b.usePackages.usersDb.IsUserInDB(b.chatID)
	if err != nil {
		b.botError(err)
		return
	}
	if !inDB {
		b.usePackages.usersDb.AddUser(b.chatID)
	}
	hasGroup, err := b.usePackages.usersDb.IsUserHasGroup(b.chatID)
	if err != nil {
		b.botError(err)
		return
	}
	if hasGroup {
		replyMarkup = kb.ChangeGroupKeyboard()
	} else {
		replyMarkup = kb.MemoryGroupKeyboard()
	}
	b.answer("Переходим в меню настроек", replyMarkup)
}

func (b *bot) rememberGroup() {
	hasGroup, err := b.usePackages.usersDb.IsUserHasGroup(b.chatID)
	if err != nil {
		b.botError(err)
		return
	}
	if !hasGroup {
		b.state = b.chooseUniversity
		b.nextFn = b.addUserGroup
		b.previousFn = b.getSettings
		b.getFirstSymbolKeyboard()
	} else {
		b.answer("Не ломайте меня, пожалуйста🙏", nil)
	}
}

func (b *bot) addUserGroup() {
	b.usePackages.usersDb.ChangeUserGroup(b.chatID, b.group)
	b.state = b.HandleMessage
	b.group = 0
	b.answer(
		"Я вас запомнил, теперь вам не нужно выбирать группу",
		kb.GreetingKeyboard(),
	)
}

func (b *bot) changeGroup() {
	hasGroup, err := b.usePackages.usersDb.IsUserHasGroup(b.chatID)
	if err != nil {
		b.botError(err)
		return
	}
	if hasGroup {
		b.state = b.chooseUniversity
		b.nextFn = b.updateUserGroup
		b.previousFn = b.getSettings
		b.getFirstSymbolKeyboard()
	} else {
		b.answer("Не ломайте меня, пожалуйста🙏", nil)
	}
}

func (b *bot) updateUserGroup() {
	b.state = b.HandleMessage
	userGroup, err := b.usePackages.usersDb.GetUserGroup(b.chatID)
	if err != nil {
		b.botError(err)
		return
	}
	if userGroup != b.group {
		b.usePackages.usersDb.ChangeUserGroup(b.chatID, b.group)
		b.group = 0
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
}

func (b *bot) deleteGroupInfo() {
	hasGroup, err := b.usePackages.usersDb.IsUserHasGroup(b.chatID)
	if err != nil {
		b.botError(err)
		return
	}
	if hasGroup {
		b.usePackages.usersDb.DeleteGroup(b.chatID)
		replyMarkup := kb.GreetingKeyboard()
		b.answer("Данные о вашей группе успешно удалены", replyMarkup)
	} else {
		b.answer("Не ломайте меня, пожалуйста🙏", nil)
	}
}

func (b *bot) help() {
	b.answer(
		"Бот, упрощающий получение расписания студениами ЧГУ.\n\n"+
			"Получение расписания - можно получать расписание как на сегодня/завтра, "+
			"так и на произвольную дату или произвольный промежуток."+
			"Есть функция запоминания группы пользователя для получения "+
			"расписания по нажатию *одной кнопки.\n\n"+
			"Исходный код выложен на GitHub "+
			"https://github.com/BobaUbisoft17/chsuBot\n"+
			"Связаться с автором проекта:\n"+
			"Телеграм @BobaUbisoft\n"+
			"VK vk.com/bobaubisoft\n"+
			"Почта aksud2316@gmail.com\n\n"+
			"Поддержать проект: 5536 9137 8142 8269",
		nil,
	)
}

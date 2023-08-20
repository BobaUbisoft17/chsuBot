package bot

import (
	"sort"

	"github.com/NicoNex/echotron/v3"
)

func GreetingKeyboard() echotron.ReplyKeyboardMarkup {
	keyboard := [][]echotron.KeyboardButton{
		add("Узнать расписание"),
		add("Настройки"),
		add("Помощь"),
	}
	return echotron.ReplyKeyboardMarkup{
		Keyboard:       keyboard,
		ResizeKeyboard: true,
	}
}

func ChooseDateMarkup() echotron.ReplyKeyboardMarkup {
	keyboard := [][]echotron.KeyboardButton{
		row("На сегодня", "На завтра"),
		row("Выбрать другой день", "Выбрать диапазон"),
		add("Назад"),
	}
	return echotron.ReplyKeyboardMarkup{
		Keyboard:       keyboard,
		ResizeKeyboard: true,
	}
}

func MemoryGroupKeyboard() echotron.ReplyKeyboardMarkup {
	keyboard := [][]echotron.KeyboardButton{
		add("Запомнить группу"),
		add("Назад"),
	}
	return echotron.ReplyKeyboardMarkup{
		Keyboard:       keyboard,
		ResizeKeyboard: true,
	}
}

func ChangeGroupKeyboard() echotron.ReplyKeyboardMarkup {
	keyboard := [][]echotron.KeyboardButton{
		add("Изменить группу"),
		add("Удалить данные о группе"),
		add("Назад"),
	}
	return echotron.ReplyKeyboardMarkup{
		Keyboard:       keyboard,
		ResizeKeyboard: true,
	}
}

func FirstPartGroups(groupNames []string) echotron.ReplyKeyboardMarkup {
	sort.Strings(groupNames)
	keyboard := [][]echotron.KeyboardButton{
		add("Дальше »"),
		add("Назад"),
	}
	for i := 0; i < (len(groupNames)+6)/2; i++ {
		keyboard = append(keyboard, add(groupNames[i]))
	}
	keyboard = append(keyboard, add("Дальше »"))
	return echotron.ReplyKeyboardMarkup{
		Keyboard:       keyboard,
		ResizeKeyboard: true,
	}
}

func SecondPartGroups(groupNames []string) echotron.ReplyKeyboardMarkup {
	sort.Strings(groupNames)
	keyboard := [][]echotron.KeyboardButton{
		add("« Обратно"),
		add("Назад"),
	}
	for i := (len(groupNames) + 6) / 2; i < len(groupNames); i++ {
		keyboard = append(keyboard, add(groupNames[i]))
	}
	keyboard = append(keyboard, add("« Обратно"))
	return echotron.ReplyKeyboardMarkup{
		Keyboard:       keyboard,
		ResizeKeyboard: true,
	}
}

func GetKeyboardPart(message string, groupNames []string) echotron.ReplyKeyboardMarkup {
	if message == "Дальше »" {
		return SecondPartGroups(groupNames)
	}
	return FirstPartGroups(groupNames)
}

func row(buttons ...string) []echotron.KeyboardButton {
	row := make([]echotron.KeyboardButton, 0, len(buttons))
	for _, button := range buttons {
		row = append(row, echotron.KeyboardButton{Text: button})
	}
	return row
}

func add(text string) []echotron.KeyboardButton {
	return []echotron.KeyboardButton{
		{Text: text},
	}
}

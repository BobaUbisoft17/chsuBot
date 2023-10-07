package bot

import (
	"fmt"
	"strconv"

	"github.com/BobaUbisoft17/chsuBot/internal/database"
	"github.com/NicoNex/echotron/v3"
)

func FirstSymbolKeyboard() echotron.InlineKeyboardMarkup {
	keyboard := [][]echotron.InlineKeyboardButton{}
	keyboard = append(
		keyboard,
		[]echotron.InlineKeyboardButton{
			addButton("0", "0")[0],
			addButton("1", "1")[0],
			addButton("2", "2")[0],
		},
		[]echotron.InlineKeyboardButton{
			addButton("3", "3")[0],
			addButton("4", "4")[0],
			addButton("5", "5")[0],
		},
		[]echotron.InlineKeyboardButton{
			addButton("6", "6")[0],
			addButton("7", "7")[0],
			addButton("9", "9")[0],
		},
	)

	keyboard = append(keyboard, addButton("Назад", "back"))
	return echotron.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}

func getGroupButtons(groups []database.GroupInfo) [][]echotron.InlineKeyboardButton {
	keyboard := [][]echotron.InlineKeyboardButton{}
	for i := 0; i < len(groups)-1; i += 2 {
		keyboard = append(keyboard, []echotron.InlineKeyboardButton{
			addButton(groups[i].GroupName, strconv.Itoa(groups[i].GroupID))[0],
			addButton(groups[i+1].GroupName, strconv.Itoa(groups[i+1].GroupID))[0],
		})
	}
	if len(groups)%2 != 0 {
		keyboard = append(keyboard, addButton(
			groups[len(groups)-1].GroupName,
			strconv.Itoa(groups[len(groups)-1].GroupID),
		))
	}
	return keyboard
}

func getNavButtons(university string, amountGroups, part int) []echotron.InlineKeyboardButton {
	back := addButton("Назад", "back")[0]
	previous := addButton("<", fmt.Sprintf("previous %s %d", university, part-1))[0]
	next := addButton(">", fmt.Sprintf("next %s %d", university, part+1))[0]
	if amountGroups < 18 {
		return []echotron.InlineKeyboardButton{back}
	}
	if part == 1 {
		return []echotron.InlineKeyboardButton{back, next}
	}
	if part*18 >= amountGroups {
		return []echotron.InlineKeyboardButton{previous, back}
	}
	return []echotron.InlineKeyboardButton{previous, back, next}
}

func CreateGroupKeyboard(groups []database.GroupInfo, university string, part int) echotron.InlineKeyboardMarkup {
	var end int
	start := 18 * (part - 1)
	if 18*part <= len(groups) {
		end = 18 * part
	} else {
		end = len(groups)
	}
	keyboard := getGroupButtons(groups[start:end])
	keyboard = append(keyboard, getNavButtons(university, len(groups), part))
	return echotron.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}

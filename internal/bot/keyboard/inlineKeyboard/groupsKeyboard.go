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

func CreateGroupKeyboard(groups []database.GroupInfo, university string, part int) echotron.InlineKeyboardMarkup {
	var keyboard [][]echotron.InlineKeyboardButton
	if len(groups) < 18 {
		keyboard = getGroupButtons(groups)
		keyboard = append(keyboard, addButton("Назад", "back"))
	} else if part == 1 {
		keyboard = getGroupButtons(groups[:18])
		keyboard = append(keyboard, []echotron.InlineKeyboardButton{
			addButton("Назад", "back")[0],
			addButton(">", fmt.Sprintf("next %s 2", university))[0],
		})
	} else if part*18 > len(groups) {
		keyboard = getGroupButtons(groups[(part-1)*18:])
		keyboard = append(keyboard, []echotron.InlineKeyboardButton{
			addButton("<", fmt.Sprintf("previous %s %d", university, part-1))[0],
			addButton("Назад", "back")[0],
		})
	} else {
		keyboard = getGroupButtons(groups[(part-1)*18 : part*18])
		keyboard = append(keyboard, []echotron.InlineKeyboardButton{
			addButton("<", fmt.Sprintf("previous %s %d", university, part-1))[0],
			addButton("Назад", "back")[0],
			addButton(">", fmt.Sprintf("next %s %d", university, part+1))[0],
		})
	}
	return echotron.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}

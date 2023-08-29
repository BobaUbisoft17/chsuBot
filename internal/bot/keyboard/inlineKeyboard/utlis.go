package bot

import "github.com/NicoNex/echotron/v3"

func rowButtons(buttons []string) (buttonLine []echotron.InlineKeyboardButton) {
	for _, button := range buttons {
		buttonLine = append(buttonLine, addButton(button, "nil")...)
	}
	return
}

func addButton(text string, callback string) (inlineButton []echotron.InlineKeyboardButton) {
	return []echotron.InlineKeyboardButton{
		{Text: text, CallbackData: callback},
	}
}

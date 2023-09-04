package bot

import (
	kb "github.com/BobaUbisoft17/chsuBot/internal/bot/keyboard/replyKeyboard"
	"github.com/NicoNex/echotron/v3"
)

func (b *bot) IsAdmin(nextFunc nextFn) {
	if int(b.chatID) == b.usePackages.adminId {
		nextFunc()
	}
}

func (b *bot) createPost() {
	b.state = b.choosePostType
	b.answer(
		"Выберите тип поста",
		kb.PostKeyboard(),
	)
}

func (b *bot) choosePostType(update *echotron.Update) stateFn {
	switch update.Message.Text {
	case "Текстовый пост":
		b.prepareGetPostText()
		b.nextFn = b.sendTextPost
	case "Фото":
		b.prepareGetPostPhoto()
	case "Смешанный пост":
		b.prepareGetPostText()
		b.nextFn = b.prepareGetPostPhoto
	case "Назад":
		b.state = b.HandleMessage
		b.answer(
			"Возвращаемся в главное меню",
			kb.GreetingKeyboard(),
		)
	}
	return b.state
}

func (b *bot) prepareGetPostText() {
	b.state = b.getPostText
	b.answer(
		"Напишите мне текст для поста",
		kb.BackButton(),
	)
}

func (b *bot) getPostText(update *echotron.Update) stateFn {
	message := update.Message.Text
	if message != "Назад" {
		b.postText = message
		b.nextFn()
	} else {
		b.createPost()
	}
	return b.state
}

func (b *bot) prepareGetPostPhoto() {
	b.state = b.getPostPhoto
	b.answer(
		"Отправьте мне фото для поста",
		kb.BackButton(),
	)
}

func (b *bot) getPostPhoto(update *echotron.Update) stateFn {
	if update.Message.Text != "Назад" {
		photo := update.Message.Photo[0].FileID
		postPhoto := echotron.NewInputFileID(photo)
		b.sendPostWithImage(postPhoto)
	} else {
		b.createPost()
	}
	return b.HandleMessage
}

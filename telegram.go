package main

import (
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func sendMessage(chatId int64, msg string) tgbotapi.Message {
	return sendMessageAsReply(chatId, msg, 0)
}

func sendMessageAsReply(chatId int64, msg string, replyToId int) tgbotapi.Message {
	return sendMessageWithKeyboard(chatId, msg, nil, replyToId)
}

func sendMessageWithKeyboard(chatId int64, msg string, keyboard *tgbotapi.InlineKeyboardMarkup, replyToId int) tgbotapi.Message {
	chattable := tgbotapi.NewMessage(chatId, msg)
	chattable.BaseChat.ReplyToMessageID = replyToId
	chattable.ParseMode = "HTML"
	chattable.DisableWebPagePreview = true
	if keyboard != nil {
		chattable.BaseChat.ReplyMarkup = *keyboard
	}
	message, err := bot.Send(chattable)
	if err != nil {
		if strings.Index(err.Error(), "reply message not found") != -1 {
			chattable.BaseChat.ReplyToMessageID = 0
			message, err = bot.Send(chattable)
		}

		log.Warn().Err(err).Int64("chat", chatId).Str("msg", msg).Msg("error sending message")
	}
	return message
}

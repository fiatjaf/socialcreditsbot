package main

import (
	"database/sql"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func getLastTelegramUpdate() (ltu int64, err error) {
	err = pg.Get(&ltu, `
      SELECT telegram_update
      FROM events
      ORDER BY time
      DESC LIMIT 1
    `)
	if err == sql.ErrNoRows {
		err = nil
		ltu = 0
	}

	return
}

func userName(user *tgbotapi.User) string {
	userName := user.UserName
	if userName == "" {
		userName = user.FirstName
	} else {
		userName = "@" + userName
	}
	return userName
}

package main

import (
	"encoding/hex"
	"fmt"

	"github.com/hoisie/mustache"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func handle(upd tgbotapi.Update) {
	if upd.Message != nil {
		if upd.Message.Sticker != nil && upd.Message.ReplyToMessage != nil {
			// get params
			points := 0
			switch hex.EncodeToString([]byte(upd.Message.Sticker.Emoji)) {
			case "f09f989e": // -20
				points = -20
			case "f09f9884": // +20
				points = +20
			default:
				return
			}

			// check user is admin
			chatmember, err := bot.GetChatMember(tgbotapi.ChatConfigWithUser{
				ChatID:             upd.Message.Chat.ID,
				SuperGroupUsername: upd.Message.Chat.ChatConfig().SuperGroupUsername,
				UserID:             upd.Message.From.ID,
			})
			if err != nil ||
				(chatmember.Status != "administrator" && chatmember.Status != "creator") {
				log.Print("not admin")
				return
			}

			// save scores
			_, err = pg.Exec(`
              INSERT INTO events
                (chat_id, user_id, creator_id, credits, telegram_update)
              VALUES ($1, $2, $3, $4, $5)
            `,
				upd.Message.Chat.ID,
				upd.Message.ReplyToMessage.From.ID,
				upd.Message.From.ID,
				points,
				upd.UpdateID,
			)
			if err != nil {
				log.Warn().Err(err).Msg("failed to save event")
				sendMessage(upd.Message.Chat.ID, "Error!")
				return
			}

			sendMessage(upd.Message.Chat.ID,
				fmt.Sprintf(
					"%s credits saved.",
					userName(upd.Message.ReplyToMessage.From),
				),
			)
		} else if upd.Message.Text == "/credits" {
			var scores []Score
			err = pg.Select(&scores, `
              SELECT * FROM (
                SELECT user_id, sum(credits) AS credits
                FROM events
                WHERE chat_id = $1
                GROUP BY user_id
              )x ORDER BY credits DESC
            `, upd.Message.Chat.ID)
			if err != nil {
				log.Warn().Err(err).Msg("error fetching scores")
				sendMessage(upd.Message.Chat.ID, "Error!")
				return
			}

			// fetch user names
			for i, score := range scores {
				// fallback
				scores[i].Name = fmt.Sprintf("user:%d", score.UserId)

				chatmember, err := bot.GetChatMember(tgbotapi.ChatConfigWithUser{
					ChatID:             upd.Message.Chat.ID,
					SuperGroupUsername: upd.Message.Chat.ChatConfig().SuperGroupUsername,
					UserID:             int(score.UserId),
				})
				if err != nil {
					continue
				}

				scores[i].Name = userName(chatmember.User)
			}

			sendMessage(upd.Message.Chat.ID,
				mustache.Render(`<b>Credits</b>
{{#scores}}<code>{{Credits}}</code>: {{Name}}
{{/scores}}
                `, map[string]interface{}{
					"scores": scores,
				}),
			)
		}
	}
}

type Score struct {
	UserId  int64 `db:"user_id"`
	Credits int64 `db:"credits"`
	Name    string
}

package telegram

import (
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	tele "gopkg.in/telebot.v4"
	"log"
	"os"
	"storybot/database"
	"strconv"
)

func RegisterAllHandlers(bot *tele.Bot, pool *pgxpool.Pool) {
	//bot.Handle("/start", func(c tele.Context) error {
	//	return handle(c, pool, "/start")
	//})

	bot.Handle(tele.OnText, func(c tele.Context) error {
		return BasicHandle(c, pool, c.Text(), bot)
	})
}

func BasicHandle(context tele.Context, pool *pgxpool.Pool, message string, bot *tele.Bot) error {

	result, err := database.SelectHandleFromDatabase(pool, message)
	if result == "" {
		adminID, _ := strconv.Atoi(os.Getenv("ADMIN_ID"))
		SendMessage(bot, int64(adminID), fmt.Sprint("Произошла ошибка: ", err))
		return context.Send("Ты написал неправильное сообщение")
	}
	return context.Send(result)
}

func SendMessage(bot *tele.Bot, chatID int64, message string) {
	_, err := bot.Send(&tele.Chat{ID: chatID}, message)
	if err != nil {
		log.Printf("Не удалось отправить сообщение: %v", err)
	}
}

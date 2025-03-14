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

var userState = map[int64]string{}
var msg string // Сообщение глобальное

func RegisterAllHandlers(bot *tele.Bot, pool *pgxpool.Pool) {
	//bot.Handle("/start", func(c tele.Context) error {
	//	return handle(c, pool, "/start")
	//})
	bot.Handle("/add", func(context tele.Context) error {
		return actionWithHandle(context, "/add")
	})
	bot.Handle("/delete", func(context tele.Context) error {
		return actionWithHandle(context, "/delete")
	})
	bot.Handle("/cancel", func(context tele.Context) error {
		return cancelHandle(context)
	})
	bot.Handle(tele.OnText, func(context tele.Context) error {
		return BasicHandle(context, pool, context.Text(), bot)

	})
}

func BasicHandle(context tele.Context, pool *pgxpool.Pool, message string, bot *tele.Bot) error {

	chatID := context.Chat().ID
	if state, exists := userState[chatID]; exists {
		switch state {
		case "remove_handle":
			err := database.RemoveHandleFromDatabase(pool, message)
			delete(userState, chatID)
			if err != nil {
				return context.Send("Произошла ошибка при удалении. Возможно такого handle не существует.")
			}
			return context.Send("Вы успешно удалили handle: " + message)
		case "wait_handle":
			userState[chatID] = "wait_answer"
			msg = message
			return context.Send("Введите answer: ")
		case "wait_answer":
			handle := msg
			answer := message
			msg = "" // На всякий обнулим msg
			err := database.AddHandleFromDatabase(pool, handle, answer)
			delete(userState, chatID)
			if err != nil {
				return context.Send("Произошла ошибка при добавлении handle")
			}
			return context.Send("Вы успешно добавили новый handle")
		}

	}

	result, err := database.SelectHandleFromDatabase(pool, message)
	if result == "" {
		adminID, _ := strconv.Atoi(os.Getenv("ADMIN_ID"))
		SendMessage(bot, int64(adminID), fmt.Sprint("Произошла ошибка: ", err))
		return context.Send("Ты написал неправильное сообщение")
	}
	return context.Send(result)
}

func actionWithHandle(context tele.Context, message string) error {
	chatID := context.Chat().ID
	adminID, _ := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	if chatID != adminID {
		return context.Send("Вы не обладаете правами для выполнения данной команды")
	}

	if message == "/add" {
		userState[chatID] = "wait_handle"
	} else {
		userState[chatID] = "remove_handle"
	}

	return context.Send("Введите handle: ")
}

func SendMessage(bot *tele.Bot, chatID int64, message string) {
	_, err := bot.Send(&tele.Chat{ID: chatID}, message)
	if err != nil {
		log.Printf("Не удалось отправить сообщение: %v", err)
	}
}

func cancelHandle(context tele.Context) error {
	chatID := context.Chat().ID
	if _, exists := userState[chatID]; exists {
		delete(userState, chatID)
		return context.Send("Операция отменена")
	}

	return context.Send("Нет активных операций для отмены.")
}

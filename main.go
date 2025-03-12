package main

import (
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
	"log"
	"os"
	"storybot/database"
	"storybot/telegram"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN не найден в файле .env")
	}

	pool := database.CreateDatabasePool()

	settings := tele.Settings{
		Token: token,
	}

	bot, err := tele.NewBot(settings)
	if err != nil {
		log.Fatalf("Бот не поднялся, что то не работает: %v", err)
	}

	telegram.RegisterAllHandlers(bot, pool)
	bot.Start()

}

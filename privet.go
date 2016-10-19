package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func process(update tgbotapi.Update) (answer tgbotapi.MessageConfig, err error) {
	if update.Message == nil {
		err = errors.New("empty")
		return
	}

	if update.Message.LeftChatMember != nil {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Ну ты чего? Нормально же общались"), nil
	}

	if update.Message.NewChatMember != nil {
		var name string
		if update.Message.NewChatMember.FirstName != "" {
			name = update.Message.NewChatMember.FirstName
		} else {
			name = update.Message.NewChatMember.UserName
		}
		message := fmt.Sprintf("Привет, %s, расскажи пару слов о себе и о проекте, о котором нам почти не рассказывают!", name)
		return tgbotapi.NewMessage(update.Message.Chat.ID, message), nil
	}

	msgText := strings.ToLower(strings.TrimSpace(update.Message.Text))

	if strings.Contains(msgText, "к ноге") && update.Message.From.UserName == "just_alyosha" {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "🐕"), nil
	} else if strings.Contains(msgText, "привет") {
		var msg tgbotapi.MessageConfig
		if update.Message.From.UserName == "ainomc" {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Миша, отстань")
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "привет")
		}
		return msg, nil
	}
	err = errors.New("empty")
	return
}

func main() {
	token := flag.String("token", "", "Telegram api token.")
	flag.Parse()
	if *token == "" {
		log.Fatal("token is empty")
	}
	bot, err := tgbotapi.NewBotAPI(*token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if msg, err := process(update); err == nil {
			bot.Send(msg)
		}
	}
}

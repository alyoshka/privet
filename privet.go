package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	fatherName = "just_alyosha"
	privet     = "привет"
	sad        = "/sad"
	come       = "к ноге"
	goodBoy    = "🐕"
	leave      = "Ну ты чего? Нормально же общались"
	greet      = "Привет, %s, расскажи пару слов о себе и о проекте, о котором нам почти не рассказывают!"
	dontBeSad  = "%s, главное - не расстраиваться!"
)

var myName string

func getName(u *tgbotapi.User) string {
	if u.FirstName != "" {
		return u.FirstName
	}
	return u.UserName
}

func process(update tgbotapi.Update) (answer tgbotapi.MessageConfig, err error) {
	if update.Message == nil {
		err = errors.New("empty")
		return
	}

	switch {
	// somebody leaves group
	case update.Message.LeftChatMember != nil:
		return tgbotapi.NewMessage(update.Message.Chat.ID, leave), nil

	// newby
	case update.Message.NewChatMember != nil:
		if update.Message.NewChatMember.UserName == myName {
			return tgbotapi.NewMessage(update.Message.Chat.ID, privet), nil
		}
		message := fmt.Sprintf(greet, getName(update.Message.NewChatMember))
		return tgbotapi.NewMessage(update.Message.Chat.ID, message), nil

	// somebody is sad
	case strings.Contains(update.Message.Text, sad):
		name := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(update.Message.Text, sad), "@"+myName))
		if name == "" || name == "@"+myName {
			name = getName(update.Message.From)
		}
		return tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(dontBeSad, name)), nil
	// father
	case strings.Contains(strings.ToLower(strings.TrimSpace(update.Message.Text)), come) && update.Message.From.UserName == fatherName:
		return tgbotapi.NewMessage(update.Message.Chat.ID, goodBoy), nil

	// say hello
	case strings.Contains(strings.ToLower(strings.TrimSpace(update.Message.Text)), privet):
		var msg tgbotapi.MessageConfig
		if update.Message.From.UserName == "ainomc" {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Миша, отстань")
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, privet)
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

	myName = bot.Self.UserName
	log.Printf("Authorized on account %s", myName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if msg, err := process(update); err == nil {
			bot.Send(msg)
		}
	}
}

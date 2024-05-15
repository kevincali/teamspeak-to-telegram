package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"
)

const (
	tgPrefix = "[Telegram]\t"
)

// newTelegramBot creates a Telegram bot using the Telegram API
func (tgConfig *Telegram) newTelegramBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(tgConfig.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = false

	log.Printf("%s authorized on account %s", tgPrefix, bot.Self.UserName)
	return bot
}

// initMessage sets a messageId in the config it it's not yet present
// it does so by sending a message to the chat and pinning it afterwards
func (config *Config) initMessage(bot *tgbotapi.BotAPI, configPath string) {
	// check if we already have a messageId specified
	if config.Telegram.MessageId == 0 {
		log.Printf("no messageId specified in config")

		// send message
		initChattable := tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:              config.Telegram.ChatId,
				DisableNotification: true,
			},
			Text: "init",
		}
		initMsg, err := bot.Send(initChattable)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s sent init message", tgPrefix)

		// save messageId to config
		config.Telegram.MessageId = initMsg.MessageID
		yamlData, err := yaml.Marshal(config)
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile(configPath, yamlData, 0644)
		if err != nil {
			log.Printf("%s unable to save messageId %d to config, %s", tgPrefix, config.Telegram.MessageId, err)
		} else {
			log.Printf("%s saved messageId %d to config", tgPrefix, config.Telegram.MessageId)
		}

		// pin message
		pinConfig := tgbotapi.PinChatMessageConfig{
			ChatID:              config.Telegram.ChatId,
			ChannelUsername:     "",
			MessageID:           config.Telegram.MessageId,
			DisableNotification: true,
		}
		bot.Send(pinConfig)
		log.Printf("%s pinned message", tgPrefix)
	}
}

// updateMessage edits the message to display the currently online users
func (tgConfig *Telegram) updateMessage(bot *tgbotapi.BotAPI, onlineUsers []string) {
	content := strings.Join(onlineUsers, tgConfig.Separator)

	if len(onlineUsers) == 0 {
		content = tgConfig.ZeroUsers
	}

	edit := tgbotapi.NewEditMessageText(tgConfig.ChatId, tgConfig.MessageId, content)
	message, err := bot.Send(edit)
	if err != nil {
		// don't log expected error
		if strings.Contains(err.Error(), "exactly the same") {
			return
		}
		log.Printf("%s unable to update message, %s", tgPrefix, err)
		return
	}

	log.Printf("%s updated message with online users: [%s]", tgPrefix, content)

	if tgConfig.UpdateTitle {
		tgConfig.updateTitle(bot, message.Chat.Title, len(onlineUsers))
	}
}

// updateTitle prepends the chat title with the current amount of online users
func (tgConfig *Telegram) updateTitle(bot *tgbotapi.BotAPI, originalTitle string, userAmount int) {
	// remove online users from original title
	titleFields := strings.Fields(originalTitle)
	if _, err := strconv.Atoi(titleFields[0]); err == nil {
		originalTitle = strings.Join(titleFields[1:], " ")
	}

	newTitle := originalTitle
	if userAmount != 0 {
		newTitle = fmt.Sprintf("%d %s", userAmount, originalTitle)
	}

	params := tgbotapi.Params{}
	params.AddFirstValid("chat_id", tgConfig.ChatId)
	params.AddBool("disable_notification", true)
	params["title"] = newTitle

	_, err := bot.MakeRequest("setChatTitle", params)
	if err != nil {
		log.Printf("%s unable to update chat title, %s", tgPrefix, err)
	}

	updates, err := bot.GetUpdates(tgbotapi.UpdateConfig{})
	if err != nil {
		log.Printf("%s unable to receive updates, %s", tgPrefix, err)
	}

	update := updates[len(updates)-1]
	if update.Message.From.UserName == bot.Self.UserName {
		deleteChattable := tgbotapi.NewDeleteMessage(tgConfig.ChatId, update.Message.MessageID)
		_, err = bot.Request(deleteChattable)
		if err != nil {
			log.Printf("%s unable to delete chat title update message, %s", tgPrefix, err)
		}
	}
}

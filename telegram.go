package main

import (
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"
)

// newTelegramBot creates a Telegram bot using the Telegram API
func (tgConfig *Telegram) newTelegramBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(tgConfig.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = false

	log.Printf("[Telegram]\t authorized on account %s", bot.Self.UserName)
	return bot
}

// initMessage sets a messageId in the config it it's not yet present
// it does so by sending a message to the chat and pinning it afterwards
func (config *Config) initMessage(bot *tgbotapi.BotAPI) {
	// check if we already have a messageId specified
	if config.Telegram.MessageId == 0 {
		log.Printf("no messageId specified in %s", configFile)

		// send message
		initChattable := tgbotapi.NewMessage(config.Telegram.ChatId, "init")
		initMsg, err := bot.Send(initChattable)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("[Telegram]\t sent message")

		// save messageId to config
		config.Telegram.MessageId = initMsg.MessageID
		yamlData, err := yaml.Marshal(config)
		if err != nil {
			log.Fatal(err)
		}
		os.WriteFile(configFile, yamlData, 0644)
		log.Printf("[Telegram]\t saved messageId to %s", configFile)

		// pin message
		pinConfig := tgbotapi.PinChatMessageConfig{
			ChatID:              config.Telegram.ChatId,
			ChannelUsername:     "",
			MessageID:           config.Telegram.MessageId,
			DisableNotification: false,
		}
		bot.Send(pinConfig)
		log.Println("[Telegram]\t pinned message")
	}
}

// updateMessage edits the message to display the currently online users
func (tgConfig *Telegram) updateMessage(bot *tgbotapi.BotAPI, onlineUsers []string) {
	content := strings.Join(onlineUsers, tgConfig.Separator)

	if len(onlineUsers) == 0 {
		content = tgConfig.ZeroUsers
	}

	edit := tgbotapi.NewEditMessageText(tgConfig.ChatId, tgConfig.MessageId, content)
	_, err := bot.Send(edit)
	if err != nil {
		// don't log expected error
		if strings.Contains(err.Error(), "exactly the same") {
			return
		}
		log.Printf("[Telegram]\t unable to update message, %s", err)
		return
	}

	log.Printf("[Telegram]\t updated message with online users: [%s]", content)
}

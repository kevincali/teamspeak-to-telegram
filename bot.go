package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"
)

const (
	configFile = "config.yaml"
)

type Config struct {
	BotToken  string `yaml:"botToken"`
	ChatId    int64  `yaml:"chatId"`
	MessageId int    `yaml:"messageId"`
}

func loadConfig() Config {
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var config Config

	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Fatal(err)
	}

	return config
}

func (config *Config) initMessage(bot *tgbotapi.BotAPI) {
	if config.MessageId == 0 {
		log.Printf("no messageId specified in %s", configFile)

		initChattable := tgbotapi.NewMessage(config.ChatId, "init")
		// send message
		initMsg, err := bot.Send(initChattable)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("sent message")

		config.MessageId = initMsg.MessageID
		// save messageId to config
		yamlData, err := yaml.Marshal(config)
		if err != nil {
			log.Fatal(err)
		}
		os.WriteFile(configFile, yamlData, 0644)
		log.Printf("saved messageId to %s", configFile)

		// pin message
		pinConfig := tgbotapi.PinChatMessageConfig{
			ChatID:              config.ChatId,
			ChannelUsername:     "",
			MessageID:           config.MessageId,
			DisableNotification: false,
		}
		bot.Send(pinConfig)
		log.Println("pinned message")
	}
}

func (config *Config) updateMessage(bot *tgbotapi.BotAPI, content string) {
	edit := tgbotapi.NewEditMessageText(config.ChatId, config.MessageId, content)
	_, err := bot.Send(edit)
	if err != nil {
		log.Printf("unable to update message, %s", err)
		return
	}
	log.Println("updated message")
}

func main() {
	config := loadConfig()

	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	// bot.Debug = true

	log.Printf("authorized on account %s", bot.Self.UserName)

	config.initMessage(bot)
	config.updateMessage(bot, "b")
}

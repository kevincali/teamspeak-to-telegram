package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const configFile = "config.yaml"

type Config struct {
	IntervalSeconds int       `yaml:"intervalSeconds"`
	TeamSpeak       TeamSpeak `yaml:"teamSpeak"`
	Telegram        Telegram  `yaml:"telegram"`
}

type TeamSpeak struct {
	FavoriteUsers []string    `yaml:"favoriteUsers"`
	Address       string      `yaml:"address"`
	ServerQuery   ServerQuery `yaml:"serverQuery"`
}

type ServerQuery struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Telegram struct {
	BotToken  string `yaml:"botToken"`
	ChatId    int64  `yaml:"chatId"`
	MessageId int    `yaml:"messageId"`
	Separator string `yaml:"separator"`
	ZeroUsers string `yaml:"zeroUsers"`
}

// loadConfig reads the config file
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

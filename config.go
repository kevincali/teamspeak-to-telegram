package main

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

const configFile = "config.yaml"

type Config struct {
	IntervalSeconds int       `yaml:"intervalSeconds" validate:"required"`
	TeamSpeak       TeamSpeak `yaml:"teamSpeak"`
	Telegram        Telegram  `yaml:"telegram"`
}

type TeamSpeak struct {
	FavoriteUsers []string    `yaml:"favoriteUsers"`
	Address       string      `yaml:"address" validate:"required"`
	ServerQuery   ServerQuery `yaml:"serverquery"`
}

type ServerQuery struct {
	Username string `yaml:"username" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}

type Telegram struct {
	BotToken  string `yaml:"botToken" validate:"required"`
	ChatId    int64  `yaml:"chatId" validate:"required"`
	MessageId int    `yaml:"messageId" validate:"required"`
	Separator string `yaml:"separator" validate:"required"`
	ZeroUsers string `yaml:"zeroUsers" validate:"required"`
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

func (config *Config) validate() {
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		validationErrors := err.(validator.ValidationErrors)

		log.Fatalf("missing values in %s\n%s", configFile, validationErrors)
	}
}
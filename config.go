package main

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Config struct {
	IntervalSeconds int         `yaml:"intervalSeconds" validate:"required"`
	TeamSpeak3      *TeamSpeak3 `yaml:"teamSpeak3"`
	TeamSpeak6      *TeamSpeak6 `yaml:"teamSpeak6"`
	Telegram        Telegram    `yaml:"telegram"`
}

type TeamSpeak3 struct {
	Address         string   `yaml:"address" validate:"required"`
	Username        string   `yaml:"username" validate:"required"`
	Password        string   `yaml:"password" validate:"required"`
	VirtualServerId string   `yaml:"virtualServerId" validate:"required"`
	FavoriteUsers   []string `yaml:"favoriteUsers"`
}

type TeamSpeak6 struct {
	Address       string   `yaml:"address" validate:"required"`
	ApiKey        string   `yaml:"apiKey" validate:"required"`
	FavoriteUsers []string `yaml:"favoriteUsers"`
}

type Telegram struct {
	BotToken    string `yaml:"botToken" validate:"required"`
	ChatId      int64  `yaml:"chatId" validate:"required"`
	MessageId   int    `yaml:"messageId"`
	Separator   string `yaml:"separator" validate:"required"`
	ZeroUsers   string `yaml:"zeroUsers" validate:"required"`
	UpdateTitle bool   `yaml:"updateTitle"`
}

func loadConfig(configPath string) Config {
	file, err := os.ReadFile(configPath)
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

	hasTS3 := config.TeamSpeak3 != nil
	hasTS6 := config.TeamSpeak6 != nil

	if hasTS3 && hasTS6 {
		log.Fatal("cannot use both teamSpeak3 and teamSpeak6 at the same time, please use only one")
	}

	if !hasTS3 && !hasTS6 {
		log.Fatal("please specify either teamSpeak3 or teamSpeak6 in config")
	}

	if hasTS3 {
		if err := validate.Struct(config.TeamSpeak3); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			log.Fatalf("missing values in teamSpeak3 config\n%s", validationErrors)
		}
	}

	if hasTS6 {
		if err := validate.Struct(config.TeamSpeak6); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			log.Fatalf("missing values in teamSpeak6 config\n%s", validationErrors)
		}
	}

	if err := validate.Struct(config); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		log.Fatalf("missing values in config\n%s", validationErrors)
	}
}

package main

import (
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	IntervalSeconds int `validate:"required"`
	TeamSpeak3      *TeamSpeak3
	TeamSpeak6      *TeamSpeak6
	Telegram        Telegram
}

type TeamSpeak3 struct {
	Address         string `validate:"required"`
	Username        string `validate:"required"`
	Password        string `validate:"required"`
	VirtualServerId string `validate:"required"`
	FavoriteUsers   []string
}

type TeamSpeak6 struct {
	Address       string `validate:"required"`
	ApiKey        string `validate:"required"`
	FavoriteUsers []string
}

type Telegram struct {
	BotToken  string `validate:"required"`
	ChatId    int64  `validate:"required"`
	MessageId int
	Separator string `validate:"required"`
	ZeroUsers string `validate:"required"`
}

func loadConfig() Config {
	config := Config{
		IntervalSeconds: 60,

		Telegram: Telegram{
			Separator: " | ",
			ZeroUsers: `¯\_(ツ)_/¯`,
		},
	}

	if v := os.Getenv("TS_POLLING_INTERVAL_SECONDS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			slog.Error("invalid TS_POLLING_INTERVAL_SECONDS", "error", err)
			os.Exit(1)
		}
		config.IntervalSeconds = n
	}

	// TeamSpeak3
	if os.Getenv("TS3_ADDRESS") != "" {
		config.TeamSpeak3 = &TeamSpeak3{
			Address:         os.Getenv("TS3_ADDRESS"),
			Username:        os.Getenv("TS3_USERNAME"),
			Password:        os.Getenv("TS3_PASSWORD"),
			VirtualServerId: os.Getenv("TS3_VIRTUAL_SERVER_ID"),
			FavoriteUsers:   parseCommaSeparated(os.Getenv("TS_FAVORITE_USERS")),
		}
	}

	// TeamSpeak6
	if os.Getenv("TS6_ADDRESS") != "" {
		config.TeamSpeak6 = &TeamSpeak6{
			Address:       os.Getenv("TS6_ADDRESS"),
			ApiKey:        os.Getenv("TS6_API_KEY"),
			FavoriteUsers: parseCommaSeparated(os.Getenv("TS_FAVORITE_USERS")),
		}
	}

	// Telegram
	config.Telegram.BotToken = os.Getenv("TELEGRAM_BOT_TOKEN")

	if v := os.Getenv("TELEGRAM_CHAT_ID"); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			slog.Error("invalid TELEGRAM_CHAT_ID", "error", err)
			os.Exit(1)
		}
		config.Telegram.ChatId = n
	}

	if v := os.Getenv("TELEGRAM_SEPARATOR"); v != "" {
		config.Telegram.Separator = v
	}

	if v := os.Getenv("TELEGRAM_ZERO_USERS"); v != "" {
		config.Telegram.ZeroUsers = v
	}

	return config
}

func parseCommaSeparated(s string) []string {
	if s == "" {
		return nil
	}
	var result []string
	for part := range strings.SplitSeq(s, ",") {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func (config *Config) validate() {
	validate := validator.New()

	hasTS3 := config.TeamSpeak3 != nil
	hasTS6 := config.TeamSpeak6 != nil

	if hasTS3 && hasTS6 {
		slog.Error("cannot use both teamspeak3 and teamspeak6 at the same time, please use only one")
		os.Exit(1)
	}

	if !hasTS3 && !hasTS6 {
		slog.Error("please set either TS3_ADDRESS or TS6_ADDRESS")
		os.Exit(1)
	}

	if hasTS3 {
		if err := validate.Struct(config.TeamSpeak3); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			slog.Error("missing required TS3_* environment variables", "scope", "TeamSpeak3", "errors", validationErrors.Error())
			os.Exit(1)
		}
	}

	if hasTS6 {
		if err := validate.Struct(config.TeamSpeak6); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			slog.Error("missing required TS6_* environment variables", "scope", "TeamSpeak6", "errors", validationErrors.Error())
			os.Exit(1)
		}
	}

	if err := validate.Struct(config); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		slog.Error("missing required environment variables", "errors", validationErrors.Error())
		os.Exit(1)
	}
}

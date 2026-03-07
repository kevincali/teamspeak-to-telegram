package main

import (
	"os"
	"time"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}
	config := loadConfig(configPath)
	config.validate()

	telegramBot := config.Telegram.newTelegramBot()
	config.initMessage(telegramBot, configPath)

	var onlineUsers []string

	if config.TeamSpeak3 != nil {
		tsConn := config.TeamSpeak3.newTeamSpeakConn()
		for {
			onlineUsers = config.TeamSpeak3.getOnlineUsers(tsConn)
			config.Telegram.updateMessage(telegramBot, onlineUsers)
			time.Sleep(time.Duration(config.IntervalSeconds) * time.Second)
		}
	}

	if config.TeamSpeak6 != nil {
		tsConn := config.TeamSpeak6.newTeamSpeakConn()
		for {
			onlineUsers = config.TeamSpeak6.getOnlineUsers(tsConn)
			config.Telegram.updateMessage(telegramBot, onlineUsers)
			time.Sleep(time.Duration(config.IntervalSeconds) * time.Second)
		}
	}
}

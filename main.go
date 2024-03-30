package main

import (
	"log"
	"os"
	"time"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("please specify a CONFIG_PATH envvar")
	}
	config := loadConfig(configPath)
	config.validate()

	telegramBot := config.Telegram.newTelegramBot()
	config.initMessage(telegramBot, configPath)

	teamspeakConn := config.TeamSpeak.newTeamSpeakConn()

	for {
		onlineUsers := config.TeamSpeak.getOnlineUsers(teamspeakConn)
		config.Telegram.updateMessage(telegramBot, onlineUsers)

		time.Sleep(time.Duration(config.IntervalSeconds) * time.Second)
	}
}

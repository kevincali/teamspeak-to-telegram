package main

import (
	"log"
	"time"
)

func main() {
	config := loadConfig()
	config.validate()

	// Load persisted state
	state, err := loadState(stateFile)
	if err != nil {
		log.Printf("warning: failed to load state: %s", err)
	}
	if state != nil && state.TelegramMessageId != 0 {
		config.Telegram.MessageId = state.TelegramMessageId
		log.Printf("loaded message ID %d from state file", state.TelegramMessageId)
	}

	telegramBot := config.Telegram.newTelegramBot()
	config.initMessage(telegramBot)

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

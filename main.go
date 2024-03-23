package main

import "time"

func main() {
	config := loadConfig()
	config.validate()

	telegramBot := config.Telegram.newTelegramBot()
	config.initMessage(telegramBot)

	teamspeakConn := config.TeamSpeak.newTeamSpeakConn()

	for {
		onlineUsers := config.TeamSpeak.getOnlineUsers(teamspeakConn)
		config.Telegram.updateMessage(telegramBot, onlineUsers)

		time.Sleep(time.Duration(config.IntervalSeconds) * time.Second)
	}
}

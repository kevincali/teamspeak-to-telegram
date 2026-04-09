package main

import (
	"context"
	"log/slog"
	"os"
	"time"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := loadConfig()
	config.validate()

	tgLog := slog.With("scope", "Telegram")

	state, err := loadState(stateFile)
	if err != nil {
		slog.Warn("failed to load state", "error", err)
	}
	if state != nil && state.TelegramMessageId != 0 {
		config.Telegram.MessageId = state.TelegramMessageId
		tgLog.Info("loaded message ID from state file", "message_id", state.TelegramMessageId)
	}

	telegramBot := config.Telegram.newTelegramBot(tgLog)
	config.initMessage(ctx, telegramBot, tgLog)

	var onlineUsers []string

	if config.TeamSpeak3 != nil {
		ts3Log := slog.With("scope", "TeamSpeak3")
		tsConn := config.TeamSpeak3.newTeamSpeakConn(ts3Log)
		for {
			onlineUsers = config.TeamSpeak3.getOnlineUsers(tsConn)
			config.Telegram.updateMessage(ctx, telegramBot, onlineUsers, tgLog)
			time.Sleep(time.Duration(config.IntervalSeconds) * time.Second)
		}
	}

	if config.TeamSpeak6 != nil {
		ts6Log := slog.With("scope", "TeamSpeak6")
		tsConn := config.TeamSpeak6.newTeamSpeakConn(ts6Log)
		for {
			onlineUsers = config.TeamSpeak6.getOnlineUsers(tsConn)
			config.Telegram.updateMessage(ctx, telegramBot, onlineUsers, tgLog)
			time.Sleep(time.Duration(config.IntervalSeconds) * time.Second)
		}
	}
}

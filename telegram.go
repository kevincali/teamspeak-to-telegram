package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/go-telegram/bot"
)

// newTelegramBot creates a Telegram bot using the Telegram API
func (tgConfig *Telegram) newTelegramBot(logger *slog.Logger) *bot.Bot {
	telegramBot, err := bot.New(tgConfig.BotToken)
	if err != nil {
		logger.Error("failed to create bot", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	self, err := telegramBot.GetMe(ctx)
	if err != nil {
		logger.Info("authorized bot", "bot_id", telegramBot.ID(), "username_error", err)
	} else {
		logger.Info("authorized bot", "username", self.Username)
	}

	return telegramBot
}

// initMessage sets a messageId in the config if it's not yet present
// by sending a message to the chat and pinning it afterwards
func (config *Config) initMessage(ctx context.Context, telegramBot *bot.Bot, logger *slog.Logger) {
	if config.Telegram.MessageId != 0 {
		return
	}

	logger.Info("no message ID in state, creating new message")

	initMsg, err := telegramBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:              config.Telegram.ChatId,
		DisableNotification: true,
		Text:                "init",
	})
	if err != nil {
		logger.Error("failed to send init message", "error", err, "chat_id", config.Telegram.ChatId)
		os.Exit(1)
	}

	config.Telegram.MessageId = initMsg.ID
	logger.Info("created message", "message_id", initMsg.ID)

	_, err = telegramBot.PinChatMessage(ctx, &bot.PinChatMessageParams{
		ChatID:              config.Telegram.ChatId,
		MessageID:           config.Telegram.MessageId,
		DisableNotification: true,
	})
	if err != nil {
		logger.Warn("failed to pin message", "error", err)
	} else {
		logger.Info("pinned message")
	}

	err = saveState(stateFile, &State{TelegramMessageId: config.Telegram.MessageId})
	if err != nil {
		logger.Warn("failed to save state, message will be re-created on next restart", "error", err)
	} else {
		logger.Info("saved message ID to state file")
	}
}

// updateMessage edits the message to display the currently online users
func (tgConfig *Telegram) updateMessage(ctx context.Context, telegramBot *bot.Bot, onlineUsers []string, logger *slog.Logger) {
	content := strings.Join(onlineUsers, tgConfig.Separator)

	if len(onlineUsers) == 0 {
		content = tgConfig.ZeroUsers
	}

	_, err := telegramBot.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    tgConfig.ChatId,
		MessageID: tgConfig.MessageId,
		Text:      content,
	})
	if err != nil {
		if strings.Contains(err.Error(), "message is not modified") || strings.Contains(err.Error(), "exactly the same") {
			return
		}
		logger.Error("failed to update message", "error", err, "chat_id", tgConfig.ChatId, "message_id", tgConfig.MessageId)
		return
	}

	logger.Info("updated message", "online_users", onlineUsers, "count", len(onlineUsers))
}

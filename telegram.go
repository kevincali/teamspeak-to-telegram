package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/go-telegram/bot"
)

const (
	tgPrefix = "[Telegram]\t"
)

// newTelegramBot creates a Telegram bot using the Telegram API
func (tgConfig *Telegram) newTelegramBot() *bot.Bot {
	telegramBot, err := bot.New(tgConfig.BotToken)
	if err != nil {
		log.Fatalf("%s failed to create bot: %s (botToken=%s)", tgPrefix, err, tgConfig.BotToken)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	self, err := telegramBot.GetMe(ctx)
	if err != nil {
		log.Printf("%s authorized bot with ID %d (failed to fetch username: %s)", tgPrefix, telegramBot.ID(), err)
	} else {
		log.Printf("%s authorized on account %s", tgPrefix, self.Username)
	}

	return telegramBot
}

// initMessage sets a messageId in the config it it's not yet present
// it does so by sending a message to the chat and pinning it afterwards
func (config *Config) initMessage(ctx context.Context, telegramBot *bot.Bot) {
	// check if we already have a messageId specified
	if config.Telegram.MessageId == 0 {
		log.Printf("%s no message ID in state, creating new message", tgPrefix)

		// send message
		initMsg, err := telegramBot.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:              config.Telegram.ChatId,
			DisableNotification: true,
			Text:                "init",
		})
		if err != nil {
			log.Fatalf("%s failed to send init message: %s (botToken=%s, chatId=%d)", tgPrefix, err, config.Telegram.BotToken, config.Telegram.ChatId)
		}
		log.Printf("%s sent init message", tgPrefix)

		config.Telegram.MessageId = initMsg.ID
		log.Printf("%s created message with ID %d", tgPrefix, initMsg.ID)

		// pin message
		_, err = telegramBot.PinChatMessage(ctx, &bot.PinChatMessageParams{
			ChatID:              config.Telegram.ChatId,
			MessageID:           config.Telegram.MessageId,
			DisableNotification: true,
		})
		if err != nil {
			log.Printf("%s failed to pin message: %s", tgPrefix, err)
		} else {
			log.Printf("%s pinned message", tgPrefix)
		}

		// persist message ID to state file
		err = saveState(stateFile, &State{TelegramMessageId: config.Telegram.MessageId})
		if err != nil {
			log.Printf("%s failed to save state: %s (message will be re-created on next restart)", tgPrefix, err)
		} else {
			log.Printf("%s saved message ID to state file", tgPrefix)
		}
	}
}

// updateMessage edits the message to display the currently online users
func (tgConfig *Telegram) updateMessage(ctx context.Context, telegramBot *bot.Bot, onlineUsers []string) {
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
		// don't log expected error
		if strings.Contains(err.Error(), "message is not modified") || strings.Contains(err.Error(), "exactly the same") {
			return
		}
		log.Printf("%s failed to update message: %s (chatId=%d, messageId=%d)", tgPrefix, err, tgConfig.ChatId, tgConfig.MessageId)
		return
	}

	log.Printf("%s updated message with online users: [%s]", tgPrefix, content)
}

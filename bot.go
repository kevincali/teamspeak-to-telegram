package main

import (
	"cmp"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ziutek/telnet"
	"gopkg.in/yaml.v3"
)

const configFile = "config.yaml"

type Config struct {
	IntervalSeconds int       `yaml:"intervalSeconds"`
	TeamSpeak       TeamSpeak `yaml:"teamSpeak"`
	Telegram        Telegram  `yaml:"telegram"`
}

type TeamSpeak struct {
	FavoriteUsers []string    `yaml:"favoriteUsers"`
	Address       string      `yaml:"address"`
	ServerQuery   ServerQuery `yaml:"serverQuery"`
}

type ServerQuery struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Telegram struct {
	BotToken  string `yaml:"botToken"`
	ChatId    int64  `yaml:"chatId"`
	MessageId int    `yaml:"messageId"`
	Separator string `yaml:"separator"`
	ZeroUsers string `yaml:"zeroUsers"`
}

func main() {
	config := loadConfig()

	telegramBot := config.Telegram.newTelegramBot()
	config.initMessage(telegramBot)

	teamspeakConn := config.TeamSpeak.newTeamSpeakConn()

	for {
		onlineUsers := config.TeamSpeak.getOnlineUsers(teamspeakConn)
		config.Telegram.updateMessage(telegramBot, onlineUsers)

		time.Sleep(time.Duration(config.IntervalSeconds) * time.Second)
	}
}

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

// TeamSpeak

func (tsConfig *TeamSpeak) newTeamSpeakConn() *telnet.Conn {
	conn, err := telnet.DialTimeout("tcp", tsConfig.Address, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	// skip until banner end is reached
	err = conn.SkipUntil("command.")
	if err != nil {
		log.Fatal(err)
	}

	conn.Write([]byte(fmt.Sprintf("login %s %s\n", tsConfig.ServerQuery.Username, tsConfig.ServerQuery.Password)))
	conn.Write([]byte("use 3\n"))

	// skip first message
	err = conn.SkipUntil("msg=ok")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[TeamSpeak]\t connected to telnet")
	return conn
}

func (tsConfig *TeamSpeak) getOnlineUsers(conn *telnet.Conn) []string {
	// login, select server and return clientlist
	conn.Write([]byte("clientlist\n"))
	time.Sleep(1000 * time.Millisecond)

	// get clientlist
	result, err := conn.ReadUntil("msg=ok")
	if err != nil {
		log.Fatal(err)
	}
	// log.Println("[TeamSpeak]\t received clientlist")

	// parse clientlist
	clients := strings.Split(string(result), "|")
	var users []string
	for _, client := range clients {
		username := strings.Fields(client)[3]
		username = strings.Split(username, "=")[1]
		username = strings.ReplaceAll(username, "\\s", " ")

		if len(tsConfig.FavoriteUsers) != 0 {
			if slices.Contains(tsConfig.FavoriteUsers, username) {
				users = append(users, username)
			}
			continue
		}
		// skip serverquery client
		if username == "Unknown" {
			continue
		}
		users = append(users, username)
	}

	// sort by name length to display as much as possible on smaller screens
	slices.SortFunc(users, func(a, b string) int {
		return cmp.Compare(len(a), len(b))
	})

	// log.Printf("[TeamSpeak]\t %d connected users", len(users))

	return users
}

// Telegram

func (tgConfig *Telegram) newTelegramBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(tgConfig.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = false

	log.Printf("[Telegram]\t authorized on account %s", bot.Self.UserName)
	return bot
}

func (config *Config) initMessage(bot *tgbotapi.BotAPI) {
	// check if we already have a messageId specified
	if config.Telegram.MessageId == 0 {
		log.Printf("no messageId specified in %s", configFile)

		// send message
		initChattable := tgbotapi.NewMessage(config.Telegram.ChatId, "init")
		initMsg, err := bot.Send(initChattable)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("[Telegram]\t sent message")

		// save messageId to config
		config.Telegram.MessageId = initMsg.MessageID
		yamlData, err := yaml.Marshal(config)
		if err != nil {
			log.Fatal(err)
		}
		os.WriteFile(configFile, yamlData, 0644)
		log.Printf("[Telegram]\t saved messageId to %s", configFile)

		// pin message
		pinConfig := tgbotapi.PinChatMessageConfig{
			ChatID:              config.Telegram.ChatId,
			ChannelUsername:     "",
			MessageID:           config.Telegram.MessageId,
			DisableNotification: false,
		}
		bot.Send(pinConfig)
		log.Println("[Telegram]\t pinned message")
	}
}

func (tgConfig *Telegram) updateMessage(bot *tgbotapi.BotAPI, onlineUsers []string) {
	content := strings.Join(onlineUsers, tgConfig.Separator)

	if len(onlineUsers) == 0 {
		content = tgConfig.ZeroUsers
	}

	edit := tgbotapi.NewEditMessageText(tgConfig.ChatId, tgConfig.MessageId, content)
	_, err := bot.Send(edit)
	if err != nil {
		if strings.Contains(err.Error(), "exactly the same") {
			// don't log expected error
			return
		}
		log.Printf("[Telegram]\t unable to update message, %s", err)
		return
	}

	log.Printf("[Telegram]\t updated message with online users: [%s]", content)
}

package main

import (
	"cmp"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/ziutek/telnet"
)

type TS3Connection struct {
	conn   *telnet.Conn
	logger *slog.Logger
}

func (tsConfig *TeamSpeak3) newTeamSpeakConn(logger *slog.Logger) *TS3Connection {
	conn, err := telnet.DialTimeout("tcp", tsConfig.Address, 5*time.Second)
	if err != nil {
		logger.Error("failed to connect", "address", tsConfig.Address, "error", err)
		os.Exit(1)
	}

	// skip until banner end is reached
	if err = conn.SkipUntil("command."); err != nil {
		logger.Error("failed to read banner", "error", err)
		os.Exit(1)
	}

	fmt.Fprintf(conn, "login %s %s\n", tsConfig.Username, tsConfig.Password)
	if err = conn.SkipUntil("msg=ok"); err != nil {
		logger.Error("failed to authenticate", "error", err)
		os.Exit(1)
	}

	fmt.Fprintf(conn, "use %s\n", tsConfig.VirtualServerId)
	if err = conn.SkipUntil("msg=ok"); err != nil {
		logger.Error("failed to select virtual server", "virtual_server_id", tsConfig.VirtualServerId, "error", err)
		os.Exit(1)
	}

	logger.Info("connected", "address", tsConfig.Address)
	return &TS3Connection{conn: conn, logger: logger}
}

func (tsConfig *TeamSpeak3) getOnlineUsers(tsConn *TS3Connection) []string {
	tsConn.conn.Write([]byte("clientlist\n"))
	time.Sleep(1 * time.Second)

	// get clientlist
	result, err := tsConn.conn.ReadUntil("msg=ok")
	if err != nil {
		tsConn.logger.Error("failed to read client list", "error", err)
		os.Exit(1)
	}

	// parse clientlist
	clients := strings.Split(string(result), "|")
	var users []string
	for _, client := range clients {
		fields := strings.Fields(client)

		if len(fields) < 5 {
			continue
		}

		if clientType := strings.Split(fields[4], "=")[1]; clientType != "0" {
			continue
		}

		username := strings.Split(fields[3], "=")[1]
		username = strings.ReplaceAll(username, "\\s", " ")

		if len(tsConfig.FavoriteUsers) != 0 {
			if slices.Contains(tsConfig.FavoriteUsers, username) {
				users = append(users, username)
			}
			continue
		}
		users = append(users, username)
	}

	// sort by name length to display as much as possible on smaller screens
	slices.SortFunc(users, func(a, b string) int {
		return cmp.Compare(len(a), len(b))
	})

	return users
}

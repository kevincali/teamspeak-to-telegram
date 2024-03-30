package main

import (
	"cmp"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/ziutek/telnet"
)

// newTeamSpeakConn initiates the telnet connection to a TeamSpeak server
func (tsConfig *TeamSpeak) newTeamSpeakConn() *telnet.Conn {
	conn, err := telnet.DialTimeout("tcp", tsConfig.Address, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	// skip until banner end is reached
	if err = conn.SkipUntil("command."); err != nil {
		log.Fatal(err)
	}

	conn.Write([]byte(fmt.Sprintf("login %s %s\n", tsConfig.Username, tsConfig.Password)))
	if err = conn.SkipUntil("msg=ok"); err != nil {
		log.Fatal(err)
	}

	conn.Write([]byte(fmt.Sprintf("use %s\n", tsConfig.VirtualServerId)))
	if err = conn.SkipUntil("msg=ok"); err != nil {
		log.Fatal(err)
	}

	log.Println("[TeamSpeak]\t connected to telnet")
	return conn
}

// getOnlineUsers checks the online users and filters by favorites if they're defined
func (tsConfig *TeamSpeak) getOnlineUsers(conn *telnet.Conn) []string {
	// login, select server and return clientlist
	conn.Write([]byte("clientlist\n"))
	time.Sleep(1 * time.Second)

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
		fields := strings.Fields(client)

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

	// log.Printf("[TeamSpeak]\t %d connected users", len(users))

	return users
}

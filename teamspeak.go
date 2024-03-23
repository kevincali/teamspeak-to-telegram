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

// getOnlineUsers checks the online users and filters by favorites if they're defined
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

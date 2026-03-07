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

const (
	ts3Prefix = "[TeamSpeak3]\t"
)

type TS3Connection struct {
	conn *telnet.Conn
}

func (tsConfig *TeamSpeak3) newTeamSpeakConn() *TS3Connection {
	conn, err := telnet.DialTimeout("tcp", tsConfig.Address, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	// skip until banner end is reached
	if err = conn.SkipUntil("command."); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(conn, "login %s %s\n", tsConfig.Username, tsConfig.Password)
	if err = conn.SkipUntil("msg=ok"); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(conn, "use %s\n", tsConfig.VirtualServerId)
	if err = conn.SkipUntil("msg=ok"); err != nil {
		log.Fatal(err)
	}

	log.Printf("%s connected to telnet", ts3Prefix)
	return &TS3Connection{conn: conn}
}

func (tsConfig *TeamSpeak3) getOnlineUsers(tsConn *TS3Connection) []string {
	tsConn.conn.Write([]byte("clientlist\n"))
	time.Sleep(1 * time.Second)

	// get clientlist
	result, err := tsConn.conn.ReadUntil("msg=ok")
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("%s received clientlist", ts3Prefix)

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

	// log.Printf("%s %d connected users", ts3Prefix, len(users))
	return users
}

package main

import (
	"cmp"
	"encoding/json"
	"log"
	"net/http"
	"slices"
)

const (
	ts6Prefix = "[TeamSpeak6]\t"
)

type TS6Client struct {
	ClientNickname string `json:"client_nickname"`
	ClientType     string `json:"client_type"`
}

type TS6ResponseBody struct {
	Body   []TS6Client `json:"body"`
	Status TS6Status   `json:"status"`
}

type TS6Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type TS6Connection struct {
	address string
	apiKey  string
	client  *http.Client
}

func (tsConfig *TeamSpeak6) newTeamSpeakConn() *TS6Connection {
	log.Printf("%s connecting to %s", ts6Prefix, tsConfig.Address)

	return &TS6Connection{
		address: tsConfig.Address,
		apiKey:  tsConfig.ApiKey,
		client:  &http.Client{},
	}
}

func (tsConfig *TeamSpeak6) getOnlineUsers(tsConn *TS6Connection) []string {
	req, err := http.NewRequest("GET", tsConn.address+"/1/clientlist", nil)
	if err != nil {
		log.Printf("%s error creating request: %s", ts6Prefix, err)
		return nil
	}

	req.Header.Set("x-api-key", tsConn.apiKey)

	resp, err := tsConn.client.Do(req)
	if err != nil {
		log.Printf("%s error making request: %s", ts6Prefix, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("%s unexpected status code: %d", ts6Prefix, resp.StatusCode)
		return nil
	}

	var result TS6ResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("%s error decoding response: %s", ts6Prefix, err)
		return nil
	}

	if result.Status.Code != 0 {
		log.Printf("%s API error: %s", ts6Prefix, result.Status.Message)
		return nil
	}

	var users []string
	for _, client := range result.Body {
		if client.ClientType == "1" {
			continue
		}
		username := client.ClientNickname

		if len(tsConfig.FavoriteUsers) != 0 {
			if slices.Contains(tsConfig.FavoriteUsers, username) {
				users = append(users, username)
			}
			continue
		}
		users = append(users, username)
	}

	slices.SortFunc(users, func(a, b string) int {
		return cmp.Compare(len(a), len(b))
	})

	// log.Printf("%s %d connected users", ts6Prefix, len(users))
	return users
}

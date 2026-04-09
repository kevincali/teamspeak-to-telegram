package main

import (
	"cmp"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
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
	logger  *slog.Logger
}

func (tsConfig *TeamSpeak6) newTeamSpeakConn(logger *slog.Logger) *TS6Connection {
	u, err := url.Parse(tsConfig.Address)
	if err != nil || u.Port() == "" {
		logger.Error("address must include a port (e.g. http://host:10080)", "address", tsConfig.Address)
		os.Exit(1)
	}

	logger.Info("connecting", "address", tsConfig.Address)

	return &TS6Connection{
		address: tsConfig.Address,
		apiKey:  tsConfig.ApiKey,
		client:  &http.Client{},
		logger:  logger,
	}
}

func (tsConfig *TeamSpeak6) getOnlineUsers(tsConn *TS6Connection) []string {
	req, err := http.NewRequest("GET", tsConn.address+"/1/clientlist", nil)
	if err != nil {
		tsConn.logger.Error("failed to create request", "error", err)
		return nil
	}

	req.Header.Set("x-api-key", tsConn.apiKey)

	resp, err := tsConn.client.Do(req)
	if err != nil {
		tsConn.logger.Error("failed to make request", "error", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		tsConn.logger.Error("unexpected status code", "status", resp.StatusCode)
		return nil
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		tsConn.logger.Error("unexpected content type, check address includes the query port (e.g. :10080)", "content_type", contentType)
		os.Exit(1)
	}

	var result TS6ResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		tsConn.logger.Error("failed to decode response", "error", err)
		return nil
	}

	if result.Status.Code != 0 {
		tsConn.logger.Error("API error", "code", result.Status.Code, "message", result.Status.Message)
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

	return users
}

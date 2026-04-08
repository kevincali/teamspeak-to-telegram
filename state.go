package main

import (
	"encoding/json"
	"errors"
	"os"
)

const stateFile = "state.json"

type State struct {
	TelegramMessageId int `json:"telegram_message_id"`
}

func loadState(filePath string) (*State, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &State{}, nil
		}
		return nil, err
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func saveState(filePath string, state *State) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

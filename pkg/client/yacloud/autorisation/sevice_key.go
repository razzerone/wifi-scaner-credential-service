package autorisation

import (
	"encoding/json"
	"os"
)

type Key struct {
	ID         string `json:"id"`
	AccountID  string `json:"service_account_id"`
	PrivateKey string `json:"private_key"`
}

func GetKeyFromFile(path string) (*Key, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var key Key

	err = json.Unmarshal(content, &key)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

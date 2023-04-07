package autorisation

import (
	"encoding/json"
	"log"
	"os"
)

type Key struct {
	ID         string `json:"id"`
	AccountID  string `json:"service_account_id"`
	PrivateKey string `json:"private_key"`
}

func GetKeyFromFile(path string) *Key {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Panic("unable read key file")
		return nil
	}

	var key Key

	err = json.Unmarshal(content, &key)
	if err != nil {
		log.Panic("unable parse key file")
		return nil
	}

	return &key
}

package keyring

import (
	"log"

	"github.com/zalando/go-keyring"
)

const service = "gmail_checker"

func GetEntry(key string) string {
	// get password
	secret, err := keyring.Get(service, key)
	if err != nil {
		log.Fatal(err)
	}
	return secret
}

func SetEntry(key string, data string) bool {

	err := keyring.Set(service, key, data)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

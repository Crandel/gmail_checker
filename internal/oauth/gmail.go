package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Crandel/gmail/internal/keyring"
	"golang.org/x/oauth2"
)

const (
	tokKey    = "tokenKey"
	configKey = "configKey"
)

func SaveConfig(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	config := &oauth2.Config{}
	err = json.NewDecoder(f).Decode(config)
	if err != nil {
		return err
	}
	slog.Debug("Saving config to keyring")
	key := configKey + config.ClientID
	configByte, err := json.Marshal(config)
	if err != nil {
		slog.Debug("Error during marshalling config", slog.Any("error", err))
		return err
	}
	return keyring.SetEntry(key, string(configByte))
}

func GetConfig(clientID string) (*oauth2.Config, error) {
	config := &oauth2.Config{}
	key := configKey + clientID
	configStr, err := keyring.GetEntry(key)
	if err != nil {
		slog.Debug("can't get config from keyring", slog.Any("error", err))
		return nil, err
	}
	err = json.Unmarshal([]byte(configStr), config)
	return config, err
}

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	key := tokKey + config.ClientID
	tok, err := tokenFromKeyring(key)
	if err != nil {
		tok = getTokenFromWeb(config)
		err = saveToken(key, tok)
		if err != nil {
			slog.Debug("can't save token to keyring", slog.Any("error", err))
		}
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	slog.Info("Go to the following link in your browser then type", slog.String("authorization code", authURL))

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		slog.Debug("Unable to read authorization code", slog.Any("error", err))
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		slog.Debug("Unable to retrieve token from web", slog.Any("error", err))
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromKeyring(file string) (*oauth2.Token, error) {
	tok := &oauth2.Token{}
	tokString, err := keyring.GetEntry(tokKey)
	if err != nil {
		slog.Debug("can't get token from keyring", slog.Any("error", err))
		return nil, err
	}
	err = json.Unmarshal([]byte(tokString), tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(key string, token *oauth2.Token) error {
	slog.Debug("Saving credentials  to keyring")
	tokenByte, err := json.Marshal(token)
	if err != nil {
		slog.Debug("Error during marshalling token", slog.Any("error", err))
		return err
	}
	return keyring.SetEntry(key, string(tokenByte))
}

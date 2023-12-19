package gmail

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Crandel/gmail/internal/keyring"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

const (
	tokKey    = "gmailTokenKey"
	configKey = "gmailConfigKey"
)

type GmailUser struct {
	*oauth2.Token
	Alias string
}

func SaveConfig(path string) error {
	credJSON, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var config *oauth2.Config
	config, err = google.ConfigFromJSON(credJSON, gmail.GmailReadonlyScope)
	if err != nil {
		slog.Error("Error during reading config from file", slog.Any("error", err))
	}
	slog.Debug("Saving config to keyring", slog.String("client id", config.ClientID))
	key := configKey + config.ClientID
	configByte, err := json.Marshal(config)
	if err != nil {
		slog.Debug("Error during marshalling config", slog.Any("error", err))
		return err
	}
	return keyring.SetEntry(key, string(configByte))
}

func SaveToken(config *oauth2.Config, alias string) error {
	key := configKey + config.ClientID

	tok := getTokenFromWeb(config)
	gmailUser := GmailUser{
		tok,
		alias,
	}

	err := saveToken(key, gmailUser)
	if err != nil {
		slog.Debug("can't save token to keyring", slog.Any("error", err))
	}
	return err
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
func GetClient(config *oauth2.Config) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	key := tokKey + config.ClientID
	gmailUser, err := tokenFromKeyring(key)
	if err != nil {
		return nil, err
	}
	return config.Client(context.Background(), gmailUser.Token), err
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

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
func tokenFromKeyring(file string) (*GmailUser, error) {
	gmailUser := &GmailUser{}
	userString, err := keyring.GetEntry(tokKey)
	if err != nil {
		slog.Debug("can't get token from keyring", slog.Any("error", err))
		return nil, err
	}
	err = json.Unmarshal([]byte(userString), gmailUser)
	return gmailUser, err
}

// Saves a token to a file path.
func saveToken(key string, gmailUser GmailUser) error {
	slog.Debug("Saving credentials  to keyring")
	tokenByte, err := json.Marshal(gmailUser)
	if err != nil {
		slog.Debug("Error during marshalling token", slog.Any("error", err))
		return err
	}
	return keyring.SetEntry(key, string(tokenByte))
}

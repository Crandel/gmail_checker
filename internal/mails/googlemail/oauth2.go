package googlemail

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Crandel/gmail/internal/keyring"
	"golang.org/x/oauth2"
)

const (
	tokKey = "gmailTokenKey"
)

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(ctx context.Context, config *oauth2.Config) (*http.Client, error) {
	token, err := tokenFromKeyring(ctx, config)
	if err != nil {
		return nil, err
	}
	return config.Client(ctx, token), err
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	authURL, err := config.DeviceAuth(ctx, oauth2.AccessTypeOffline)
	if err != nil {
		slog.Debug("unable to create auth link", slog.Any("error", err))
		return nil
	}
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
func tokenFromKeyring(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	key := tokKey + config.ClientID
	token := &oauth2.Token{}
	userString, err := keyring.GetEntry(key)
	if err != nil {
		slog.Debug("can't get token from keyring", slog.Any("error", err))
		tok := getTokenFromWeb(ctx, config)

		err = saveToken(key, tok)
		if err != nil {
			slog.Debug("can't save token to keyring", slog.Any("error", err))
		}
		return nil, err
	}
	err = json.Unmarshal([]byte(userString), token)
	return token, err
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

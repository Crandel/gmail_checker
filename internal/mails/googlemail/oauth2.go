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
	return config.Client(ctx, token), nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	codeChan := startServer()
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	authCode := <-codeChan
	tok, err := config.Exchange(ctx, authCode)
	if err != nil {
		slog.Debug("Unable to retrieve token from web", slog.Any("error", err))
	}
	return tok, nil
}

// Retrieves a token from a local file.
func tokenFromKeyring(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	key := tokKey + config.ClientID
	token := &oauth2.Token{}
	userString, err := keyring.GetEntry(key)
	if err == nil {
		err = json.Unmarshal([]byte(userString), token)
		if err == nil && token.Valid() {
			return token, nil
		}
	}

	slog.Debug("can't get token from keyring", slog.Any("error", err), slog.Bool("token is valid", token.Valid()))
	tok, err := getTokenFromWeb(ctx, config)

	if err != nil && tok == nil {
		slog.Debug("can't get token from web", slog.Any("error", err), slog.Any("token", tok))
		return nil, err
	}
	err = saveToken(key, tok)
	if err != nil {
		slog.Debug("can't save token to keyring", slog.Any("error", err))
		return nil, err
	}
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
func startServer() chan string {
	codeChan := make(chan string)
	server := &http.Server{Addr: ":8080"}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		codeChan <- code
		fmt.Fprintf(w, "Authorization successful! You can close this window now.")
		go func() {
			server.Shutdown(context.Background())
		}()
	})

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			slog.Error("HTTP server ListenAndServe: %v", err)
			return
		}
	}()

	return codeChan
}

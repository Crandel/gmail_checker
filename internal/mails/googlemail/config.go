package googlemail

import (
	"fmt"
	"log/slog"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	libGmail "google.golang.org/api/gmail/v1"

	"github.com/Crandel/gmail/internal/file"
)

// GetConfig will return oauth2 config for clientID.
func GetConfig(clientID string) (*oauth2.Config, error) {
	cacheDir, err := file.GetCacheDir()
	if err != nil {
		return nil, fmt.Errorf("unable to get cache dir: %w", err)
	}
	credentialsFile := fmt.Sprintf("%s/%s.json", cacheDir, clientID)

	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %w", err)
	}

	config, err := google.ConfigFromJSON(b, libGmail.GmailMetadataScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %w", err)
	}
	slog.Debug("config", slog.Any("conf", config))

	return config, nil
}

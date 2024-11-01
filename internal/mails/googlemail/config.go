package googlemail

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Crandel/gmail/internal/file"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	libGmail "google.golang.org/api/gmail/v1"
)

func GetConfig(clientID string) (*oauth2.Config, error) {
	cacheDir, err := file.GetCacheDir()
	if err != nil {
		return nil, fmt.Errorf("unable to get cache dir: %v", err)
	}
	credentialsFile := fmt.Sprintf("%s/%s.json", cacheDir, clientID)

	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, libGmail.GmailMetadataScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}
	slog.Debug("config", slog.Any("conf", config))

	return config, nil
}

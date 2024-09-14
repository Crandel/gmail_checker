package googlemail

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Crandel/gmail/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

const credentialsName = "credentials.json"

var credentialsFile = fmt.Sprintf("%s/%s", config.ConfigDir, credentialsName)

func GetConfig(clientID, clientSecret string) (*oauth2.Config, error) {

	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailMetadataScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}
	slog.Debug("config", slog.Any("conf", config))

	return config, nil
}

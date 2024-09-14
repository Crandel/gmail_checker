package googlemail

import (
	"fmt"
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

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	return config, nil
}

package googlemail

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	libGmail "google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"github.com/Crandel/gmail/internal/accounts"
)

const (
	url     = "https://mail.google.com"
	service = "gmail_checker"
)

// var requiredLabels = []string{
// 	"INBOX",
// 	"TRASH",
// 	"SPAM",
// 	"Github",
// 	"Delivery",
// }

// CheckOnline will check google servers availability.
func CheckOnline() bool {
	// Check if domain online
	resp, err := http.Get(url)
	if err == nil || resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusMovedPermanently {
		return true
	}
	return false
}

// GetMailCount - new goroutine for checking emails
func GetGMailCount(ctx context.Context, channel chan<- string, acc accounts.Account, systemKeyring bool) {
	config, err := GetConfig(acc.ClientID)
	if err != nil {
		slog.Debug("can't load config from credentials", slog.Any("error", err))
		return
	}
	// separate all network requests to goroutines
	client, err := GetClient(ctx, config, systemKeyring)
	if err != nil {
		slog.Debug("can't load client for client id: "+acc.ClientID, slog.Any("error", err))
		return
	}

	srv, err := libGmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		slog.Debug("Unable to retrieve Gmail client", slog.Any("error", err))
	}

	user := "me"
	var count int64
	ll, err := srv.Users.Labels.Get(user, "UNREAD").Do()
	if err != nil {
		slog.Debug("label get ", slog.Any("error", err))
	} else {
		count = ll.MessagesUnread
	}
	channel <- fmt.Sprintf("%s:%d", acc.Short, count)
}

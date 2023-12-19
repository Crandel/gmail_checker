package mails

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Crandel/gmail/internal/accounts"
	"github.com/Crandel/gmail/internal/extractors"
)

const url = "https://mail.google.com/mail/feed/atom"

// GetMailCount - new goroutine for checking emails
func GetMailCount(channel chan<- string, acc accounts.Account) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Debug("error creating new request", slog.Any("error", err))
	}
	req.Header.Set("Authorization", "Basic %s")
	resp, err := client.Do(req)
	if err != nil {
		slog.Debug("error during request", slog.Any("error", err))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Debug("error during response parsing", slog.Any("error", err))
	}
	count := extractors.ExtractCount(string(body))
	channel <- fmt.Sprintf("%[1]v:%[2]v ", acc.Alias, count)
}

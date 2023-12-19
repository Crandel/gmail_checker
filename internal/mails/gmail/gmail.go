package gmail

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Crandel/gmail/internal/accounts"
	"github.com/Crandel/gmail/internal/extractors"
	"golang.org/x/oauth2"
)

type GMailParams struct {
	accounts.Account
	oauth2.Config
}

const (
	url  = "https://mail.google.com/mail/feed/atom"
	Type = "gmail"
)

// GetMailCount - new goroutine for checking emails
func GetGMailCount(channel chan<- string, mailParams GMailParams) {
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
	channel <- fmt.Sprintf("alias:%v ", count)
}

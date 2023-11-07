package mails

import (
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Crandel/gmail/internal/accounts"
	"github.com/Crandel/gmail/internal/extractors"
)

// GetMailCount - new goroutine for checking emails
func GetMailCount(channel chan<- string, acc accounts.Account) {
	url := "https://mail.google.com/mail/feed/atom"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Debug("error creating new request", slog.Any("error", err))
	}
	tokenStr := fmt.Sprintf("%s:%s", acc.Email, acc.Password)
	token := base64.StdEncoding.EncodeToString([]byte(tokenStr))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
	resp, err := client.Do(req)
	body, err := io.ReadAll(resp.Body)
	count := extractors.ExtractCount(string(body))
	channel <- fmt.Sprintf("%[1]v:%[2]v ", acc.Short, count)
}

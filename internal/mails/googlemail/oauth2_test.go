package googlemail

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

const (
	validKey   = "key"
	missingKey = "missing key"
	errMsg     = "%s error is: %v, wantErr: %v"
)

// This is an interface from oauth2.go file
//
//	type keyringHandler interface {
//	    GetEntry(key string) (string, error)
//	    SetEntry(key string, data string) error
//	}
type mockKeyring struct {
	storage map[string]string
}

func (m mockKeyring) GetEntry(key string) (string, error) {
	v, ok := m.storage[key]
	if !ok {
		return "", errors.New("No value for key: " + key)
	}
	return v, nil
}

func (m mockKeyring) SetEntry(key string, data string) error {
	_, ok := m.storage[key]
	if ok {
		return errors.New("There are already value for key: " + key)
	}

	m.storage[key] = data
	return nil
}

func (m mockKeyring) tokenString(token *oauth2.Token) string {
	tokenByte, _ := json.Marshal(token)
	return string(tokenByte)
}

func TestTokenFromKeyring(t *testing.T) { //nolint: paralleltest // concurrent map writes.
	validToken := &oauth2.Token{AccessToken: "your_token", Expiry: time.Now().AddDate(1, 0, 0)}
	keyringH := mockKeyring{
		storage: map[string]string{},
	}
	validTokenString := keyringH.tokenString(validToken)
	_ = keyringH.SetEntry(validKey, validTokenString)
	tests := []struct {
		name       string
		key        string
		token      *oauth2.Token
		keyringH   keyringHandler
		errMessage string
	}{
		{
			name:       "Successful retrieval of token",
			key:        validKey,
			token:      validToken,
			keyringH:   keyringH,
			errMessage: "",
		},
		{
			name:       "Retrieve missing key",
			key:        missingKey,
			keyringH:   keyringH,
			errMessage: "No value for key: missing key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tokenFromKeyring(tt.key, tt.keyringH)
			if err != nil && tt.errMessage == "" {
				t.Errorf(errMsg, "tokenFromKeyring()", err, tt.errMessage)
				return
			}
			if tt.errMessage != "" {
				if tt.errMessage != err.Error() {
					t.Errorf("error is empty but expect an error: %s", tt.errMessage)
					return
				}
				return
			}

			if got.AccessToken != tt.token.AccessToken {
				t.Errorf("got token: '%s', expected: '%s'", got.AccessToken, tt.token.AccessToken)
			}
		})
	}
}

func Test_saveToken(t *testing.T) {
	t.Parallel()
	keyringH := mockKeyring{
		storage: map[string]string{},
	}
	type args struct {
		key     string
		token   *oauth2.Token
		keyring keyringHandler
	}
	tests := []struct {
		name   string
		args   args
		errMsg string
	}{
		{
			name: "Successful save of token",
			args: args{
				key: validKey,
				token: &oauth2.Token{
					AccessToken:  "access_token",
					RefreshToken: "refresh_token",
				},
				keyring: keyringH,
			},
			errMsg: "",
		},
		{
			name: "Save the same key",
			args: args{
				key: validKey,
				token: &oauth2.Token{
					AccessToken:  "access_token",
					RefreshToken: "refresh_token",
				},
				keyring: keyringH,
			},
			errMsg: "There are already value for key: key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := saveToken(tt.args.key, tt.args.token, tt.args.keyring)
			if err != nil {
				if tt.errMsg == "" {
					t.Errorf("saveToken error: %v, expect no error", err)
				}
				if tt.errMsg != err.Error() {
					t.Errorf("saveToken() error = %v, wantErr %v", err, tt.errMsg)
				}
				return
			}
			if tt.errMsg != "" {
				t.Errorf("saveToken() error is nil, wantErr\n%v", tt.errMsg)
				return
			}
			token, err := tokenFromKeyring(tt.args.key, tt.args.keyring)
			if err != nil {
				if tt.errMsg == "" {
					t.Errorf("saveToken error: %v, expect no error", err)
					return
				}
				if tt.errMsg != err.Error() {
					t.Errorf("saveToken() error:\n%v, wantErr\n%v", err, tt.errMsg)
					return
				}
			}
			if token.AccessToken != tt.args.token.AccessToken {
				t.Errorf("token is invalid. \nWant:\n%v,\ngot:\n%v", tt.args.token, token)
			}
		})
	}
}

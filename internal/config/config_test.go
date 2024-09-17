package config

import (
	"io"
	"reflect"
	"testing"

	"github.com/Crandel/gmail/internal/accounts"
)

//	func Test_addNewUser(t *testing.T) {
//		type args struct {
//			reader *bufio.Reader
//		}
//		tests := []struct {
//			name    string
//			args    args
//			want    accounts.Account
//			wantErr bool
//		}{
//			// TODO: Add test cases.
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				got, err := addNewUser(tt.args.reader)
//				if (err != nil) != tt.wantErr {
//					t.Errorf("addNewUser() error = %v, wantErr %v", err, tt.wantErr)
//					return
//				}
//				if !reflect.DeepEqual(got, tt.want) {
//					t.Errorf("addNewUser() = %v, want %v", got, tt.want)
//				}
//			})
//		}
//	}
//

// // MockReader is a struct that implements the bufio.Reader interface for testing.
// type MockReader struct {
// 	input []string
// }

// func (mr *MockReader) ReadString(delim byte) (string, error) {
// 	if len(mr.input) > 0 {
// 		s := mr.input[0]
// 		mr.input = mr.input[1:]
// 		return s + "\n", nil
// 	}
// 	return "", errors.New("no more input")
// }

// func TestAddNewUser(t *testing.T) {
// 	// Define your test data: alias, type, email and clientId in that order.
// 	testData := []struct {
// 		alias    string
// 		mailType string
// 		email    string
// 		clientID string
// 	}{
// 		{"alias1", "gmail", "example@email.com", "clientId1"},
// 		{"alias2", "notGmail", "anotherExample@email.com", "clientId2"},
// 	}

// 	for _, td := range testData {
// 		// Create a mock reader that returns the test data when ReadString is called.
// 		mockReader := &MockReader{input: []string{td.alias, td.mailType, td.email, td.clientID}}

// 		account, err := addNewUser(mockReader)
// 		if err != nil {
// 			t.Errorf("addNewUser returned error %q when it should not have", err)
// 		}

// 		// Assert that the account has the correct values.
// 		if account.Short != td.alias || account.Email != td.email || account.ClientID != td.clientID {
// 			t.Errorf("addNewUser returned incorrect account: got %+v, want Short=%s, Email=%s, ClientId =%s", account, td.alias, td.email, td.clientID)
// 		}
// 	}
// }

// MockReader is a struct that implements the io.Reader interface for testing.
type MockReader struct {
	input string
}

func (mr MockReader) Read(p []byte) (n int, err error) {
	copy(p, []byte(mr.input))
	return len(mr.input), nil
}

func TestGetAccounts(t *testing.T) {
	type testCase struct {
		name     string
		input    io.Reader
		expected accounts.ListAccounts
		err      error
	}

	var jsonStr = []byte(`[{"Short":"alias1", "MailType":"gmail", "Email":"email@example.com", "ClientID":"clientId1"}, {"Short":"alias2", "MailType":"gmail", "Email":"another@example.com", "ClientID":"clientId2"}]`)
	var accountsList = []accounts.Account{
		{
			Short:    "alias1",
			MailType: accounts.Gmail,
			Email:    "email@example.com",
			ClientID: "clientId1",
		},
		{
			Short:    "alias2",
			MailType: accounts.Gmail,
			Email:    "another@example.com",
			ClientID: "clientId2",
		},
	}

	tests := []testCase{{
		name:     "GetAccounts test",
		input:    MockReader{string(jsonStr)},
		expected: accountsList,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotAccounts, err := GetAccounts(tc.input)
			if err != nil {
				t.Errorf("GetAccounts() error = %v, wantErr %v", err, tc.err)
				return
			}
			if !reflect.DeepEqual(gotAccounts, tc.expected) {
				t.Errorf("GetAccounts() got =%v, want =%v", gotAccounts, tc.expected)
			}
		})
	}
}

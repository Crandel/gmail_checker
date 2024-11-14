package accounts

// MailType is an alias for string.
type MailType string

// Gmail is MailType for google mail.
const Gmail MailType = "gmail"

// ListAccounts - list of accounts from config file
type ListAccounts []Account

// Account type - description of account
type Account struct {
	ClientID string   `json:"client_id"`
	Email    string   `json:"email"`
	MailType MailType `json:"mail_type"`
	Short    string   `json:"short"`
}

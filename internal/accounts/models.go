package accounts

// ListAccounts - list of accounts from config file
type ListAccounts []Account

// Account type - description of account
type Account struct {
	MailType string `json:"mail_type"`
	Email    string `json:"email"`
	ClientID string `json:"client_id"`
}

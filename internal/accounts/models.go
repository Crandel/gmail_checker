package accounts

// ListAccounts - list of accounts from config file
type ListAccounts []Account

// Account type - description of account
type Account struct {
	MailType string `json:"mail_type"`
	Account  string `json:"account"`
	Short    string `json:"short_conky"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

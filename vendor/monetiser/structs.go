package monetiser

type MonetiserRequest struct {
	ID            string  `json:"id"`
	Currency      string  `json:"currency"`
	PrivateURL    string  `json:"private_url"`
	PrivateText   string  `json:"private_text"`
	PublicTitle   string  `json:"public_title"`
	Price         float64 `json:"price"`
	WalletAddress string  `json:"wallet_address"`
	DateExpiry    string  `json:"date_expiry"`
	Address       string  `json:"address"`
	Mode          string  `json:"mode"`
	Email         string  `json:"email"`
}

type ListMontiserRequest struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

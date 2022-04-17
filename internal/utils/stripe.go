package utils

type StripeConfig struct {
	SecretKey string
	PublicKey string
}

type StripePayload struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

package utils

import (
	"net/http"
	"strconv"
)

type TransactionData struct {
	FirstName       string
	LastName        string
	Email           string
	PaymentIntentID string
	PaymentMethodID string
	PaymentAmount   int
	PaymentCurrency string
	LastFour        string
	ExpiryMonth     int
	ExpiryYear      int
	BankReturnCode  string
}

// GetTransactionData get txn data from POST and strip
func GetTransactionData(r *http.Request, sc StripeConfig) (TransactionData, error) {
	var txnData TransactionData

	err := r.ParseForm()
	if err != nil {
		return txnData, err
	}

	lastName := r.Form.Get("last_name")
	firstName := r.Form.Get("first_name")
	email := r.Form.Get("cardholder_email")
	paymentCurrency := r.Form.Get("payment_currency")
	paymentIntentID := r.Form.Get("payment_intent")
	paymentMethodID := r.Form.Get("payment_method")
	paymentAmount, _ := strconv.Atoi(r.Form.Get("payment_amount"))

	card := Card{
		Key:    sc.PublicKey,
		Secret: sc.SecretKey,
	}

	pi, err := card.RetrievePaymentIntent(paymentIntentID)
	if err != nil {
		return txnData, err
	}

	pm, err := card.GetPaymentMethod(paymentMethodID)
	if err != nil {
		return txnData, err
	}

	txnData = TransactionData{
		FirstName:       firstName,
		LastName:        lastName,
		Email:           email,
		PaymentAmount:   paymentAmount,
		PaymentCurrency: paymentCurrency,
		LastFour:        pm.Card.Last4,
		BankReturnCode:  pi.Charges.Data[0].ID,
		ExpiryMonth:     int(pm.Card.ExpMonth),
		ExpiryYear:      int(pm.Card.ExpYear),
		PaymentIntentID: paymentIntentID,
		PaymentMethodID: paymentMethodID,
	}

	return txnData, nil
}

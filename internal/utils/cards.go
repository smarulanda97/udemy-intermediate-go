package utils

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/paymentmethod"
)

type Card struct {
	Secret   string
	Key      string
	Currency string
}

type Transaction struct {
	StatusID       int
	Amount         int
	Currency       string
	LastFour       string
	BankReturnCode string
}

func (c *Card) Charge(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	return c.createPaymentIntent(currency, amount)
}

func (c *Card) createPaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	stripe.Key = c.Secret

	// Payment intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount)),
		Currency: stripe.String(currency),
	}

	// params.AddMetadata("key", "value")
	paymentIntent, err := paymentintent.New(params)
	if err == nil {
		return paymentIntent, "", nil
	}

	var message string
	if stripeErr, ok := err.(*stripe.Error); ok {
		message = cardErrorMessage(stripeErr.Code)
	}

	return nil, message, err
}

// GetPaymentMethod gets the payment method by payment intent id
func (c *Card) GetPaymentMethod(s string) (*stripe.PaymentMethod, error) {
	stripe.Key = c.Secret

	pm, err := paymentmethod.Get(s, nil)
	if err != nil {
		return nil, err
	}

	return pm, nil
}

// RetrievePaymentIntent gets an existing payment intent by id
func (c *Card) RetrievePaymentIntent(id string) (*stripe.PaymentIntent, error) {
	stripe.Key = c.Secret

	pi, err := paymentintent.Get(id, nil)
	if err != nil {
		return nil, err
	}

	return pi, nil
}

func cardErrorMessage(errorCode stripe.ErrorCode) string {
	var message string

	switch errorCode {
	case stripe.ErrorCodeExpiredCard:
		message = "Your card is expired"
	case stripe.ErrorCodeIncorrectCVC:
		message = "Incorrect CVC code"
	case stripe.ErrorCodeAmountTooLarge:
		message = "The amount is too large to charge to your card"
	case stripe.ErrorCodeAmountTooSmall:
		message = "The amount is too small to charge to your card"
	default:
		message = "Your card was declined"
	}

	return message
}

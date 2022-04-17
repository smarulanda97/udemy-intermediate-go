package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/smarulanda97/app-stripe/internal/utils"
)

func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	payload, err := getPayload(w, r)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	amount, err := strconv.Atoi(payload.Amount)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	card := utils.Card{
		Currency: payload.Currency,
		Key:      app.kernel.Stripe.PublicKey,
		Secret:   app.kernel.Stripe.SecretKey,
	}

	okay := true

	paymentIntent, message, err := card.Charge(payload.Currency, amount)
	if err != nil {
		okay = false
	}

	if okay {
		out, err := json.MarshalIndent(paymentIntent, "", "   ")
		if err != nil {
			app.errorLog.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	} else {
		response := utils.JsonResponse{
			Ok:      true,
			Message: message,
			Content: "",
		}

		out, err := json.MarshalIndent(response, "", "   ")
		if err != nil {
			app.errorLog.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}

}

func getPayload(w http.ResponseWriter, r *http.Request) (utils.StripePayload, error) {
	var payload utils.StripePayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return payload, err
	}

	return payload, nil
}

func (app *application) GetProductById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	productId, _ := strconv.Atoi(id)

	product, err := app.DB.GetProduct(productId)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	out, err := json.MarshalIndent(product, "", "   ")
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

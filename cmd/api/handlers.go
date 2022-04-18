package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/smarulanda97/app-stripe/internal/models"
	"github.com/smarulanda97/app-stripe/internal/utils"
)

func (app *application) GetPaymentIntentController(w http.ResponseWriter, r *http.Request) {
	var okay = true
	var response any

	var payload utils.StripePayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	amount, _ := strconv.Atoi(payload.Amount)

	card := utils.Card{
		Currency: payload.Currency,
		Key:      app.kernel.Stripe.PublicKey,
		Secret:   app.kernel.Stripe.SecretKey,
	}

	pi, message, err := card.Charge(payload.Currency, amount)
	if err != nil {
		okay = false
	}

	if !okay {
		response = utils.JsonResponse{
			Ok:      true,
			Message: message,
			Content: "",
		}
	} else {
		response = pi
	}

	out, err := json.MarshalIndent(response, "", "   ")
	if err != nil {
		app.errorLog.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *application) GetProductController(w http.ResponseWriter, r *http.Request) {
	productID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	product, err := app.DB.GetProduct(productID)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	out, _ := json.MarshalIndent(product, "", "   ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (app *application) CreateAuthTokenController(w http.ResponseWriter, r *http.Request) {
	var userInput utils.AuthInput

	err := utils.ReadJson(w, r, &userInput)
	if err != nil {
		utils.BadRequest(w, r, err)
		return
	}

	user, err := app.DB.GetUserByEmail(userInput.Email)
	if err != nil {
		utils.InvalidCredentials(w, r)
		return
	}

	isValidPassword, err := utils.PasswordMatches(user.Password, userInput.Password)
	if err != nil || !isValidPassword {
		utils.InvalidCredentials(w, r)
		return
	}

	token, err := models.GenerateToken(user.ID, 24*time.Hour, models.ScopeAuthentication)
	if err != nil {
		utils.BadRequest(w, r, err)
	}

	// save to database
	err = app.DB.InsertToken(token, user)
	if err != nil {
		utils.BadRequest(w, r, err)
		return
	}

	_ = utils.WriteJson(w, r, http.StatusOK, utils.Payload{
		Error:   false,
		Token:   token,
		Message: "Success!",
	})
}

func (app *application) CheckAuthTokenController(w http.ResponseWriter, r *http.Request) {
	user, err := utils.AuthenticateToken(r, app.DB)
	if err != nil {
		utils.InvalidCredentials(w, r)
		return
	}

	payload := utils.Payload{
		Error:   false,
		Message: fmt.Sprintf("authenticated user %s", user.Email),
	}

	utils.WriteJson(w, r, 200, payload)
}

func (app *application) ApiTerminalController(w http.ResponseWriter, r *http.Request) {

}

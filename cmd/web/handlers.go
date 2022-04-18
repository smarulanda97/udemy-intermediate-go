package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/smarulanda97/app-stripe/internal/models"
	"github.com/smarulanda97/app-stripe/internal/utils"
)

func (app *application) HomeController(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "home", &utils.TemplateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) TerminalController(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "terminal", &utils.TemplateData{}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) PaymentController(w http.ResponseWriter, r *http.Request) {
	txnData, err := utils.GetTransactionData(r, app.kernel.Stripe)
	if err != nil {
		app.errorLog.Println(err)
	}

	if err := r.ParseForm(); err != nil {
		app.errorLog.Println(err)
	}
	productID, _ := strconv.Atoi(r.Form.Get("product_id"))

	// Create a new customer
	customerID, err := app.DB.SaveCustomer(models.Customer{
		FirstName: txnData.FirstName,
		LastName:  txnData.LastName,
		Email:     txnData.Email,
	})
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// Create transaction
	txnID, err := app.DB.SaveTransaction(models.Transaction{
		TransactionStatusID: 2,
		Amount:              txnData.PaymentAmount,
		Currency:            txnData.PaymentCurrency,
		LastFour:            txnData.LastFour,
		ExpiryMonth:         txnData.ExpiryMonth,
		ExpiryYear:          txnData.ExpiryYear,
		BankReturnCode:      txnData.BankReturnCode,
		PaymentIntent:       txnData.PaymentIntentID,
		PaymentMethod:       txnData.PaymentMethodID,
	})
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// Create a new order
	_, err = app.DB.SaveOrder(models.Order{
		CustomerID:    customerID,
		ProductID:     productID,
		TransactionID: txnID,
		StatusID:      1,
		Quantity:      1,
		Amount:        txnData.PaymentAmount,
	})
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// Write this data to session, and the redirect the user to a new page
	app.Session.Put(r.Context(), "receipt", txnData)
	http.Redirect(w, r, "/cart/receipt", http.StatusSeeOther)
}

func (app *application) ReceiptController(w http.ResponseWriter, r *http.Request) {
	txn := app.Session.Get(r.Context(), "receipt").(utils.TransactionData)
	app.Session.Remove(r.Context(), "receipt")

	data := make(map[string]interface{})
	data["txn"] = txn

	templateData := &utils.TemplateData{
		Data: data,
	}

	if err := app.renderTemplate(w, r, "receipt", templateData); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) ProductController(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	productID, _ := strconv.Atoi(id)

	product, err := app.DB.GetProduct(productID)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	data := make(map[string]interface{})
	data["product"] = product

	templateData := &utils.TemplateData{
		Data: data,
	}

	if err := app.renderTemplate(w, r, "product", templateData, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) BronzePlanController(w http.ResponseWriter, r *http.Request) {
	product, err := app.DB.GetProduct(2)
	if err != nil {
		app.errorLog.Println(product)
		return
	}

	data := make(map[string]interface{})
	data["product"] = product

	templateData := &utils.TemplateData{
		Data: data,
	}

	if err := app.renderTemplate(w, r, "bronze-plan", templateData); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) LoginController(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "login", &utils.TemplateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

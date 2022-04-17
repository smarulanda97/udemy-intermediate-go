package main

import (
	"encoding/gob"

	"github.com/alexedwards/scs/v2"
	"github.com/smarulanda97/app-stripe/internal/kernel"
	"github.com/smarulanda97/app-stripe/internal/models"
	"github.com/smarulanda97/app-stripe/internal/utils"

	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

const version = "1.0.0"

var session *scs.SessionManager

type application struct {
	infoLog     *log.Logger
	errorLog    *log.Logger
	kernel      kernel.Kernel
	renderCache map[string]*template.Template
	Session     *scs.SessionManager
	DB          models.DBModels
}

func (app *application) serve() error {
	server := &http.Server{
		Handler:           app.routes(),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		Addr:              fmt.Sprintf(":%d", app.kernel.Port),
	}
	app.infoLog.Printf("Starting HTTP server in %s mode on port %d\r\n", app.kernel.Environment, app.kernel.Port)

	return server.ListenAndServe()
}

func main() {
	gob.Register(utils.TransactionData{})

	k := kernel.Kernel{}
	k.Start(version)

	infoLog, errorLog := k.CreateLoggers()

	connDB, err := utils.OpenDBConnection(k.Database.Dsn)
	if err == nil {
		defer connDB.Close()
	} else {
		errorLog.Fatal(err)
	}

	// setup session
	session = scs.New()
	session.Lifetime = 24 * time.Hour

	app := &application{
		kernel:      k,
		infoLog:     infoLog,
		Session:     session,
		errorLog:    errorLog,
		DB:          models.DBModels{DB: connDB},
		renderCache: make(map[string]*template.Template),
	}

	if err := app.serve(); err != nil {
		app.errorLog.Println(err)
		app.errorLog.Fatal(err)
	}
}

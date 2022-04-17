package main

import (
	"github.com/smarulanda97/app-stripe/internal/kernel"
	"github.com/smarulanda97/app-stripe/internal/models"
	"github.com/smarulanda97/app-stripe/internal/utils"

	"fmt"
	"log"
	"net/http"
	"time"
)

const version = "1.0.0"

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	kernel   kernel.Kernel
	DB       models.DBModels
}

func (app *application) serve() error {
	server := &http.Server{
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		Addr:              fmt.Sprintf(":%d", app.kernel.Port),
	}

	app.infoLog.Printf("Starting Backend server in %s mode on port %d\r\n", app.kernel.Environment, app.kernel.Port)

	return server.ListenAndServe()
}

func main() {
	k := kernel.Kernel{}
	k.Start(version)

	infoLog, errorLog := k.CreateLoggers()

	conn, err := utils.OpenDBConnection(k.Database.Dsn)
	if err == nil {
		defer conn.Close()
	} else {
		errorLog.Fatal(err)
	}

	app := &application{
		kernel:   k,
		infoLog:  infoLog,
		errorLog: errorLog,
		DB:       models.DBModels{DB: conn},
	}

	if err := app.serve(); err != nil {
		app.errorLog.Println(err)
		app.errorLog.Fatal(err)
	}
}

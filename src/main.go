package main

import (
	"blockchain-service/src/database"
	"blockchain-service/src/handlers"
	"blockchain-service/src/middlewares"
	"blockchain-service/src/models"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// connect to the database
	DB, GORM := database.InitDBS()
	defer DB.Close()
	defer GORM.Close()

	// Auto migrate the tables
	GORM.AutoMigrate(&models.Wallet{}, &models.Transaction{})

	// Create a new router
	r := mux.NewRouter()
	r.Use(middlewares.CorsMiddleware)
	r.Use(middlewares.RecoverMiddleware)
	// Register the handlers
	r.HandleFunc("/create_wallet", handlers.CreateWalletHandler(DB)).Methods("POST")
	r.HandleFunc("/process_transaction", handlers.ProcessTransactionHandler(DB)).Methods("POST")
	r.HandleFunc("/transfer_tips", handlers.TransferTipHandler(DB)).Methods("POST")

	service := &http.Server{
		Addr:    ":80",
		Handler: r,
	}

	// run server in a goroutine
	go func() {
		// start the HTTP server
		log.Println("starting HTTP server on 80 port")
		if err := service.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start HTTP server", err)
		}
	}()

	// handle Ctrl+C signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	log.Println("program is exiting")

	ctx, ctxc := context.WithTimeout(context.Background(), time.Second*10)
	defer ctxc()

	if err := service.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

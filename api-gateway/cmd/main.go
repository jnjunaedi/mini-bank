package main

import (
	"log"
	"mini-bank/api-gateway/internal/client"
	"mini-bank/api-gateway/internal/delivery/http/handler"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	accountClient := client.InitAccountClient("localhost:50051")
	accountHandler := handler.NewAccountHandler(accountClient)

	transactionClient := client.InitTransactionClient("localhost:50052")
	transactionHandler := handler.NewTransactionHandler(transactionClient)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/nasabah/register", accountHandler.RegisterHandler)
		r.Get("/nasabah/{no_rek}", accountHandler.CheckProfileHandler)
		r.Post("/transaction/cash-deposit", transactionHandler.CashDepositHandler)
		r.Post("/transaction/cash-withdrawal", transactionHandler.CashWithdrawalHandler)
		r.Post("/transaction/transfer", transactionHandler.TransaferHandler)
		r.Post("/transaction/batch-gaji", transactionHandler.BatchDepositGajiHandler)
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Gagal menyalakan API Gateway: %v", err)
	}
}

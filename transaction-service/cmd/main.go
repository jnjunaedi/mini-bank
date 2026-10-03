package main

import (
	"log"
	transactionpb "mini-bank/pb/transaction"
	"mini-bank/transaction-service/config"
	"mini-bank/transaction-service/internal/delivery/grpc/handler"
	"mini-bank/transaction-service/internal/infrastructure/client"
	"mini-bank/transaction-service/internal/repository"
	"mini-bank/transaction-service/internal/usecase"
	"mini-bank/transaction-service/worker"
	"net"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Gagal memuat file .env : %v\n", err)
	}

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Gagal terhubung ke MySQL : %v\n", err)
	}
	defer db.Close()

	accountClient := client.InitAccountClient("localhost:50051")
	accountBridge := client.NewAccountClient(accountClient)

	txRepo := repository.NewTransactionRepository(db)
	txUc := usecase.NewTransactionUsacase(txRepo, accountBridge)
	txHandler := handler.NewTransactionGRPCHandler(txUc)

	worker.RegisterCronJobs(db, time.Now(), 6)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Gagal membuka port 50052: %v", err)
	}

	server := grpc.NewServer()

	transactionpb.RegisterTransactionServiceServer(server, txHandler)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Gagal menjalankan server gRPC: %v", err)
	}
}

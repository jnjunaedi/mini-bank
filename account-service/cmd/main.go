package main

import (
	"log"
	"mini-bank/account-service/config"
	"mini-bank/account-service/internal/delivery/grpc/handler"
	"mini-bank/account-service/internal/repository"
	"mini-bank/account-service/internal/usecase"
	accountpb "mini-bank/pb/account"
	"net"

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

	accountRepo := repository.NewAccountRepository(db)
	accountUc := usecase.NewAccountUsecase((accountRepo))
	accountGRPCHandler := handler.NewAccountGRPCHandler(accountUc)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Gagal membuka port 50051: %v", err)
	}

	server := grpc.NewServer()

	accountpb.RegisterAccountServiceServer(server, accountGRPCHandler)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Gagal menjalankan server gRPC: %v", err)
	}
}

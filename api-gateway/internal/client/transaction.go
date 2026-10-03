package client

import (
	"log"
	transactionpb "mini-bank/pb/transaction"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitTransactionClient(addressServer string) transactionpb.TransactionServiceClient {
	conn, err := grpc.NewClient(addressServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("API Gateway gagal tersambung ke Transaction Service gRPC: %v", err)
	}

	return transactionpb.NewTransactionServiceClient(conn)

}

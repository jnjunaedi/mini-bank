package client

import (
	"log"
	accountpb "mini-bank/pb/account"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitAccountClient(addressServer string) accountpb.AccountServiceClient {
	conn, err := grpc.NewClient(addressServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("API Gateway gagal tersambung ke Account Service gRPC: %v", err)
	}

	return accountpb.NewAccountServiceClient(conn)

}

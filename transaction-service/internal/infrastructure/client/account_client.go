package client

import (
	"context"
	accountpb "mini-bank/pb/account"
	"mini-bank/pkg/utils/errgrpc"
	"mini-bank/transaction-service/internal/domain"

	"google.golang.org/grpc/status"
)

type accountClient struct {
	grpcClient accountpb.AccountServiceClient
}

func NewAccountClient(grpcClient accountpb.AccountServiceClient) domain.AccountClientBridge {
	return &accountClient{grpcClient: grpcClient}
}

func (c *accountClient) VerifyPIN(ctx context.Context, noRek string, pin string) (domain.VerifyPINResult, error) {
	res, err := c.grpcClient.VerifyPIN(ctx, &accountpb.VerifyPINRequest{NoRek: noRek, Pin: pin})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			return domain.VerifyPINResult{Message: st.Message()}, errgrpc.ReadGRPCError(err)
		}
		return domain.VerifyPINResult{Message: "gagal terhubung ke layanan otentikasi"}, errgrpc.ReadGRPCError(err)
	}

	return domain.VerifyPINResult{IsValid: res.GetIsValid(), NasabahID: res.GetNasabahId(), Message: res.GetMessage()}, nil
}

func (c *accountClient) FindByNoRek(ctx context.Context, noRek string) (domain.AccountProfileResult, error) {
	res, err := c.grpcClient.GetAccountDetail(ctx, &accountpb.GetAccountDetailRequest{NoRek: noRek})
	if err != nil {
		return domain.AccountProfileResult{}, errgrpc.ReadGRPCError(err)
	}

	return domain.AccountProfileResult{NoRek: res.GetNoRek(), Nama: res.GetNama(), Status: res.GetStatus()}, nil
}

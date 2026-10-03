package handler

import (
	"context"
	"mini-bank/account-service/internal/domain"
	accountpb "mini-bank/pb/account"
	"mini-bank/pkg/utils/errgrpc"
)

type accountGRPCHandler struct {
	accountpb.UnimplementedAccountServiceServer
	uc domain.AccountUsecase
}

func NewAccountGRPCHandler(uc domain.AccountUsecase) accountpb.AccountServiceServer {
	return &accountGRPCHandler{uc: uc}
}

func (h *accountGRPCHandler) RegisterNasabah(ctx context.Context, req *accountpb.RegisterNasabahRequest) (*accountpb.RegisterNasabahResponse, error) {
	regReq := domain.RegisterRequest{
		NoIdentitas:   req.GetNomorIdentitas(),
		Npwp:          req.GetNpwp(),
		TipeNasabah:   req.GetTipeNasabah(),
		Nama:          req.GetNama(),
		NoHP:          req.GetNoHp(),
		Email:         req.GetEmail(),
		Pekerjaan:     req.GetPekerjaan(),
		AlamatKantor:  req.GetAlamatKantor(),
		AlamatNasabah: req.GetAlamatNasabah(),
		Pin:           req.GetPin(),
		SetoranAwal:   req.GetSetoranAwal(),
	}

	res, err := h.uc.Register(ctx, regReq)
	if err != nil {
		return nil, errgrpc.ReadGRPCError(err)
	}

	return &accountpb.RegisterNasabahResponse{Nama: res.Nama, NoRek: res.NoRek, Saldo: res.SaldoAwal}, nil
}

func (h *accountGRPCHandler) GetAccountDetail(ctx context.Context, req *accountpb.GetAccountDetailRequest) (*accountpb.GetAccountDetailResponse, error) {
	res, err := h.uc.CheckProfile(ctx, req.GetNoRek())
	if err != nil {
		return nil, errgrpc.ReadGRPCError(err)
	}

	return &accountpb.GetAccountDetailResponse{NasabahId: res.NasabahID, NoRek: res.NoRek, Nama: res.Nama, Saldo: res.Saldo, Status: res.Status}, nil
}

func (h *accountGRPCHandler) VerifyPIN(ctx context.Context, req *accountpb.VerifyPINRequest) (*accountpb.VerifyPINResponse, error) {
	res, err := h.uc.VerifyPIN(ctx, req.GetNoRek(), req.GetPin())
	if err != nil {
		return &accountpb.VerifyPINResponse{}, errgrpc.ReadGRPCError(err)
	}

	return &accountpb.VerifyPINResponse{IsValid: res.IsValid, NasabahId: res.NasabahID, Message: res.Message}, nil
}

package handler

import (
	"context"
	transactionpb "mini-bank/pb/transaction"
	"mini-bank/pkg/utils/errgrpc"
	"mini-bank/transaction-service/internal/domain"
)

type transactionGRPCHandler struct {
	transactionpb.UnimplementedTransactionServiceServer
	uc domain.TransactionUsecase
}

func NewTransactionGRPCHandler(uc domain.TransactionUsecase) transactionpb.TransactionServiceServer {
	return &transactionGRPCHandler{uc: uc}
}

func (h *transactionGRPCHandler) CashDeposit(ctx context.Context, req *transactionpb.CashDepositRequest) (*transactionpb.CashDepositResponse, error) {
	r := domain.CashTransactionRequest{
		NoRek:      req.GetNoRek(),
		PIN:        req.GetPin(),
		Nominal:    req.GetNominal(),
		Keterangan: req.GetKeterangan(),
	}
	res, err := h.uc.CashDeposit(ctx, r)
	if err != nil {
		return nil, errgrpc.ReadGRPCError(err)
	}

	return &transactionpb.CashDepositResponse{TransaksiId: res.TransaksiID, NoRek: res.NoRek, TotalSaldo: res.TotalSaldo, Message: res.Message}, nil
}

func (h *transactionGRPCHandler) CashWithdrawal(ctx context.Context, req *transactionpb.CashWithdrawalRequest) (*transactionpb.CashWithdrawalResponse, error) {
	r := domain.CashTransactionRequest{
		NoRek:      req.GetNoRek(),
		PIN:        req.GetPin(),
		Nominal:    req.GetNominal(),
		Keterangan: req.GetKeterangan(),
	}
	res, err := h.uc.CashWithdrawal(ctx, r)
	if err != nil {
		return nil, errgrpc.ReadGRPCError(err)
	}

	return &transactionpb.CashWithdrawalResponse{TransaksiId: res.TransaksiID, NoRek: res.NoRek, TotalSaldo: res.TotalSaldo, Message: res.Message}, nil
}

func (h *transactionGRPCHandler) Transafer(ctx context.Context, req *transactionpb.TransferRequest) (*transactionpb.TransferResponse, error) {
	res, err := h.uc.Transfer(ctx, domain.TransferRequest{NoRekPengirim: req.GetNoRek(), PIN: req.GetPin(), NoRekPenerima: req.GetNoRekTujuan(), Nominal: req.GetNominal()})
	if err != nil {
		return nil, errgrpc.ReadGRPCError(err)
	}

	return &transactionpb.TransferResponse{TransaksiId: res.TransaksiID, SisaSaldo: res.SisaSaldo, Message: res.Message}, nil
}

func (h *transactionGRPCHandler) BatchDepositGaji(ctx context.Context, req *transactionpb.BatchDepositGajiRequest) (*transactionpb.BatchDepositGajiResponse, error) {

	gajiItem := make([]domain.GajiItem, 0, len(req.GetDaftarGaji()))
	for _, karyawan := range req.GetDaftarGaji() {
		item := domain.GajiItem{NoRekKaryawan: karyawan.GetNoRekKaryawan(), NominalGaji: karyawan.GetNominalGaji()}
		gajiItem = append(gajiItem, item)
	}

	res, err := h.uc.ProcessBatchGaji(ctx, domain.BatchGajiRequest{NoRekPerusahaan: req.GetNoRekPerusahaan(), NamaPerusahaan: req.GetNamaPerusahaan(), PINPerusahaan: req.GetPinPerusahaan(), DaftarKaryawan: gajiItem})
	if err != nil {
		return nil, errgrpc.ReadGRPCError(err)
	}

	return &transactionpb.BatchDepositGajiResponse{Status: res.Status, TotalDiproses: res.TotalDiproses, TotalDanaKeluar: res.TotalDanaKeluar, Message: res.Message}, nil
}

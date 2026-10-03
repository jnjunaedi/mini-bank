package handler

import (
	"mini-bank/api-gateway/internal/delivery/http/request"
	"mini-bank/api-gateway/pkg/validator"
	transactionpb "mini-bank/pb/transaction"
	"mini-bank/pkg/utils"
	"mini-bank/pkg/utils/errgrpc"
	"net/http"
)

type TransactionHandler struct {
	client transactionpb.TransactionServiceClient
}

func NewTransactionHandler(client transactionpb.TransactionServiceClient) *TransactionHandler {
	return &TransactionHandler{client: client}
}

func (h *TransactionHandler) CashDepositHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req request.CashTransactionRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
		return
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"message": "validasi payload gagal",
			"error":   errs,
		})
		return
	}

	res, err := h.client.CashDeposit(ctx, &transactionpb.CashDepositRequest{NoRek: req.NoRek, Nominal: req.Nominal, Keterangan: req.Keterangan, Pin: req.PIN})
	if err != nil {
		errgrpc.WriteGRPCError(w, err)
		return
	}

	var resC = struct {
		TransaksiID string `json:"transaksi_id"`
		NoRek       string `json:"no_rek"`
		TotalSaldo  int64  `json:"total_saldo"`
		Message     string `json:"message"`
	}{
		TransaksiID: res.GetTransaksiId(),
		NoRek:       res.GetNoRek(),
		TotalSaldo:  res.GetTotalSaldo(),
		Message:     res.GetMessage(),
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": resC,
	})

}

func (h *TransactionHandler) CashWithdrawalHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req request.CashTransactionRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
		return
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"message": "validasi payload gagal",
			"error":   errs,
		})
		return
	}

	res, err := h.client.CashWithdrawal(ctx, &transactionpb.CashWithdrawalRequest{NoRek: req.NoRek, Nominal: req.Nominal, Keterangan: req.Keterangan, Pin: req.PIN})
	if err != nil {
		errgrpc.WriteGRPCError(w, err)
		return
	}

	var resC = struct {
		TransaksiID string `json:"transaksi_id"`
		NoRek       string `json:"no_rek"`
		TotalSaldo  int64  `json:"total_saldo"`
		Message     string `json:"message"`
	}{
		TransaksiID: res.GetTransaksiId(),
		NoRek:       res.GetNoRek(),
		TotalSaldo:  res.GetTotalSaldo(),
		Message:     res.GetMessage(),
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": resC,
	})

}

func (h *TransactionHandler) TransaferHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req request.TransferRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
		return
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"message": "validasi payload gagal",
			"error":   errs,
		})
		return
	}

	res, err := h.client.Transafer(ctx, &transactionpb.TransferRequest{NoRek: req.NoRek, Pin: req.PIN, NoRekTujuan: req.NoRekPenerima, Nominal: req.Nominal})
	if err != nil {
		errgrpc.WriteGRPCError(w, err)
		return
	}

	var resC = struct {
		TransaksiID string `json:"transaksi_id"`
		SisaSaldo   int64  `json:"sisa_saldo"`
		Message     string `json:"message"`
	}{
		TransaksiID: res.GetTransaksiId(),
		SisaSaldo:   res.GetSisaSaldo(),
		Message:     res.GetMessage(),
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": resC,
	})
}

func (h *TransactionHandler) BatchDepositGajiHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req request.BatchDepositGajiRequest
	if err := utils.ReadJSON(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
		return
	}

	if errs := validator.ValidateStruct(req); errs != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"message": "validasi payload gagal",
			"error":   errs,
		})
		return
	}

	txConvert := make([]*transactionpb.GajiItem, 0, len(req.GajiItem))

	for _, gaji := range req.GajiItem {
		item := &transactionpb.GajiItem{
			NoRekKaryawan: gaji.NoRekKaryawan,
			NominalGaji:   gaji.NominalGaji,
		}

		txConvert = append(txConvert, item)
	}

	res, err := h.client.BatchDepositGaji(ctx, &transactionpb.BatchDepositGajiRequest{NoRekPerusahaan: req.NoRekPerusahaan, NamaPerusahaan: req.NamaPerusahaan, PinPerusahaan: req.PINPerusahaan, DaftarGaji: txConvert})
	if err != nil {
		errgrpc.WriteGRPCError(w, err)
		return
	}

	var resC = struct {
		Status          string `json:"status"`
		TotalDiproses   int32  `json:"total_diproses"`
		TotalDanaKeluar int64  `json:"total_dana_keluar"`
		Message         string `json:"message"`
	}{
		Status:          res.GetStatus(),
		TotalDiproses:   res.GetTotalDiproses(),
		TotalDanaKeluar: res.GetTotalDanaKeluar(),
		Message:         res.GetMessage(),
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": resC,
	})
}

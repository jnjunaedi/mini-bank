package handler

import (
	"mini-bank/api-gateway/internal/delivery/http/request"
	"mini-bank/api-gateway/pkg/validator"
	accountpb "mini-bank/pb/account"
	"mini-bank/pkg/utils"
	"mini-bank/pkg/utils/errgrpc"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AccountHandler struct {
	client accountpb.AccountServiceClient
}

func NewAccountHandler(client accountpb.AccountServiceClient) *AccountHandler {
	return &AccountHandler{client: client}
}

func (h *AccountHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var accountReq request.AccountRequest

	if err := utils.ReadJSON(r, &accountReq); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
		return
	}

	if errs := validator.ValidateStruct(accountReq); errs != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"message": "validasi payload gagal",
			"error":   errs,
		})
		return
	}

	nasabah, err := h.client.RegisterNasabah(ctx, &accountpb.RegisterNasabahRequest{NomorIdentitas: accountReq.NomorIdentitas, Nama: accountReq.Nama, NoHp: accountReq.NoHP, Email: accountReq.Email, Pekerjaan: accountReq.Pekerjaan, AlamatKantor: accountReq.AlamatKantor, AlamatNasabah: accountReq.AlamatNasabah, TipeNasabah: accountReq.TipeNasabah, Npwp: accountReq.Npwp, Pin: accountReq.Pin, SetoranAwal: accountReq.SetoranAwal})
	if err != nil {
		errgrpc.WriteGRPCError(w, err)
		return
	}

	var res = struct {
		Nama        string `json:"nama"`
		NoRekening  string `json:"no_rek"`
		SetoranAwal int64  `json:"setoran_awal"`
	}{
		Nama:        nasabah.GetNama(),
		NoRekening:  nasabah.GetNoRek(),
		SetoranAwal: nasabah.GetSaldo(),
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"data": res,
	})

}

func (h *AccountHandler) CheckProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	noRek := chi.URLParam(r, "no_rek")

	resGRPC, err := h.client.GetAccountDetail(ctx, &accountpb.GetAccountDetailRequest{NoRek: noRek})
	if err != nil {
		errgrpc.WriteGRPCError(w, err)
		return
	}

	var res = struct {
		NasabahID string
		NoRek     string
		Nama      string
		Saldo     int64
		Status    string
	}{
		NasabahID: resGRPC.NasabahId,
		NoRek:     resGRPC.NoRek,
		Nama:      resGRPC.Nama,
		Saldo:     resGRPC.Saldo,
		Status:    resGRPC.Status,
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": res,
	})
}

func (h *AccountHandler) VerifyPINHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req request.VerifyPINRequest

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

	res, err := h.client.VerifyPIN(ctx, &accountpb.VerifyPINRequest{NoRek: req.NoRek, Pin: req.Pin})
	if err != nil {
		errgrpc.WriteGRPCError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data":     res.GetNasabahId(),
		"message":  res.GetMessage(),
		"is_valid": res.GetIsValid(),
	})

}

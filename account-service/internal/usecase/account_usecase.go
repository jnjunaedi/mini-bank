package usecase

import (
	"context"
	"mini-bank/account-service/internal/domain"
	"mini-bank/account-service/internal/pkg/utils"
	"mini-bank/pkg/security"
	"mini-bank/pkg/utils/apperror"

	"github.com/google/uuid"
)

type AccountUsecase struct {
	repo domain.AccountRepository
}

func NewAccountUsecase(repo domain.AccountRepository) domain.AccountUsecase {
	return &AccountUsecase{repo: repo}
}

func (u *AccountUsecase) Register(ctx context.Context, req domain.RegisterRequest) (domain.RegisterResponse, error) {

	genUUIDGoogle, err := uuid.NewV7()
	if err != nil {
		return domain.RegisterResponse{}, apperror.New(apperror.ErrInternalServer, "error internal server")
	}
	genUUID := genUUIDGoogle.String()

	genRef, err := uuid.NewV7()
	if err != nil {
		return domain.RegisterResponse{}, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	genIDRw, _ := uuid.NewV7()

	pinHash, err := security.HashPassword(req.Pin)
	if err != nil {
		return domain.RegisterResponse{}, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	noRek := utils.GenerateNomorRekening()

	var nasabah domain.Nasabah
	var rekening domain.Rekening
	var regResp domain.RegisterResponse

	if req.TipeNasabah == "PERUSAHAAN" {
		if req.AlamatKantor == "" {
			return domain.RegisterResponse{}, apperror.New(apperror.ErrBadRequest, "kolom alamat_kantor wajib diisi untuk registrasi perusahaan")
		}
		nasabah = domain.Nasabah{
			ID:            genUUID,
			TipeNasabah:   "PERUSAHAAN",
			NoIdentitas:   req.NoIdentitas,
			Npwp:          &req.Npwp,
			Nama:          req.Nama,
			NoHP:          req.NoHP,
			Email:         req.Email,
			AlamatKantor:  &req.AlamatKantor,
			Pin:           pinHash,
			Pekerjaan:     nil,
			AlamatNasabah: nil,
		}
		rekening = domain.Rekening{NoRek: noRek, NasabahID: genUUID, Status: "AKTIF", Saldo: req.SetoranAwal}

		regResp = domain.RegisterResponse{
			Nama:      req.Nama,
			NoRek:     rekening.NoRek,
			SaldoAwal: req.SetoranAwal,
		}

	} else {
		if req.AlamatNasabah == "" {
			return domain.RegisterResponse{}, apperror.New(apperror.ErrBadRequest, "kolom alamat_nasabah wajib diisi untuk registrasi akun nasabah")
		}
		nasabah = domain.Nasabah{
			ID:            genUUID,
			TipeNasabah:   "PERORANGAN",
			NoIdentitas:   req.NoIdentitas,
			Nama:          req.Nama,
			NoHP:          req.NoHP,
			Email:         req.Email,
			Pekerjaan:     &req.Pekerjaan,
			AlamatNasabah: &req.AlamatNasabah,
			Pin:           pinHash,
			Npwp:          nil,
			AlamatKantor:  nil,
		}
		rekening = domain.Rekening{NoRek: noRek, NasabahID: genUUID, Status: "AKTIF", Saldo: req.SetoranAwal}

		regResp = domain.RegisterResponse{
			Nama:      req.Nama,
			NoRek:     rekening.NoRek,
			SaldoAwal: req.SetoranAwal,
		}
	}

	err = u.repo.CreateWithTransaction(ctx, nasabah, rekening, genRef.String(), genIDRw.String())
	if err != nil {
		return domain.RegisterResponse{}, err
	}

	return regResp, nil
}

func (u *AccountUsecase) CheckProfile(ctx context.Context, noRek string) (domain.AccountDetailResponse, error) {
	nasabah, rekening, err := u.repo.FindByNoRek(ctx, noRek)
	if err != nil {
		return domain.AccountDetailResponse{}, err
	}

	res := domain.AccountDetailResponse{
		NasabahID: nasabah.ID,
		NoRek:     rekening.NoRek,
		Nama:      nasabah.Nama,
		Saldo:     rekening.Saldo,
		Status:    rekening.Status,
	}

	return res, nil
}

func (u *AccountUsecase) VerifyPIN(ctx context.Context, noRek string, pin string) (domain.VerifyPINResponse, error) {
	nasabah, rekening, err := u.repo.FindByNoRek(ctx, noRek)
	if err != nil {
		return domain.VerifyPINResponse{}, err
	}

	if rekening.Status == "DIBLOKIR" {
		return domain.VerifyPINResponse{}, apperror.New(apperror.ErrForbidden, "rekening Anda sedang DIBLOKIR. Hubungi CS")
	}

	if rekening.Status == "PASIF" {
		return domain.VerifyPINResponse{}, apperror.New(apperror.ErrForbidden, "rekening berstatus PASIF. Silakan aktivasi dahulu")
	}

	isMatch := security.ComparePassword(pin, nasabah.Pin)
	if !isMatch {
		nextAttempts := rekening.FailedAttempts + 1
		if nextAttempts >= 3 {
			_ = u.repo.UpdateStatusDanAttempts(ctx, noRek, "DIBLOKIR", nextAttempts)
			return domain.VerifyPINResponse{}, apperror.New(apperror.ErrForbidden, "rekening Anda otomatis DIBLOKIR karena 3x salah memasukkan PIN")
		}
		_ = u.repo.UpdateStatusDanAttempts(ctx, noRek, rekening.Status, nextAttempts)
		return domain.VerifyPINResponse{}, apperror.New(apperror.ErrUnauthorized, "PIN yang Anda masukkan salah")
	}

	if rekening.FailedAttempts > 0 {
		_ = u.repo.UpdateStatusDanAttempts(ctx, noRek, "AKTIF", 0)
	}
	return domain.VerifyPINResponse{IsValid: true, NasabahID: nasabah.ID, Message: "Berhasil login"}, nil
}

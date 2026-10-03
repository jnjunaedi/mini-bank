package usecase

import (
	"context"
	"fmt"
	"mini-bank/pkg/utils/apperror"
	"mini-bank/transaction-service/internal/domain"

	"github.com/google/uuid"
)

type TransactionUsecase struct {
	repo         domain.TransactionRepository
	accounBridge domain.AccountClientBridge
}

func NewTransactionUsacase(repo domain.TransactionRepository, accounBridge domain.AccountClientBridge) domain.TransactionUsecase {
	return &TransactionUsecase{repo: repo, accounBridge: accounBridge}
}

func (u *TransactionUsecase) CashDeposit(ctx context.Context, req domain.CashTransactionRequest) (domain.CashTransactionResponse, error) {
	res, err := u.accounBridge.VerifyPIN(ctx, req.NoRek, req.PIN)
	if err != nil {
		return domain.CashTransactionResponse{}, err
	}

	if !res.IsValid {
		return domain.CashTransactionResponse{}, apperror.New(apperror.ErrForbidden, res.Message)
	}

	genUUIDGoogle, err := uuid.NewV7()
	if err != nil {
		return domain.CashTransactionResponse{}, apperror.New(apperror.ErrInternalServer, "error internal server")
	}
	txIDStr := genUUIDGoogle.String()

	genRef, err := uuid.NewV7()
	if err != nil {
		return domain.CashTransactionResponse{}, apperror.New(apperror.ErrInternalServer, "error internal server")
	}
	referensiIDStr := genRef.String()

	sisaSaldo, err := u.repo.SaveTransaction(ctx, domain.RiwayatTx{ID: txIDStr, ReferensiID: referensiIDStr, NoRekNasabah: req.NoRek, JenisTransaksi: "MASUK", KategoriTx: "SETOR_TUNAI", Nominal: req.Nominal, NoRekTujuan: "", Keterangan: "Setor Tunai Mandiri"})
	if err != nil {
		return domain.CashTransactionResponse{}, err
	}

	resTx := domain.CashTransactionResponse{
		TransaksiID: txIDStr,
		NoRek:       req.NoRek,
		TotalSaldo:  sisaSaldo,
		Message:     "Setor tunai berhasil diproses",
	}

	return resTx, nil

}

func (u *TransactionUsecase) CashWithdrawal(ctx context.Context, req domain.CashTransactionRequest) (domain.CashTransactionResponse, error) {
	res, err := u.accounBridge.VerifyPIN(ctx, req.NoRek, req.PIN)
	if err != nil {
		return domain.CashTransactionResponse{}, err
	}

	if !res.IsValid {
		return domain.CashTransactionResponse{}, apperror.New(apperror.ErrForbidden, res.Message)
	}

	genUUIDGoogle, err := uuid.NewV7()
	if err != nil {
		return domain.CashTransactionResponse{}, apperror.New(apperror.ErrInternalServer, "error internal server")
	}
	txIDStr := genUUIDGoogle.String()

	genRef, err := uuid.NewV7()
	if err != nil {
		return domain.CashTransactionResponse{}, apperror.New(apperror.ErrInternalServer, "error internal server")
	}
	referensiIDStr := genRef.String()

	sisaSaldo, err := u.repo.CashWithdrawal(ctx, domain.RiwayatTx{ID: txIDStr, ReferensiID: referensiIDStr, NoRekNasabah: req.NoRek, JenisTransaksi: "KELUAR", KategoriTx: "TARIK_TUNAI", Nominal: req.Nominal, NoRekTujuan: "", TotalSaldo: req.Nominal, Keterangan: "Tarik Tunai Mandiri"})
	if err != nil {
		return domain.CashTransactionResponse{}, err
	}

	resTx := domain.CashTransactionResponse{
		TransaksiID: txIDStr,
		NoRek:       req.NoRek,
		TotalSaldo:  sisaSaldo,
		Message:     "Tarik tunai berhasil diproses",
	}

	return resTx, nil
}

func (u *TransactionUsecase) Transfer(ctx context.Context, req domain.TransferRequest) (domain.TransferResponse, error) {
	res, err := u.accounBridge.VerifyPIN(ctx, req.NoRekPengirim, req.PIN)
	if err != nil {
		return domain.TransferResponse{}, err
	}

	if !res.IsValid {
		return domain.TransferResponse{}, apperror.New(apperror.ErrForbidden, res.Message)
	}

	pengirim, err := u.accounBridge.FindByNoRek(ctx, req.NoRekPengirim)
	if err != nil {
		return domain.TransferResponse{}, apperror.New(apperror.ErrNotFound, "No rekening tidak valid")
	}

	penerima, err := u.accounBridge.FindByNoRek(ctx, req.NoRekPenerima)
	if err != nil {
		return domain.TransferResponse{}, apperror.New(apperror.ErrNotFound, "Rekening tujuan/penerima tidak ditemukan")
	}

	if penerima.Status == "DIBLOKIR" {
		return domain.TransferResponse{}, apperror.New(apperror.ErrForbidden, "Rekening tujuan sedang diblokir, tidak bisa menerima dana")
	}

	if penerima.Status == "PASIF" {
		return domain.TransferResponse{}, apperror.New(apperror.ErrForbidden, "Status akun nomor rekening tujuan sedang pasif, tidak bisa melakukan transfer. Silakan hubungi layanan pelanggan.")
	}

	keteranganPengirim := "Transfer ke " + penerima.NoRek + " " + penerima.Nama
	keteranganPenerima := "Transfer dari " + pengirim.NoRek + " " + pengirim.Nama

	genRef, err := uuid.NewV7()
	if err != nil {
		return domain.TransferResponse{}, apperror.New(apperror.ErrInternalServer, "error internal server")
	}
	referensiIDStr := genRef.String()

	uuidPengirim, err := uuid.NewV7()
	if err != nil {
		return domain.TransferResponse{}, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	uuidPenerima, err := uuid.NewV7()
	if err != nil {
		return domain.TransferResponse{}, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	saldoPengirim, err := u.repo.Transfer(ctx, domain.RiwayatTx{ID: uuidPengirim.String(), ReferensiID: referensiIDStr, NoRekNasabah: req.NoRekPengirim, JenisTransaksi: "KELUAR", KategoriTx: "TRANSFER", Nominal: req.Nominal, Keterangan: keteranganPengirim}, domain.RiwayatTx{ID: uuidPenerima.String(), ReferensiID: referensiIDStr, NoRekNasabah: req.NoRekPenerima, JenisTransaksi: "MASUK", KategoriTx: "TRANSFER", Nominal: req.Nominal, Keterangan: keteranganPenerima})
	if err != nil {
		return domain.TransferResponse{}, err
	}

	return domain.TransferResponse{TransaksiID: referensiIDStr, SisaSaldo: saldoPengirim, Message: fmt.Sprintf("Transfer sukses ke %s sebesar Rp%d", penerima.Nama, req.Nominal)}, nil
}

func (u *TransactionUsecase) ProcessBatchGaji(ctx context.Context, req domain.BatchGajiRequest) (domain.BatchGajiResponse, error) {
	resVerify, err := u.accounBridge.VerifyPIN(ctx, req.NoRekPerusahaan, req.PINPerusahaan)
	if err != nil {
		return domain.BatchGajiResponse{}, err
	}

	if !resVerify.IsValid {
		return domain.BatchGajiResponse{}, apperror.New(apperror.ErrForbidden, resVerify.Message)
	}

	genRef, err := uuid.NewV7()
	if err != nil {
		return domain.BatchGajiResponse{}, apperror.New(apperror.ErrInternalServer, "error internal server")
	}
	referensiIDStr := genRef.String()

	var danaKeluar int64
	kumupulanKaryawan := make([]domain.RiwayatTx, 0, len(req.DaftarKaryawan))

	for _, karyawan := range req.DaftarKaryawan {

		resKaryawan, err := u.accounBridge.FindByNoRek(ctx, karyawan.NoRekKaryawan)
		if err != nil {
			return domain.BatchGajiResponse{}, err
		}

		if resKaryawan.Status == "DIBLOKIR" {
			return domain.BatchGajiResponse{}, apperror.New(apperror.ErrForbidden, "Rekening tujuan"+karyawan.NoRekKaryawan+" sedang diblokir")
		}

		if resKaryawan.Status == "PASIF" {
			return domain.BatchGajiResponse{}, apperror.New(apperror.ErrForbidden, "Rekening tujuan "+karyawan.NoRekKaryawan+" berstatus pasif")
		}

		uuidRw, err := uuid.NewV7()
		if err != nil {
			return domain.BatchGajiResponse{}, apperror.New(apperror.ErrInternalServer, "error internal server")
		}

		keteranganKaryawan := "Payroll Gaji " + req.NamaPerusahaan

		rwTx := domain.RiwayatTx{
			ID:             uuidRw.String(),
			ReferensiID:    referensiIDStr,
			NoRekNasabah:   karyawan.NoRekKaryawan,
			JenisTransaksi: "MASUK",
			KategoriTx:     "GAJI",
			Nominal:        karyawan.NominalGaji,
			NoRekTujuan:    req.NoRekPerusahaan,
			Keterangan:     keteranganKaryawan,
		}

		kumupulanKaryawan = append(kumupulanKaryawan, rwTx)
		danaKeluar += karyawan.NominalGaji

	}

	jumlahSukses, err := u.repo.SaveBatchGaji(ctx, req.NoRekPerusahaan, danaKeluar, kumupulanKaryawan)
	if err != nil {
		return domain.BatchGajiResponse{}, err
	}

	return domain.BatchGajiResponse{Status: "Success", TotalDiproses: jumlahSukses, TotalDanaKeluar: danaKeluar, Message: fmt.Sprintf("Sukses! Gaji massal telah berhasil ditransfer ke %d karyawan.", jumlahSukses)}, nil

}

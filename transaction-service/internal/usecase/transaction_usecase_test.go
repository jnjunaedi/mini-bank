package usecase_test

import (
	"context"
	"fmt"
	"mini-bank/transaction-service/internal/domain"
	"mini-bank/transaction-service/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTransactionRepository struct {
	mock.Mock
}

func (m *mockTransactionRepository) SaveTransaction(ctx context.Context, riwayatTx domain.RiwayatTx) (int64, error) {
	args := m.Called(ctx, riwayatTx)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockTransactionRepository) CashWithdrawal(ctx context.Context, riwayatTx domain.RiwayatTx) (int64, error) {
	args := m.Called(ctx, riwayatTx)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockTransactionRepository) Transfer(ctx context.Context, riwayatTxPengirim, riwayatTxPenerima domain.RiwayatTx) (int64, error) {
	args := m.Called(ctx, riwayatTxPengirim, riwayatTxPenerima)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockTransactionRepository) SaveBatchGaji(ctx context.Context, noRekPerusahaan string, totalDanaGaji int64, listTx []domain.RiwayatTx) (int32, error) {
	args := m.Called(ctx, noRekPerusahaan, totalDanaGaji, listTx)
	return args.Get(0).(int32), args.Error(1)
}

type mockAccountClientBridge struct {
	mock.Mock
}

func (m *mockAccountClientBridge) FindByNoRek(ctx context.Context, noRek string) (domain.AccountProfileResult, error) {
	args := m.Called(ctx, noRek)
	return args.Get(0).(domain.AccountProfileResult), args.Error(1)
}

func (m *mockAccountClientBridge) VerifyPIN(ctx context.Context, noRek, pin string) (domain.VerifyPINResult, error) {
	args := m.Called(ctx, noRek, pin)
	return args.Get(0).(domain.VerifyPINResult), args.Error(1)
}

func TestCashDeposit(t *testing.T) {
	mockRepo := new(mockTransactionRepository)
	mockAccountBridge := new(mockAccountClientBridge)

	u := usecase.NewTransactionUsacase(mockRepo, mockAccountBridge)

	ctx := context.Background()

	reqCashTx := domain.CashTransactionRequest{
		NoRek:      "32103094930",
		PIN:        "123123",
		Nominal:    120000,
		Keterangan: "Setor Tunai pendaftaran awal",
	}

	reqRw := domain.RiwayatTx{
		NoRekNasabah:   "32103094930",
		KategoriTx:     "MASUK",
		JenisTransaksi: "SETOR_TUNAI",
		Nominal:        reqCashTx.Nominal,
	}

	mockAccountBridge.On("VerifyPIN", ctx, reqCashTx.NoRek, reqCashTx.PIN).Return(domain.VerifyPINResult{
		IsValid:   true,
		NasabahID: "01a0ece9-e8d5-78b3-80d1-042c452ddc4f",
		Message:   "PIN valid",
	}, nil)

	mockRepo.On("SaveTransaction", ctx, mock.MatchedBy(func(rw domain.RiwayatTx) bool {
		return rw.NoRekNasabah == reqRw.NoRekNasabah && rw.KategoriTx == "SETOR_TUNAI" && rw.JenisTransaksi == "MASUK" && rw.Nominal == reqCashTx.Nominal
	})).Return(int64(90000000), nil)

	res, err := u.CashDeposit(ctx, reqCashTx)

	assert.NoError(t, err)
	assert.Equal(t, int64(90000000), res.TotalSaldo)
	assert.Equal(t, "Setor tunai berhasil diproses", res.Message)

	mockRepo.AssertExpectations(t)
	mockAccountBridge.AssertExpectations(t)

}

func TestCashWithdrawal(t *testing.T) {
	mockRepo := new(mockTransactionRepository)
	mockAccountBridge := new(mockAccountClientBridge)

	u := usecase.NewTransactionUsacase(mockRepo, mockAccountBridge)

	ctx := context.Background()

	reqCashTx := domain.CashTransactionRequest{
		NoRek:   "321039837401",
		PIN:     "123123",
		Nominal: 10000000,
	}

	resVerifyPIN := domain.VerifyPINResult{
		IsValid:   true,
		NasabahID: "01a0ece9-e8d5-78b3-80d1-042c452ddc4f",
		Message:   "PIN Valid",
	}

	mockAccountBridge.On("VerifyPIN", ctx, reqCashTx.NoRek, reqCashTx.PIN).Return(resVerifyPIN, nil)

	mockRepo.On("CashWithdrawal", ctx, mock.MatchedBy(func(rw domain.RiwayatTx) bool {
		return rw.NoRekNasabah == reqCashTx.NoRek && rw.JenisTransaksi == "KELUAR" && rw.KategoriTx == "TARIK_TUNAI" && rw.Nominal == reqCashTx.Nominal
	})).Return(int64(100000), nil)

	res, err := u.CashWithdrawal(ctx, reqCashTx)

	assert.NoError(t, err)
	assert.Equal(t, reqCashTx.NoRek, res.NoRek)
	assert.Equal(t, int64(100000), res.TotalSaldo)
	assert.Equal(t, "Tarik tunai berhasil diproses", res.Message)

	mockRepo.AssertExpectations(t)
	mockAccountBridge.AssertExpectations(t)
}

func TestTransfer(t *testing.T) {
	mockRepo := new(mockTransactionRepository)
	mockAccountBridge := new(mockAccountClientBridge)

	u := usecase.NewTransactionUsacase(mockRepo, mockAccountBridge)

	ctx := context.Background()

	reqCashTx := domain.TransferRequest{
		NoRekPengirim: "321039837401",
		NoRekPenerima: "324243536451",
		PIN:           "123123",
		Nominal:       100000,
	}

	resVerifyPIN := domain.VerifyPINResult{
		IsValid:   true,
		NasabahID: "01a0ece9-e8d5-78b3-80d1-042c452ddc4f",
		Message:   "PIN Valid",
	}

	mockAccountBridge.On("VerifyPIN", ctx, reqCashTx.NoRekPengirim, reqCashTx.PIN).Return(resVerifyPIN, nil)

	pengirim := domain.AccountProfileResult{NoRek: reqCashTx.NoRekPengirim, Nama: "Badrul", Status: "AKTIF"}

	penerima := domain.AccountProfileResult{NoRek: reqCashTx.NoRekPenerima, Nama: "Aziz", Status: "AKTIF"}

	mockAccountBridge.On("FindByNoRek", ctx, reqCashTx.NoRekPengirim).Return(pengirim, nil)

	mockAccountBridge.On("FindByNoRek", ctx, reqCashTx.NoRekPenerima).Return(penerima, nil)

	mockRepo.On("Transfer", ctx, mock.MatchedBy(func(rw domain.RiwayatTx) bool {
		return rw.NoRekNasabah == reqCashTx.NoRekPengirim && rw.JenisTransaksi == "KELUAR" && rw.KategoriTx == "TRANSFER" && rw.Nominal == reqCashTx.Nominal
	}), mock.MatchedBy(func(rw domain.RiwayatTx) bool {
		return rw.NoRekNasabah == reqCashTx.NoRekPenerima && rw.JenisTransaksi == "MASUK" && rw.KategoriTx == "TRANSFER" && rw.Nominal == reqCashTx.Nominal
	})).Return(int64(100000), nil)

	res, err := u.Transfer(ctx, reqCashTx)

	assert.NoError(t, err)
	assert.Equal(t, int64(100000), res.SisaSaldo)
	assert.Equal(t, fmt.Sprintf("Transfer sukses ke %s sebesar Rp%d", penerima.Nama, reqCashTx.Nominal), res.Message)

	mockRepo.AssertExpectations(t)
	mockAccountBridge.AssertExpectations(t)
}

func TestBatchGaji(t *testing.T) {
	mockRepo := new(mockTransactionRepository)
	mockAccountBridge := new(mockAccountClientBridge)

	u := usecase.NewTransactionUsacase(mockRepo, mockAccountBridge)

	ctx := context.Background()

	reqBatchGaji := domain.BatchGajiRequest{
		NoRekPerusahaan: "32129137193",
		NamaPerusahaan:  "PT Lokal Jaya",
		PINPerusahaan:   "123123",
		DaftarKaryawan: []domain.GajiItem{
			domain.GajiItem{NoRekKaryawan: "3210313033", NominalGaji: 10000000},
			domain.GajiItem{NoRekKaryawan: "3213538781", NominalGaji: 10000000},
		},
	}

	resVerifyPIN := domain.VerifyPINResult{
		IsValid: true,
		Message: "PIN Valid",
	}

	var total int64 = 20000000

	mockAccountBridge.On("VerifyPIN", ctx, reqBatchGaji.NoRekPerusahaan, reqBatchGaji.PINPerusahaan).Return(resVerifyPIN, nil)

	listKaryawan := []domain.AccountProfileResult{
		domain.AccountProfileResult{
			NoRek:  "3210313033",
			Nama:   "Jihad",
			Status: "AKTIF",
		},
		domain.AccountProfileResult{
			NoRek:  "3213538781",
			Nama:   "Samir",
			Status: "AKTIF",
		},
	}

	for _, v := range listKaryawan {
		mockAccountBridge.On("FindByNoRek", ctx, v.NoRek).Return(v, nil)
	}

	mockRepo.On("SaveBatchGaji", ctx, reqBatchGaji.NoRekPerusahaan, total, mock.MatchedBy(func(rw []domain.RiwayatTx) bool {
		if len(rw) != len(reqBatchGaji.DaftarKaryawan) {
			return false
		}
		return rw[0].JenisTransaksi == "MASUK" && rw[0].KategoriTx == "GAJI"
	})).Return(int32(2), nil)

	res, err := u.ProcessBatchGaji(ctx, reqBatchGaji)

	assert.NoError(t, err)
	assert.Equal(t, "Success", res.Status)
	assert.Equal(t, int32(2), res.TotalDiproses)
	assert.Equal(t, total, res.TotalDanaKeluar)

	mockRepo.AssertExpectations(t)
	mockAccountBridge.AssertExpectations(t)

}

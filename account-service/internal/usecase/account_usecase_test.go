package usecase_test

import (
	"context"
	"mini-bank/account-service/internal/domain"
	"mini-bank/account-service/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAccountRepository struct {
	mock.Mock
}

func (m *mockAccountRepository) CreateWithTransaction(ctx context.Context, nasabah domain.Nasabah, rekening domain.Rekening, refID string, rwTxID string) error {
	args := m.Called(ctx, nasabah, rekening, refID, rwTxID)
	return args.Error(0)
}

func (m *mockAccountRepository) FindByNoRek(ctx context.Context, noRek string) (domain.Nasabah, domain.Rekening, error) {
	args := m.Called(ctx, noRek)
	return args.Get(0).(domain.Nasabah), args.Get(1).(domain.Rekening), args.Error(2)
}

func (m *mockAccountRepository) UpdateStatusDanAttempts(ctx context.Context, noRek string, status string, attempts int) error {
	args := m.Called(ctx, noRek, status, attempts)
	return args.Error(0)
}

func TestRegister_Perorangan_Success(t *testing.T) {
	mockRepo := new(mockAccountRepository)
	u := usecase.NewAccountUsecase(mockRepo)

	reqData := domain.RegisterRequest{
		TipeNasabah:   "PERORANGAN",
		NoIdentitas:   "3201021507950001",
		Nama:          "Dewi Lestari",
		NoHP:          "081345678901",
		Email:         "dewi.lestari@gmail.com",
		AlamatNasabah: "Jl. Tunjungan No. 45",
		Pin:           "120480",
		SetoranAwal:   10000,
	}

	ctx := context.Background()

	mockRepo.On("CreateWithTransaction", ctx, mock.MatchedBy(func(n domain.Nasabah) bool {
		return n.TipeNasabah == "PERORANGAN" && n.Nama == reqData.Nama && n.NoIdentitas == reqData.NoIdentitas
	}), mock.MatchedBy(func(r domain.Rekening) bool {
		return r.Saldo == reqData.SetoranAwal
	}),
		mock.Anything,
		mock.Anything).Return(nil)

	resp, err := u.Register(ctx, reqData)

	assert.NoError(t, err)
	assert.Equal(t, reqData.Nama, resp.Nama)
	assert.NotEmpty(t, resp.NoRek)
	assert.Equal(t, reqData.SetoranAwal, resp.SaldoAwal)

	mockRepo.AssertExpectations(t)
}

func TestRegister_Perusahaan_Success(t *testing.T) {
	mockRepo := new(mockAccountRepository)
	u := usecase.NewAccountUsecase(mockRepo)

	ctx := context.Background()

	reqData := domain.RegisterRequest{
		TipeNasabah:  "PERUSAHAAN",
		Npwp:         "1243213131",
		AlamatKantor: "Jln Mangkuboro",
		NoIdentitas:  "1425367235",
		Nama:         "Perusahaan Maju Mundur",
		NoHP:         "08775663763",
		Email:        "majumundur@gmail.com",
		Pin:          "12330",
		SetoranAwal:  50000,
	}

	mockRepo.On("CreateWithTransaction", ctx, mock.MatchedBy(func(n domain.Nasabah) bool {
		return n.TipeNasabah == "PERUSAHAAN" && n.Nama == reqData.Nama && n.NoIdentitas == reqData.NoIdentitas
	}), mock.MatchedBy(func(r domain.Rekening) bool {
		return r.Saldo == reqData.SetoranAwal
	}), mock.Anything, mock.Anything).Return(nil)

	res, err := u.Register(ctx, reqData)

	assert.NoError(t, err)
	assert.Equal(t, reqData.Nama, res.Nama)
	assert.NotEmpty(t, res.NoRek)
	assert.Equal(t, reqData.SetoranAwal, res.SaldoAwal)

	mockRepo.AssertExpectations(t)

}

func TestCheckProfile(t *testing.T) {
	mockRepo := new(mockAccountRepository)
	u := usecase.NewAccountUsecase(mockRepo)

	ctx := context.Background()

	a := "Jln Mangunsakoro"
	nasabah := domain.Nasabah{
		TipeNasabah:   "PERORANGAN",
		NoIdentitas:   "32131376378",
		Nama:          "Siti Zahwira",
		NoHP:          "086337738",
		Email:         "zahwiraaa@gmail.com",
		AlamatNasabah: &a,
		Pin:           "123123",
	}
	rekening := domain.Rekening{NoRek: "32103094930", Saldo: 10000000}

	_ = mockRepo.On("FindByNoRek", ctx, mock.Anything).Return(nasabah, rekening, nil)

	res, err := u.CheckProfile(ctx, "32103094930")

	assert.NoError(t, err)
	assert.Equal(t, nasabah.Nama, res.Nama)
	assert.Equal(t, rekening.NoRek, res.NoRek)
	assert.Equal(t, rekening.Saldo, res.Saldo)
	assert.Equal(t, rekening.Status, res.Status)

	mockRepo.AssertExpectations(t)
}

func TestVerifyPIN(t *testing.T) {
	mockRepo := new(mockAccountRepository)
	u := usecase.NewAccountUsecase(mockRepo)

	ctx := context.Background()

	a := "Jln Mangunsakoro"
	nasabah := domain.Nasabah{
		TipeNasabah:   "PERORANGAN",
		NoIdentitas:   "32131376378",
		Nama:          "Siti Zahwira",
		NoHP:          "086337738",
		Email:         "zahwiraaa@gmail.com",
		AlamatNasabah: &a,
		Pin:           "$2a$10$AzX",
	}
	rekening := domain.Rekening{NoRek: "32103094930", Saldo: 10000000}

	mockRepo.On("FindByNoRek", ctx, mock.Anything).Return(nasabah, rekening, nil)

	res, err := u.CheckProfile(ctx, "32103094930")

	assert.NoError(t, err)
	assert.Equal(t, nasabah.Nama, res.Nama)
	assert.Equal(t, rekening.NoRek, res.NoRek)
	assert.Equal(t, rekening.Saldo, res.Saldo)
	assert.Equal(t, rekening.Status, res.Status)

	mockRepo.On("UpdateStatusDanAttempts", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err = u.VerifyPIN(ctx, "32103094930", "123123")

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

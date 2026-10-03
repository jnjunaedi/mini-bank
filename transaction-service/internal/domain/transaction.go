package domain

import (
	"context"
	"time"
)

type RiwayatTx struct {
	ID             string
	ReferensiID    string
	NoRekNasabah   string
	JenisTransaksi string
	KategoriTx     string
	Nominal        int64
	NoRekTujuan    string
	TotalSaldo     int64
	Keterangan     string
	CreatedAt      time.Time
}

type AccountProfileResult struct {
	NoRek  string
	Nama   string
	Status string
}

type VerifyPINResult struct {
	IsValid   bool
	NasabahID string
	Message   string
}

type AccountClientBridge interface {
	FindByNoRek(ctx context.Context, noRek string) (AccountProfileResult, error)
	VerifyPIN(ctx context.Context, noRek, pin string) (VerifyPINResult, error)
}

type TransactionRepository interface {
	SaveTransaction(ctx context.Context, riwayatTx RiwayatTx) (int64, error)
	CashWithdrawal(ctx context.Context, riwayatTx RiwayatTx) (int64, error)
	Transfer(ctx context.Context, riwayatTxPengirim, riwayatTxPenerima RiwayatTx) (int64, error)
	SaveBatchGaji(ctx context.Context, noRekPerusahaan string, totalDanaGaji int64, listTx []RiwayatTx) (int32, error)
}

type TransactionUsecase interface {
	CashDeposit(ctx context.Context, req CashTransactionRequest) (CashTransactionResponse, error)
	CashWithdrawal(ctx context.Context, req CashTransactionRequest) (CashTransactionResponse, error)
	Transfer(ctx context.Context, req TransferRequest) (TransferResponse, error)
	ProcessBatchGaji(ctx context.Context, req BatchGajiRequest) (BatchGajiResponse, error)
}

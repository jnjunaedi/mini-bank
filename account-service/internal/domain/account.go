package domain

import (
	"context"
	"time"
)

type Nasabah struct {
	ID            string    `json:"id"`
	TipeNasabah   string    `json:"tipe_nasabah"`
	Npwp          *string   `json:"npwp"`
	AlamatKantor  *string   `json:"alamat_kantor"`
	NoIdentitas   string    `json:"no_identitas"`
	Nama          string    `json:"nama"`
	NoHP          string    `json:"no_hp"`
	Email         string    `json:"email"`
	Pekerjaan     *string   `json:"pekerjaan"`
	AlamatNasabah *string   `json:"alamat_nasabah"`
	Pin           string    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Rekening struct {
	NoRek          string    `json:"no_rek"`
	NasabahID      string    `json:"nasabah_id"`
	Saldo          int64     `json:"saldo"`
	Status         string    `json:"status"`
	FailedAttempts int       `json:"failed_attempts"`
	IsBlocked      bool      `json:"is_blocked"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type AccountRepository interface {
	CreateWithTransaction(ctx context.Context, nasabah Nasabah, rekening Rekening, refID string, rwTxID string) error
	FindByNoRek(ctx context.Context, noRek string) (Nasabah, Rekening, error)
	UpdateStatusDanAttempts(ctx context.Context, noRek string, status string, attempts int) error
}

type AccountUsecase interface {
	Register(ctx context.Context, req RegisterRequest) (RegisterResponse, error)
	CheckProfile(ctx context.Context, noRek string) (AccountDetailResponse, error)
	VerifyPIN(ctx context.Context, noRek string, pin string) (VerifyPINResponse, error)
}

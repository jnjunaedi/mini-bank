package domain

type RegisterResponse struct {
	Nama      string
	NoRek     string
	SaldoAwal int64
}

type AccountDetailResponse struct {
	NasabahID string
	NoRek     string
	Nama      string
	Saldo     int64
	Status    string
}

type VerifyPINResponse struct {
	IsValid   bool
	NasabahID string
	Message   string
}

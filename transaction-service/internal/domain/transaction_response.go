package domain

type CashTransactionResponse struct {
	TransaksiID string
	NoRek       string
	TotalSaldo  int64
	Message     string
}
type BatchGajiResponse struct {
	Status          string
	TotalDiproses   int32
	TotalDanaKeluar int64
	Message         string
}

type TransferResponse struct {
	TransaksiID string
	SisaSaldo   int64
	Message     string
}

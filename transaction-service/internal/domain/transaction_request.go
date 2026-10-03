package domain

type CashTransactionRequest struct {
	NoRek      string
	PIN        string
	Nominal    int64
	Keterangan string
}

type TransferRequest struct {
	NoRekPengirim string
	PIN           string
	NoRekPenerima string
	Nominal       int64
}

type GajiItem struct {
	NoRekKaryawan string
	NominalGaji   int64
}

type BatchGajiRequest struct {
	NoRekPerusahaan string
	NamaPerusahaan  string
	PINPerusahaan   string
	DaftarKaryawan  []GajiItem
}

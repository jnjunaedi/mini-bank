package request

type CashTransactionRequest struct {
	NoRek      string `json:"no_rek" validate:"required,numeric,len=15"`
	PIN        string `json:"pin" validate:"required,numeric,len=6"`
	Nominal    int64  `json:"nominal" validate:"required,min=50000"`
	Keterangan string `json:"keterangan" validate:"omitempty"`
}

type TransferRequest struct {
	NoRek         string `json:"no_rek" validate:"required,numeric,len=15"`
	PIN           string `json:"pin" validate:"required,numeric,len=6"`
	NoRekPenerima string `json:"no_rek_tujuan" validate:"required,numeric,len=15"`
	Nominal       int64  `json:"nominal" validate:"required,min=50000"`
}

type GajiItem struct {
	NoRekKaryawan string `json:"no_rek_karyawan" validate:"required,numeric"`
	NominalGaji   int64  `json:"nominal_gaji" validate:"required,gt=0"`
}

type BatchDepositGajiRequest struct {
	NoRekPerusahaan string     `json:"no_rek" validate:"required,numeric,len=15"`
	NamaPerusahaan  string     `json:"nama_perusahaan" validate:"required"`
	PINPerusahaan   string     `json:"pin" validate:"required,numeric,len=6"`
	GajiItem        []GajiItem `json:"gaji_item" validate:"required,dive"`
}

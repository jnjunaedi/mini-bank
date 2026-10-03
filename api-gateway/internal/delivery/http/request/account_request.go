package request

type AccountRequest struct {
	TipeNasabah    string `json:"tipe_nasabah" validate:"required,oneof=PERORANGAN PERUSAHAAN"`
	Nama           string `json:"nama" validate:"required,min=3"`
	NoHP           string `json:"no_hp" validate:"required,numeric"`
	Email          string `json:"email" validate:"required,email"`
	Pin            string `json:"pin" validate:"required,numeric,len=6"`
	SetoranAwal    int64  `json:"setoran_awal" validate:"required,gt=50000"`
	NomorIdentitas string `json:"nomor_identitas" validate:"required,numeric,min=13,max=16"`
	Npwp           string `json:"npwp,omitempty"`
	Pekerjaan      string `json:"pekerjaan,omitempty"`
	AlamatKantor   string `json:"alamat_kantor,omitempty"`
	AlamatNasabah  string `json:"alamat_nasabah,omitempty"`
}
type VerifyPINRequest struct {
	NoRek string `json:"no_rek" validate:"required,numeric,len=15"`
	Pin   string `json:"pin" validate:"required,numeric,len=6"`
}

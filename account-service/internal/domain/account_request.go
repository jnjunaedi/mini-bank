package domain

type RegisterRequest struct {
	Npwp          string
	AlamatKantor  string
	NoIdentitas   string
	Nama          string
	NoHP          string
	Email         string
	Pekerjaan     string
	AlamatNasabah string
	Pin           string
	SetoranAwal   int64
	TipeNasabah   string
}

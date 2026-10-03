package repository

import (
	"context"
	"database/sql"
	"log"
	"mini-bank/pkg/utils/apperror"
	"mini-bank/transaction-service/internal/domain"

	"github.com/google/uuid"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) domain.TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) SaveTransaction(ctx context.Context, riwayatTx domain.RiwayatTx) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}
	defer tx.Rollback()

	var saldoNasabah int64
	queryCekRekNasabah := `SELECT saldo FROM rekening WHERE no_rek = ?`
	err = tx.QueryRowContext(ctx, queryCekRekNasabah, riwayatTx.NoRekNasabah).Scan(&saldoNasabah)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, apperror.New(apperror.ErrNotFound, "no rekening : "+riwayatTx.NoRekNasabah+" tidak ditemukan")
		}
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	querySimpanSaldo := `UPDATE rekening SET saldo = saldo + ? WHERE no_rek = ?`
	_, err = tx.ExecContext(ctx, querySimpanSaldo, riwayatTx.Nominal, riwayatTx.NoRekNasabah)
	if err != nil {
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	totalSaldo := riwayatTx.Nominal + saldoNasabah
	querySimpanRiwayatTx := `INSERT INTO riwayat_tx (id, referensi_id, no_rek_nasabah, jenis_transaksi, kategori_tx, nominal, no_rek_tujuan, total_saldo, keterangan) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, querySimpanRiwayatTx, riwayatTx.ID, riwayatTx.ReferensiID, riwayatTx.NoRekNasabah, riwayatTx.JenisTransaksi, riwayatTx.KategoriTx, riwayatTx.Nominal, nil, totalSaldo, riwayatTx.Keterangan)
	if err != nil {
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	if err := tx.Commit(); err != nil {
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	return totalSaldo, nil
}

func (r *TransactionRepository) CashWithdrawal(ctx context.Context, riwayatTx domain.RiwayatTx) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}
	defer tx.Rollback()

	var saldoNasabah int64
	queryCekSaldoNasabah := `SELECT saldo FROM rekening WHERE no_rek = ? FOR UPDATE`
	err = tx.QueryRowContext(ctx, queryCekSaldoNasabah, riwayatTx.NoRekNasabah).Scan(&saldoNasabah)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, apperror.New(apperror.ErrNotFound, "no rekening : "+riwayatTx.NoRekNasabah+" tidak ditemukan")
		}
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	if saldoNasabah < riwayatTx.Nominal {
		return 0, apperror.New(apperror.ErrBadRequest, "tidak bisa melakukan penarikan, saldo anda kurang")
	}

	queryUpdateSaldo := `UPDATE rekening SET saldo = saldo - ? WHERE no_rek = ?`
	_, err = tx.ExecContext(ctx, queryUpdateSaldo, riwayatTx.Nominal, riwayatTx.NoRekNasabah)
	if err != nil {
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	saldoAkhir := saldoNasabah - riwayatTx.Nominal
	querySimpanRiwayatTx := `INSERT INTO riwayat_tx (id, referensi_id, no_rek_nasabah, jenis_transaksi, kategori_tx, nominal, no_rek_tujuan, total_saldo, keterangan) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, querySimpanRiwayatTx, riwayatTx.ID, riwayatTx.ReferensiID, riwayatTx.NoRekNasabah, riwayatTx.JenisTransaksi, riwayatTx.KategoriTx, riwayatTx.Nominal, nil, saldoAkhir, riwayatTx.Keterangan)
	if err != nil {
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	if err := tx.Commit(); err != nil {
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	return saldoAkhir, nil

}

func (r *TransactionRepository) Transfer(ctx context.Context, riwayatTxPengirim, riwayatTxPenerima domain.RiwayatTx) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}
	defer tx.Rollback()

	var saldoAwalPengirim int64
	queryCekSaldoPengirim := `SELECT saldo FROM rekening WHERE no_rek = ? FOR UPDATE`
	err = tx.QueryRowContext(ctx, queryCekSaldoPengirim, riwayatTxPengirim.NoRekNasabah).Scan(&saldoAwalPengirim)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, apperror.New(apperror.ErrNotFound, "Nomor rekening pengirim tidak ditemukan")
		}
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	if saldoAwalPengirim < riwayatTxPengirim.Nominal {
		return 0, apperror.New(apperror.ErrBadRequest, "tidak bisa melakukan transfer, saldo anda kurang")
	}

	saldoAkhirPengirim := saldoAwalPengirim - riwayatTxPengirim.Nominal
	queryUpdateSaldo := `UPDATE rekening SET saldo = ? WHERE no_rek = ?`
	_, err = tx.ExecContext(ctx, queryUpdateSaldo, saldoAkhirPengirim, riwayatTxPengirim.NoRekNasabah)
	if err != nil {
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	var saldoAwalPenerima int64
	queryCekSaldoPenerima := `SELECT saldo FROM rekening WHERE no_rek = ? FOR UPDATE`
	err = tx.QueryRowContext(ctx, queryCekSaldoPenerima, riwayatTxPenerima.NoRekNasabah).Scan(&saldoAwalPenerima)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, apperror.New(apperror.ErrNotFound, "no rekening tujuan tidak ditemukan")
		}
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	saldoAkhirPenerima := saldoAwalPenerima + riwayatTxPenerima.Nominal
	queryTransfer := `UPDATE rekening SET saldo = ? WHERE no_rek = ?`
	_, err = tx.ExecContext(ctx, queryTransfer, saldoAkhirPenerima, riwayatTxPenerima.NoRekNasabah)
	if err != nil {
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	querySimpanRiwayatTx := `INSERT INTO riwayat_tx (id, referensi_id, no_rek_nasabah, jenis_transaksi, kategori_tx, nominal, no_rek_tujuan, total_saldo, keterangan) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, querySimpanRiwayatTx, riwayatTxPengirim.ID, riwayatTxPengirim.ReferensiID, riwayatTxPengirim.NoRekNasabah, riwayatTxPengirim.JenisTransaksi, riwayatTxPengirim.KategoriTx, riwayatTxPengirim.Nominal, riwayatTxPenerima.NoRekNasabah, saldoAkhirPengirim, riwayatTxPengirim.Keterangan)
	if err != nil {
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	_, err = tx.ExecContext(ctx, querySimpanRiwayatTx, riwayatTxPenerima.ID, riwayatTxPenerima.ReferensiID, riwayatTxPenerima.NoRekNasabah, riwayatTxPenerima.JenisTransaksi, riwayatTxPenerima.KategoriTx, riwayatTxPenerima.Nominal, riwayatTxPengirim.NoRekNasabah, saldoAkhirPenerima, riwayatTxPenerima.Keterangan)
	if err != nil {
		log.Printf("error : %v\n", err.Error())
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	if err := tx.Commit(); err != nil {
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	return saldoAkhirPengirim, nil
}

func (r *TransactionRepository) SaveBatchGaji(ctx context.Context, noRekPerusahaan string, totalDanaGaji int64, listTx []domain.RiwayatTx) (int32, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}
	defer tx.Rollback()

	var saldoPerusahaan int64
	querCekSaldoPerusahaan := `SELECT saldo FROM rekening WHERE no_rek = ? FOR UPDATE`
	err = tx.QueryRowContext(ctx, querCekSaldoPerusahaan, noRekPerusahaan).Scan(&saldoPerusahaan)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, apperror.New(apperror.ErrInternalServer, "No rekening perusahaan tidak ditemukan")
		}
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	if saldoPerusahaan < totalDanaGaji {
		return 0, apperror.New(apperror.ErrBadRequest, "Saldo perusahaan tidak mencukupi untuk batch payroll")
	}

	queryCekSaldoKaryawan := `SELECT saldo FROM rekening WHERE no_rek = ?`
	stmtCekSaldoKaryawan, err := tx.PrepareContext(ctx, queryCekSaldoKaryawan)
	if err != nil {
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}
	defer stmtCekSaldoKaryawan.Close()

	queryUpdateSaldoKaryawan := `UPDATE rekening SET saldo = saldo + ? WHERE no_rek = ?`
	stmtUpdateSaldoKaryawan, err := tx.PrepareContext(ctx, queryUpdateSaldoKaryawan)
	if err != nil {
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}
	defer stmtUpdateSaldoKaryawan.Close()

	querySimpanRiwayatTx := `INSERT INTO riwayat_tx (id, referensi_id, no_rek_nasabah, jenis_transaksi, kategori_tx, nominal, no_rek_tujuan, total_saldo, keterangan) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	stmtSimpanRwTx, err := tx.PrepareContext(ctx, querySimpanRiwayatTx)
	if err != nil {
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}
	defer stmtSimpanRwTx.Close()

	var jumlahSukses int32

	for _, list := range listTx {
		var totalSaldoKaryawan int64

		err := stmtCekSaldoKaryawan.QueryRowContext(ctx, list.NoRekNasabah).Scan(&totalSaldoKaryawan)
		if err != nil {
			if err == sql.ErrNoRows {
				return 0, apperror.New(apperror.ErrInternalServer, "No rekening "+list.NoRekNasabah+" tidak ditemukan")
			}
			return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
		}

		_, err = stmtUpdateSaldoKaryawan.ExecContext(ctx, list.Nominal, list.NoRekNasabah)
		if err != nil {
			return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
		}
		totalSaldoKaryawan += list.Nominal

		_, err = stmtSimpanRwTx.ExecContext(ctx, list.ID, list.ReferensiID, noRekPerusahaan, list.JenisTransaksi, list.KategoriTx, list.Nominal, list.NoRekNasabah, totalSaldoKaryawan, list.Keterangan)
		if err != nil {
			return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
		}

		saldoPerusahaan -= list.Nominal
		ketKantor := "Transfer Gaji ke " + list.NoRekNasabah

		uuidKantor, _ := uuid.NewV7()
		stmtSimpanRwTx.ExecContext(ctx, uuidKantor, list.ReferensiID, noRekPerusahaan, "KELUAR", "GAJI", list.Nominal, list.NoRekNasabah, totalSaldoKaryawan, ketKantor)

		jumlahSukses++
		log.Println("13")
	}

	queryPotongSaldoPerusahaan := `UPDATE rekening SET saldo = saldo - ? WHERE no_rek = ?`
	_, err = tx.ExecContext(ctx, queryPotongSaldoPerusahaan, totalDanaGaji, noRekPerusahaan)
	if err != nil {
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	if err := tx.Commit(); err != nil {
		return 0, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	return jumlahSukses, nil

}

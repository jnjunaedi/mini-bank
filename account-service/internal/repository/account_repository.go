package repository

import (
	"context"
	"database/sql"
	"errors"
	"mini-bank/account-service/internal/domain"
	"mini-bank/pkg/utils/apperror"
	"time"

	"github.com/go-sql-driver/mysql"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) domain.AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) CreateWithTransaction(ctx context.Context, nasabah domain.Nasabah, rekening domain.Rekening, refID string, rwTxID string) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperror.New(apperror.ErrInternalServer, "error internal server")
	}
	defer tx.Rollback()

	queryNasabah := `INSERT INTO nasabah (id, tipe_nasabah, npwp, alamat_kantor, nomor_identitas, nama, no_hp, email, pekerjaan, alamat_nasabah, pin) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, queryNasabah, nasabah.ID, nasabah.TipeNasabah, nasabah.Npwp, nasabah.AlamatKantor, nasabah.NoIdentitas, nasabah.Nama, nasabah.NoHP, nasabah.Email, nasabah.Pekerjaan, nasabah.AlamatNasabah, nasabah.Pin)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				if nasabah.TipeNasabah == "PERUSAHAAN" {
					return apperror.New(apperror.ErrBadRequest, "nik / nohp / npwp / email sudah digunakan, silahkan gunakan yang lain")
				}
				return apperror.New(apperror.ErrBadRequest, "nik / nohp / email sudah digunakan, silahkan gunakan yang lain")
			}
		}
		return apperror.New(apperror.ErrInternalServer, "error internal server saat menyimpan data nasabah")
	}

	queryRekening := `INSERT INTO rekening (no_rek, nasabah_id, saldo, status) VALUES (?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, queryRekening, rekening.NoRek, rekening.NasabahID, rekening.Saldo, rekening.Status)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				return apperror.New(apperror.ErrBadRequest, "nik / nohp / email sudah digunakan, silahkan gunakan yang lain")
			}
		}
		return apperror.New(apperror.ErrInternalServer, "error internal server saat membuat nomor rekening")
	}

	queryRiwayat := `INSERT INTO riwayat_tx (id, referensi_id, no_rek_nasabah, jenis_transaksi, kategori_tx, nominal, no_rek_tujuan, total_saldo, keterangan) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, queryRiwayat, rwTxID, refID, rekening.NoRek, "MASUK", "SETOR_TUNAI", rekening.Saldo, nil, rekening.Saldo, "Setoran Awal Pendaftaran")
	if err != nil {
		return apperror.New(apperror.ErrInternalServer, "error internal server saat mencatat riwayat transaksi")
	}

	if err := tx.Commit(); err != nil {
		return apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	return nil

}

func (r *AccountRepository) FindByNoRek(ctx context.Context, noRek string) (domain.Nasabah, domain.Rekening, error) {
	query := `SELECT n.id, n.nomor_identitas, n.nama, n.pin, n.tipe_nasabah, r.no_rek, r.saldo, r.status, r.failed_attempts FROM rekening r INNER JOIN nasabah n ON r.nasabah_id = n.id WHERE r.no_rek = ?`

	var nasabah domain.Nasabah
	var rekening domain.Rekening

	err := r.db.QueryRowContext(ctx, query, noRek).Scan(&nasabah.ID, &nasabah.NoIdentitas, &nasabah.Nama, &nasabah.Pin, &nasabah.TipeNasabah, &rekening.NoRek, &rekening.Saldo, &rekening.Status, &rekening.FailedAttempts)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Nasabah{}, domain.Rekening{}, apperror.New(apperror.ErrNotFound, "Nomor rekening tidak valid")
		}
		return domain.Nasabah{}, domain.Rekening{}, apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	return nasabah, rekening, nil
}
func (r *AccountRepository) UpdateStatusDanAttempts(ctx context.Context, noRek string, status string, attempts int) error {
	query := `UPDATE rekening SET status = ?, failed_attempts = ?, updated_at = ? WHERE no_rek = ?`
	_, err := r.db.ExecContext(ctx, query, status, attempts, time.Now(), noRek)
	if err != nil {
		return apperror.New(apperror.ErrInternalServer, "error internal server")
	}

	return nil
}

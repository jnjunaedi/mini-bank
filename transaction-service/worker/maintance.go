package worker

import (
	"database/sql"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

func RegisterCronJobs(db *sql.DB, waktu time.Time, limit int) {
	c := cron.New()

	_, err := c.AddFunc("0 2 * * *", func() {
		log.Println("[Cron] Memulai pemeliharaan status rekening...")

		query := `
			UPDATE rekening r
			SET r.status = 'PASIF', r.updated_at = NOW()
			WHERE r.status = 'AKTIF'
			  AND NOT EXISTS (
			      SELECT 1 
			      FROM riwayat_tx t 
			      WHERE t.no_rek_nasabah = r.nomor_rekening 
			        AND t.created_at >= DATE_SUB(?, INTERVAL ? MONTH)
			  )`

		result, err := db.Exec(query, waktu, limit)
		if err != nil {
			log.Printf("[Cron] Gagal memperbarui status: %v", err)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		log.Printf("[Cron] Selesai. %d rekening diubah menjadi PASIF.", rowsAffected)
	})

	if err != nil {
		log.Fatal("[Cron] Gagal menjadwalkan cron:", err)
	}

	c.Start()
	log.Println("[Cron] Seluruh background job berhasil dijalankan.")
}

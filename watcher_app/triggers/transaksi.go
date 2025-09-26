package trigger

import (
	"fmt"

	"gorm.io/gorm"
)

func SetupTransaksiTriggers(db *gorm.DB) error {
	drops := []string{
		`DROP TRIGGER IF EXISTS transaksi_trigger ON transaksi;
		 DROP FUNCTION IF EXISTS notify_transaksi_change();`,
	}

	for _, drop := range drops {
		if err := db.Exec(drop).Error; err != nil {
			return fmt.Errorf("gagal hapus trigger/function: %w", err)
		}
	}

	triggertransaksi := [...]string{
		`
			CREATE OR REPLACE FUNCTION notify_transaksi_change()
			RETURNS trigger AS $$
			DECLARE
				payload JSON;
			BEGIN
				-- hanya tangani UPDATE
				IF TG_OP = 'UPDATE' THEN
					payload := json_build_object(
						'table', TG_TABLE_NAME,
						'action', TG_OP,
						'id_transaksi', NEW.id,
						'id_pengguna_transaksi', NEW.id_pengguna,
						'id_seller_transaksi', NEW.id_seller,
						'id_barang_induk_transaksi', NEW.id_barang_induk,
						'id_alamat_pengguna_transaksi', NEW.id_alamat_pengguna_transaksi,
						'id_pembayaran_transaksi', NEW.id_pembayaran,
						'kode_order_transaksi', NEW.kode_order,
						'metode_transaksi', NEW.metode,
						'status_transaksi', NEW.status,
						'layanan_pengiriman_kurir_transaksi', NEW.layanan_pengiriman_kurir,
						'jenis_pengiriman_transaksi', NEW.jenis_pengiriman,
						'catatan_transaksi', NEW.catatan,
						'kuantitas_barang_transaksi', NEW.kuantitas_barang,
						'total_transaksi', NEW.total,
						'created_at', NEW.created_at,
						'updated_at', NEW.updated_at,
						'deleted_at', NEW.deleted_at
					);
					PERFORM pg_notify('transaksi_channel', payload::text);
				END IF;

				RETURN NEW;
			END;
			$$ LANGUAGE plpgsql;

			CREATE TRIGGER transaksi_trigger
			AFTER UPDATE ON transaksi
			FOR EACH ROW
			EXECUTE FUNCTION notify_transaksi_change();
		`,
	}

	for _, trig := range triggertransaksi {
		if err := db.Exec(trig).Error; err != nil {
			return fmt.Errorf("gagal buat trigger/function: %w", err)
		}
	}

	fmt.Println("Semua trigger berhasil dibuat ulang!")
	return nil
}

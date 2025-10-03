package trigger

import (
	"fmt"

	"gorm.io/gorm"
)

func SetupPengirimanTriggers(db *gorm.DB) error {
	drops := []string{
		`DROP TRIGGER IF EXISTS info_pengiriman_trigger ON pengiriman;
		 DROP FUNCTION IF EXISTS notify_info_pengiriman_change();`,
	}

	for _, drop := range drops {
		if err := db.Exec(drop).Error; err != nil {
			return fmt.Errorf("gagal hapus trigger/function: %w", err)
		}
	}

	// Trigger untuk info kurir
	triggerPengiriman := []string{
		`
		CREATE OR REPLACE FUNCTION notify_info_pengiriman_change()
		RETURNS trigger AS $$
		DECLARE
			payload JSON;
		BEGIN
			IF TG_OP = 'UPDATE' THEN
				payload := json_build_object(
					'table', TG_TABLE_NAME,
					'action', TG_OP,
					'id_pengiriman', NEW.id,
					'id_transaksi_pengiriman', NEW.id_transaksi,
					'id_alamat_pengambilan_pengiriman', NEW.id_alamat_pengambilan,
					'id_alamat_pengiriman', NEW.id_alamat_pengiriman,
					'id_kurir_pengiriman', NEW.id_kurir,
					'nomor_resi_pengiriman', NEW.nomor_resi,
					'layanan_pengiriman', NEW.layanan_pengiriman_kurir,
					'jenis_pengiriman_transaksi', NEW.jenis_pengiriman,
					'status_pengiriman', NEW.status,
					'biaya_kirim_pengiriman', NEW.biaya_kirim,
					'kurir_paid_pengiriman', NEW.kurir_paid,
					'berat_total_kg_pengiriman', NEW.berat_total_kg
				);
				PERFORM pg_notify('pengiriman_channel', payload::text);
			END IF;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;

		CREATE TRIGGER info_pengiriman_trigger
		AFTER UPDATE ON pengiriman
		FOR EACH ROW
		EXECUTE FUNCTION notify_info_pengiriman_change();
		`,
	}

	for _, trig := range triggerPengiriman {
		if err := db.Exec(trig).Error; err != nil {
			return fmt.Errorf("gagal buat trigger/function: %w", err)
		}
	}

	fmt.Println("Semua trigger berhasil dibuat ulang!")
	return nil
}

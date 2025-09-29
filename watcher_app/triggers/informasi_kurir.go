package trigger

import (
	"fmt"

	"gorm.io/gorm"
)

func SetupInformasiKurirTriggers(db *gorm.DB) error {
	drops := []string{
		// perbaiki typo: IF EXIST -> IF EXISTS
		`DROP TRIGGER IF EXISTS info_kurir_trigger ON informasi_kurir;
		 DROP FUNCTION IF EXISTS notify_info_kurir_change();`,

		`DROP TRIGGER IF EXISTS info_kendaraan_trigger ON informasi_kendaraan_kurir;
		 DROP FUNCTION IF EXISTS notify_info_kendaraan_change();`,
	}

	for _, drop := range drops {
		if err := db.Exec(drop).Error; err != nil {
			return fmt.Errorf("gagal hapus trigger/function: %w", err)
		}
	}

	// Trigger untuk info kurir
	triggerInfoKurir := []string{
		`
		CREATE OR REPLACE FUNCTION notify_info_kurir_change()
		RETURNS trigger AS $$
		DECLARE
			payload JSON;
		BEGIN
			IF TG_OP = 'UPDATE' THEN
				payload := json_build_object(
					'table', TG_TABLE_NAME,
					'action', TG_OP,
					'informasi_id_kurir', NEW.id_kurir,
					'status_perizinan_kurir', NEW.informasi_status_perizinan
				);
				PERFORM pg_notify('informasi_kurir_channel', payload::text);
			END IF;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;

		CREATE TRIGGER info_kurir_trigger
		AFTER UPDATE ON informasi_kurir
		FOR EACH ROW
		EXECUTE FUNCTION notify_info_kurir_change();
		`,
	}

	// Trigger untuk info kendaraan
	triggerInfoKendaraan := []string{
		`
		CREATE OR REPLACE FUNCTION notify_info_kendaraan_change()
		RETURNS trigger AS $$
		DECLARE
			payload JSON;
		BEGIN
			IF TG_OP = 'UPDATE' THEN
				payload := json_build_object(
					'table', TG_TABLE_NAME,
					'action', TG_OP,
					'informasi_id_kurir', NEW.id_kurir,
					'status_perizinan_kurir', NEW.informasi_status_perizinan
				);
				PERFORM pg_notify('informasi_kurir_channel', payload::text);
			END IF;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;

		CREATE TRIGGER info_kendaraan_trigger
		AFTER UPDATE ON informasi_kendaraan_kurir
		FOR EACH ROW
		EXECUTE FUNCTION notify_info_kendaraan_change();
		`,
	}

	// eksekusi semua trigger
	for _, trig := range append(triggerInfoKurir, triggerInfoKendaraan...) {
		if err := db.Exec(trig).Error; err != nil {
			return fmt.Errorf("gagal buat trigger/function: %w", err)
		}
	}

	fmt.Println("Semua trigger berhasil dibuat ulang!")
	return nil
}

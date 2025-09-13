package trigger

import (
	"fmt"

	"gorm.io/gorm"
)

func SetupKomentarTriggers(db *gorm.DB) error {
	drops := []string{
		`DROP TRIGGER IF EXISTS komentar_trigger ON komentar;
		 DROP FUNCTION IF EXISTS notify_komentar_change();`,
	}

	for _, drop := range drops {
		if err := db.Exec(drop).Error; err != nil {
			return fmt.Errorf("gagal hapus trigger/function: %w", err)
		}
	}

	triggerkomentar := [...]string{
		`
		CREATE OR REPLACE FUNCTION notify_komentar_change()
		RETURNS trigger AS $$
		DECLARE
			payload JSON;
		BEGIN
			-- payload berdasarkan operasi
			IF TG_OP = 'INSERT' THEN
				payload := json_build_object(
					'table', TG_TABLE_NAME,
					'action', TG_OP,
					'id_komentar', NEW.id,
					'id_barang_induk', NEW.id_barang_induk,
					'id_entity', NEW.id_entity,
					'jenis_entity', NEW.jenis_entity,
					'isi_komentar', NEW.komentar,
					'parent_id', NEW.parent_id,
					'created_at', NEW.created_at,
					'updated_at', NEW.updated_at,
					'deleted_at', NEW.deleted_at
				);
			ELSIF TG_OP = 'UPDATE' THEN
				payload := json_build_object(
					'table', TG_TABLE_NAME,
					'action', TG_OP,
					'id_komentar', NEW.id,
					'isi_komentar', NEW.komentar
				);
			ELSIF TG_OP = 'DELETE' THEN
				payload := json_build_object(
					'table', TG_TABLE_NAME,
					'action', TG_OP,
					'id_komentar', OLD.id
				);
			END IF;

			-- kirim ke channel komentar
			PERFORM pg_notify('komentar_channel', payload::text);
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;

		CREATE TRIGGER komentar_trigger
		AFTER INSERT OR UPDATE OR DELETE ON komentar
		FOR EACH ROW
		EXECUTE FUNCTION notify_komentar_change();
	`,
	}

	for _, trig := range triggerkomentar {
		if err := db.Exec(trig).Error; err != nil {
			return fmt.Errorf("gagal buat trigger/function: %w", err)
		}
	}

	fmt.Println("Semua trigger berhasil dibuat ulang!")
	return nil
}

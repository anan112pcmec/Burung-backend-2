package trigger

import (
	"fmt"

	"gorm.io/gorm"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk follower
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func FollowerDropper() string {
	return `
	DROP TRIGGER IF EXISTS follower_trigger ON follower;
	DROP FUNCTION IF EXISTS notify_follower_change();
	`
}

func FollowerTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_follower_change()
		RETURNS trigger AS $$
		DECLARE
			payload JSON;
		BEGIN
			-- payload berdasarkan operasi
			IF TG_OP = 'INSERT' THEN
				payload := json_build_object(
					'table', TG_TABLE_NAME,
					'action', TG_OP,
					'id_follower', NEW.id_follower,
					'id_followed', NEW.id_followed
				);
			ELSIF TG_OP = 'DELETE' THEN
				payload := json_build_object(
					'table', TG_TABLE_NAME,
					'action', TG_OP,
					'id_follower', OLD.id_follower,
					'id_followed', OLD.id_followed
				);
			END IF;

			-- kirim ke channel follower
			PERFORM pg_notify('follower_channel', payload::text);
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;

		CREATE TRIGGER follower_trigger
		AFTER INSERT OR DELETE ON follower
		FOR EACH ROW
		EXECUTE FUNCTION notify_follower_change();
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk komentar
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func KomentarDropper() string {
	return `
	DROP TRIGGER IF EXISTS komentar_trigger ON komentar;
	DROP FUNCTION IF EXISTS notify_komentar_change();
	`
}

func KomentarTrigger() string {
	return `
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
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk Informasi Kurir
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func InformasiKurirDropper() string {
	return `
	DROP TRIGGER IF EXISTS info_kurir_trigger ON informasi_kurir;
	DROP FUNCTION IF EXISTS notify_info_kurir_change();
	`
}

func InformasiKurirTrigger() string {
	return `
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
					'status_perizinan_kurir', NEW.informasi_status_perizinan,
					'jenis_kendaraan', 'kosong'
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
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk informasi_kendaraan_kurir
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func InformasiKendaraanKurirDropper() string {
	return `
	DROP TRIGGER IF EXISTS info_kendaraan_trigger ON informasi_kendaraan_kurir;
	DROP FUNCTION IF EXISTS notify_info_kendaraan_change();
	`
}

func InformasiKendaraanKurirTrigger() string {
	return `
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
					'status_perizinan_kurir', NEW.informasi_status_perizinan,
					'jenis_kendaraan', NEW.jenis_kendaraan
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
		`
}

func SetupEngagementEntityTriggers(db *gorm.DB) error {
	drops := []string{
		FollowerDropper(),
		KomentarDropper(),
		InformasiKurirDropper(),
		InformasiKendaraanKurirDropper(),
	}
	for _, drop := range drops {
		if err := db.Exec(drop).Error; err != nil {
			return fmt.Errorf("gagal hapus trigger/function: %w", err)
		}
	}
	trigger := []string{
		FollowerTrigger(),
		KomentarTrigger(),
		InformasiKurirTrigger(),
		InformasiKendaraanKurirTrigger(),
	}
	for _, trig := range trigger {
		if err := db.Exec(trig).Error; err != nil {
			return fmt.Errorf("gagal buat trigger/function: %w", err)
		}
	}

	return nil
}

package trigger

import (
	"fmt"

	"gorm.io/gorm"

)

func PengirimanDropper() string {
	return `
	DROP TRIGGER IF EXISTS pengiriman_trigger ON pengiriman;
	DROP FUNCTION IF EXISTS notify_pengiriman_change();
	`
}

func PengirimanTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_pengiriman_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
		changed_columns JSONB := '{}'::jsonb;
		column_change_name TEXT[] := ARRAY[]::TEXT[];
	BEGIN
		-- =====================================================================
		-- HANDLE UPDATE: kumpulkan perubahan kolom
		-- =====================================================================
		IF TG_OP = 'UPDATE' THEN
			IF OLD.id_transaksi IS DISTINCT FROM NEW.id_transaksi THEN
				changed_columns := jsonb_set(changed_columns, '{id_transaksi}', to_jsonb(NEW.id_transaksi));
				column_change_name := array_append(column_change_name, 'id_transaksi');
			END IF;

			IF OLD.id_alamat_pengambilan IS DISTINCT FROM NEW.id_alamat_pengambilan THEN
				changed_columns := jsonb_set(changed_columns, '{id_alamat_pengambilan}', to_jsonb(NEW.id_alamat_pengambilan));
				column_change_name := array_append(column_change_name, 'id_alamat_pengambilan');
			END IF;

			IF OLD.id_alamat_pengiriman IS DISTINCT FROM NEW.id_alamat_pengiriman THEN
				changed_columns := jsonb_set(changed_columns, '{id_alamat_pengiriman}', to_jsonb(NEW.id_alamat_pengiriman));
				column_change_name := array_append(column_change_name, 'id_alamat_pengiriman');
			END IF;

			IF OLD.id_kurir IS DISTINCT FROM NEW.id_kurir THEN
				changed_columns := jsonb_set(changed_columns, '{id_kurir}', to_jsonb(NEW.id_kurir));
				column_change_name := array_append(column_change_name, 'id_kurir');
			END IF;

			IF OLD.nomor_resi IS DISTINCT FROM NEW.nomor_resi THEN
				changed_columns := jsonb_set(changed_columns, '{nomor_resi}', to_jsonb(NEW.nomor_resi));
				column_change_name := array_append(column_change_name, 'nomor_resi');
			END IF;

			IF OLD.layanan_pengiriman_kurir IS DISTINCT FROM NEW.layanan_pengiriman_kurir THEN
				changed_columns := jsonb_set(changed_columns, '{layanan_pengiriman_kurir}', to_jsonb(NEW.layanan_pengiriman_kurir));
				column_change_name := array_append(column_change_name, 'layanan_pengiriman_kurir');
			END IF;

			IF OLD.jenis_pengiriman IS DISTINCT FROM NEW.jenis_pengiriman THEN
				changed_columns := jsonb_set(changed_columns, '{jenis_pengiriman}', to_jsonb(NEW.jenis_pengiriman));
				column_change_name := array_append(column_change_name, 'jenis_pengiriman');
			END IF;

			IF OLD.status IS DISTINCT FROM NEW.status THEN
				changed_columns := jsonb_set(changed_columns, '{status}', to_jsonb(NEW.status));
				column_change_name := array_append(column_change_name, 'status');
			END IF;

			IF OLD.biaya_kirim IS DISTINCT FROM NEW.biaya_kirim THEN
				changed_columns := jsonb_set(changed_columns, '{biaya_kirim}', to_jsonb(NEW.biaya_kirim));
				column_change_name := array_append(column_change_name, 'biaya_kirim');
			END IF;

			IF OLD.kurir_paid IS DISTINCT FROM NEW.kurir_paid THEN
				changed_columns := jsonb_set(changed_columns, '{kurir_paid}', to_jsonb(NEW.kurir_paid));
				column_change_name := array_append(column_change_name, 'kurir_paid');
			END IF;

			IF OLD.berat_total_kg IS DISTINCT FROM NEW.berat_total_kg THEN
				changed_columns := jsonb_set(changed_columns, '{berat_total_kg}', to_jsonb(NEW.berat_total_kg));
				column_change_name := array_append(column_change_name, 'berat_total_kg');
			END IF;

			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_pengiriman', NEW.id,
				'id_transaksi_pengiriman', NEW.id_transaksi,
				'id_alamat_pengambilan_pengiriman', NEW.id_alamat_pengambilan,
				'id_alamat_pengiriman_pengiriman', NEW.id_alamat_pengiriman,
				'id_kurir_pengiriman', NEW.id_kurir,
				'nomor_resi_pengiriman', NEW.nomor_resi,
				'layanan_pengiriman_kurir', NEW.layanan_pengiriman_kurir,
				'jenis_pengiriman_transaksi', NEW.jenis_pengiriman,
				'status_pengiriman', NEW.status,
				'biaya_kirim_pengiriman', NEW.biaya_kirim,
				'kurir_paid_pengiriman', NEW.kurir_paid,
				'berat_total_kg_pengiriman', NEW.berat_total_kg,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', changed_columns,
				'column_change_name', column_change_name
			);

			PERFORM pg_notify('pengiriman_channel', payload::text);
			RETURN NEW;

		-- =====================================================================
		-- HANDLE INSERT
		-- =====================================================================
		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_pengiriman', NEW.id,
				'id_transaksi_pengiriman', NEW.id_transaksi,
				'id_alamat_pengambilan_pengiriman', NEW.id_alamat_pengambilan,
				'id_alamat_pengiriman_pengiriman', NEW.id_alamat_pengiriman,
				'id_kurir_pengiriman', NEW.id_kurir,
				'nomor_resi_pengiriman', NEW.nomor_resi,
				'layanan_pengiriman_kurir', NEW.layanan_pengiriman_kurir,
				'jenis_pengiriman_transaksi', NEW.jenis_pengiriman,
				'status_pengiriman', NEW.status,
				'biaya_kirim_pengiriman', NEW.biaya_kirim,
				'kurir_paid_pengiriman', NEW.kurir_paid,
				'berat_total_kg_pengiriman', NEW.berat_total_kg,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);

			PERFORM pg_notify('pengiriman_channel', payload::text);
			RETURN NEW;

		-- =====================================================================
		-- HANDLE DELETE
		-- =====================================================================
		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_pengiriman', OLD.id,
				'id_transaksi_pengiriman', OLD.id_transaksi,
				'id_alamat_pengambilan_pengiriman', OLD.id_alamat_pengambilan,
				'id_alamat_pengiriman_pengiriman', OLD.id_alamat_pengiriman,
				'id_kurir_pengiriman', OLD.id_kurir,
				'nomor_resi_pengiriman', OLD.nomor_resi,
				'layanan_pengiriman_kurir', OLD.layanan_pengiriman_kurir,
				'jenis_pengiriman_transaksi', OLD.jenis_pengiriman,
				'status_pengiriman', OLD.status,
				'biaya_kirim_pengiriman', OLD.biaya_kirim,
				'kurir_paid_pengiriman', OLD.kurir_paid,
				'berat_total_kg_pengiriman', OLD.berat_total_kg,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);

			PERFORM pg_notify('pengiriman_channel', payload::text);
			RETURN OLD;
		END IF;

		-- Return value: NEW for INSERT/UPDATE, OLD for DELETE
		IF TG_OP = 'DELETE' THEN
			RETURN OLD;
		ELSE
			RETURN NEW;
		END IF;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER pengiriman_trigger
	AFTER INSERT OR UPDATE OR DELETE ON pengiriman
	FOR EACH ROW
	EXECUTE FUNCTION notify_pengiriman_change();
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk jejak_pengiriman
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func JejakPengirimanDropper() string {
	return `
	DROP TRIGGER IF EXISTS jejak_pengiriman_trigger ON jejak_pengiriman;
	DROP FUNCTION IF EXISTS notify_jejak_pengiriman_change();
	`
}

func JejakPengirimanTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_jejak_pengiriman_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
		changed_columns JSONB := '{}'::jsonb;
		column_change_name TEXT[] := ARRAY[]::TEXT[];
	BEGIN
		-- =====================================================================
		-- HANDLE UPDATE
		-- =====================================================================
		IF TG_OP = 'UPDATE' THEN
			IF OLD.lokasi IS DISTINCT FROM NEW.lokasi THEN
				changed_columns := jsonb_set(changed_columns, '{lokasi}', to_jsonb(NEW.lokasi));
				column_change_name := array_append(column_change_name, 'lokasi');
			END IF;

			IF OLD.keterangan IS DISTINCT FROM NEW.keterangan THEN
				changed_columns := jsonb_set(changed_columns, '{keterangan}', to_jsonb(NEW.keterangan));
				column_change_name := array_append(column_change_name, 'keterangan');
			END IF;

			IF OLD.dicatat_pada IS DISTINCT FROM NEW.dicatat_pada THEN
				changed_columns := jsonb_set(changed_columns, '{dicatat_pada}', to_jsonb(NEW.dicatat_pada));
				column_change_name := array_append(column_change_name, 'dicatat_pada');
			END IF;

			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_jejak_pengiriman', NEW.id,
				'id_pengiriman_jejak_pengiriman', NEW.id_pengiriman,
				'lokasi_jejak_pengiriman', NEW.lokasi,
				'keterangan_jejak_pengiriman', NEW.keterangan,
				'dicatat_pada_jejak_pengiriman', NEW.dicatat_pada,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', changed_columns,
				'column_change_name', column_change_name
			);

			PERFORM pg_notify('jejak_pengiriman_channel', payload::text);
			RETURN NEW;

		-- =====================================================================
		-- HANDLE INSERT
		-- =====================================================================
		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_jejak_pengiriman', NEW.id,
				'id_pengiriman_jejak_pengiriman', NEW.id_pengiriman,
				'lokasi_jejak_pengiriman', NEW.lokasi,
				'keterangan_jejak_pengiriman', NEW.keterangan,
				'dicatat_pada_jejak_pengiriman', NEW.dicatat_pada,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);

			PERFORM pg_notify('jejak_pengiriman_channel', payload::text);
			RETURN NEW;

		-- =====================================================================
		-- HANDLE DELETE
		-- =====================================================================
		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_jejak_pengiriman', OLD.id,
				'id_pengiriman_jejak_pengiriman', OLD.id_pengiriman,
				'lokasi_jejak_pengiriman', OLD.lokasi,
				'keterangan_jejak_pengiriman', OLD.keterangan,
				'dicatat_pada_jejak_pengiriman', OLD.dicatat_pada,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);

			PERFORM pg_notify('jejak_pengiriman_channel', payload::text);
			RETURN OLD;
		END IF;

		-- Default return
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER jejak_pengiriman_trigger
	AFTER INSERT OR UPDATE OR DELETE
	ON jejak_pengiriman
	FOR EACH ROW
	EXECUTE FUNCTION notify_jejak_pengiriman_change();
	`
}

func SetupPengirimanTriggers(db *gorm.DB) error {
	drops := []string{
		PengirimanDropper(),
		JejakPengirimanDropper(),
	}

	for _, drop := range drops {
		if err := db.Exec(drop).Error; err != nil {
			return fmt.Errorf("gagal hapus trigger/function: %w", err)
		}
	}

	// Trigger untuk info kurir
	triggerPengiriman := []string{
		PengirimanTrigger(),
		JejakPengirimanTrigger(),
	}

	for _, trig := range triggerPengiriman {
		if err := db.Exec(trig).Error; err != nil {
			return fmt.Errorf("gagal buat trigger/function: %w", err)
		}
	}

	fmt.Println("Semua trigger berhasil dibuat ulang!")
	return nil
}

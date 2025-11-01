package trigger

import (
	"log"

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
		changed_columns JSONB := '{}'::jsonb;
		column_change_name TEXT[] := ARRAY[]::TEXT[];
	BEGIN
		IF TG_OP = 'UPDATE' THEN
			-- Deteksi perubahan kolom

			IF OLD.id_follower IS DISTINCT FROM NEW.id_follower THEN
				changed_columns := jsonb_set(changed_columns, '{id_follower}', to_jsonb(NEW.id_follower));
				column_change_name := array_append(column_change_name, 'id_follower');
			END IF;

			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_follower', OLD.id_follower,
				'id_followed', OLD.id_followed,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,
				'column_change_name', column_change_name,
				'changed_columns', changed_columns
			);

			PERFORM pg_notify('follower_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_follower', NEW.id_follower,
				'id_followed', NEW.id_followed,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'column_change_name', column_change_name,
				'changed_columns', changed_columns
			);

			PERFORM pg_notify('follower_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_follower', OLD.id_follower,
				'id_followed', OLD.id_followed,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at
			);

			PERFORM pg_notify('follower_channel', payload::text);
			RETURN OLD;
		END IF;

		RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER follower_trigger
	AFTER INSERT OR UPDATE OR DELETE ON follower
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
		changed_columns JSONB := '{}'::jsonb;
		column_change_name TEXT[] := ARRAY[]::TEXT[];
	BEGIN
		IF TG_OP = 'UPDATE' THEN
			-- Deteksi perubahan kolom komentar
			IF OLD.komentar IS DISTINCT FROM NEW.komentar THEN
				changed_columns := jsonb_set(changed_columns, '{isi_komentar}', to_jsonb(NEW.komentar));
				column_change_name := array_append(column_change_name, 'isi_komentar');
			END IF;

			IF OLD.updated_at IS DISTINCT FROM NEW.updated_at THEN
				changed_columns := jsonb_set(changed_columns, '{updated_at}', to_jsonb(NEW.updated_at));
				column_change_name := array_append(column_change_name, 'updated_at');
			END IF;

			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_komentar', NEW.id,
				'id_barang_induk_komentar', NEW.id_barang_induk,
				'id_entity_komentar', NEW.id_entity,
				'jenis_entity_komentar', NEW.jenis_entity,
				'isi_komentar', NEW.komentar,
				'parent_id_komentar', NEW.parent_id,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', changed_columns
				'column_change_name', column_change_name,
			);

			PERFORM pg_notify('komentar_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_komentar', NEW.id,
				'id_barang_induk_komentar', NEW.id_barang_induk,
				'id_entity_komentar', NEW.id_entity,
				'jenis_entity_komentar', NEW.jenis_entity,
				'isi_komentar', NEW.komentar,
				'parent_id_komentar', NEW.parent_id,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);

			PERFORM pg_notify('komentar_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_komentar', OLD.id,
				'id_barang_induk_komentar', OLD.id_barang_induk,
				'id_entity_komentar', OLD.id_entity,
				'jenis_entity_komentar', OLD.jenis_entity,
				'isi_komentar', OLD.komentar,
				'parent_id_komentar', OLD.parent_id,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);

			PERFORM pg_notify('komentar_channel', payload::text);
			RETURN OLD;
		END IF;

		RETURN NULL;
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
	DROP TRIGGER IF EXISTS informasi_kurir_trigger ON informasi_kurir;
	DROP FUNCTION IF EXISTS notify_informasi_kurir_change();
	`
}

func InformasiKurirTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_informasi_kurir_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
		changed_columns JSONB := '{}'::jsonb;
		column_change_name TEXT[] := ARRAY[]::TEXT[];
	BEGIN
		-- ======================================================
		-- HANDLE UPDATE: deteksi perubahan per kolom
		-- ======================================================
		IF TG_OP = 'UPDATE' THEN
			IF OLD.umur IS DISTINCT FROM NEW.umur THEN
				changed_columns := jsonb_set(changed_columns, '{umur_informasi_kurir}', to_jsonb(NEW.umur));
				column_change_name := array_append(column_change_name, 'umur_informasi_kurir');
			END IF;

			IF OLD.alasan IS DISTINCT FROM NEW.alasan THEN
				changed_columns := jsonb_set(changed_columns, '{alasan_informasi_kurir}', to_jsonb(NEW.alasan));
				column_change_name := array_append(column_change_name, 'alasan_informasi_kurir');
			END IF;

			IF OLD.informasi_ktp IS DISTINCT FROM NEW.informasi_ktp THEN
				changed_columns := jsonb_set(changed_columns, '{informasi_ktp_informasi_kurir}', to_jsonb(NEW.informasi_ktp));
				column_change_name := array_append(column_change_name, 'informasi_ktp_informasi_kurir');
			END IF;

			IF OLD.alamat IS DISTINCT FROM NEW.alamat THEN
				changed_columns := jsonb_set(changed_columns, '{alamat_informasi_kurir}', to_jsonb(NEW.alamat));
				column_change_name := array_append(column_change_name, 'alamat_informasi_kurir');
			END IF;

			IF OLD.status IS DISTINCT FROM NEW.status THEN
				changed_columns := jsonb_set(changed_columns, '{status_informasi_kurir}', to_jsonb(NEW.status));
				column_change_name := array_append(column_change_name, 'status_informasi_kurir');
			END IF;

			IF OLD.informasi_status_perizinan IS DISTINCT FROM NEW.informasi_status_perizinan THEN
				changed_columns := jsonb_set(changed_columns, '{status_perizinan_kurir}', to_jsonb(NEW.informasi_status_perizinan));
				column_change_name := array_append(column_change_name, 'status_perizinan_kurir');
			END IF;

			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_informasi_kurir', OLD.id,
				'id_kurir_informasi_kurir', OLD.id_kurir,
				'umur_informasi_kurir', OLD.umur,
				'alasan_informasi_kurir', OLD.alasan,
				'informasi_ktp_informasi_kurir', OLD.informasi_ktp,
				'alamat_informasi_kurir', OLD.alamat,
				'status_informasi_kurir', OLD.status,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,
				'changed_columns', changed_columns,
				'column_change_name', column_change_name
			);

			PERFORM pg_notify('informasi_kurir_channel', payload::text);
			RETURN NEW;

		-- ======================================================
		-- HANDLE INSERT
		-- ======================================================
		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_informasi_kurir', NEW.id,
				'id_kurir_informasi_kurir', NEW.id_kurir,
				'umur_informasi_kurir', NEW.umur,
				'alasan_informasi_kurir', NEW.alasan,
				'informasi_ktp_informasi_kurir', NEW.informasi_ktp,
				'alamat_informasi_kurir', NEW.alamat,
				'status_informasi_kurir', NEW.status,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', changed_columns,
				'column_change_name', column_change_name
			);

			PERFORM pg_notify('informasi_kurir_channel', payload::text);
			RETURN NEW;

		-- ======================================================
		-- HANDLE DELETE
		-- ======================================================
		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_informasi_kurir', OLD.id,
				'id_kurir_informasi_kurir', OLD.id_kurir,
				'umur_informasi_kurir', OLD.umur,
				'alasan_informasi_kurir', OLD.alasan,
				'informasi_ktp_informasi_kurir', OLD.informasi_ktp,
				'alamat_informasi_kurir', OLD.alamat,
				'status_informasi_kurir', OLD.status,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at
			);

			PERFORM pg_notify('informasi_kurir_channel', payload::text);
			RETURN OLD;
		END IF;

		RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER informasi_kurir_trigger
	AFTER INSERT OR UPDATE OR DELETE ON informasi_kurir
	FOR EACH ROW
	EXECUTE FUNCTION notify_informasi_kurir_change();
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk informasi_kendaraan_kurir
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func InformasiKendaraanKurirDropper() string {
	return `
	DROP TRIGGER IF EXISTS informasi_kendaraan_trigger ON informasi_kendaraan_kurir;
	DROP FUNCTION IF EXISTS notify_informasi_kendaraan_change();
	`
}

func InformasiKendaraanKurirTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_informasi_kendaraan_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
		changed_columns JSONB := '{}'::jsonb;
		column_change_name TEXT[] := ARRAY[]::TEXT[];
	BEGIN
		-- ======================================================
		-- HANDLE INSERT
		-- ======================================================
		IF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_informasi_kendaraan_kurir', NEW.id,
				'id_kurir_informasi_kendaraan_kurir', NEW.id_kurir,
				'jenis_kendaraan_informasi_kendaraan_kurir', NEW.jenis_kendaraan_kurir,
				'nama_kendaraan_informasi_kendaraan_kurir', NEW.nama_kendaraan,
				'roda_kendaraan_informasi_kendaraan_kurir', NEW.roda_kendaraan,
				'informasi_stnk_informasi_kendaraan_kurir', NEW.informasi_stnk,
				'informasi_bpkb_informasi_kendaraan_kurir', NEW.informasi_bpkb,
				'status_informasi_kendaraan_kurir', NEW.status,
				'created_at', NOW(),
				'updated_at', NOW(),
				'deleted_at', NULL,
				'changed_columns', changed_columns,
				'column_change_name', column_change_name
			);
			PERFORM pg_notify('informasi_kendaraan_kurir_channel', payload::text);
			RETURN NEW;
		END IF;

		-- ======================================================
		-- HANDLE UPDATE (deteksi perubahan kolom)
		-- ======================================================
		IF TG_OP = 'UPDATE' THEN
			IF OLD.id_kurir IS DISTINCT FROM NEW.id_kurir THEN
				changed_columns := jsonb_set(changed_columns, '{id_kurir_informasi_kendaraan_kurir}', to_jsonb(NEW.id_kurir));
				column_change_name := array_append(column_change_name, 'id_kurir_informasi_kendaraan_kurir');
			END IF;

			IF OLD.jenis_kendaraan_kurir IS DISTINCT FROM NEW.jenis_kendaraan_kurir THEN
				changed_columns := jsonb_set(changed_columns, '{jenis_kendaraan_informasi_kendaraan_kurir}', to_jsonb(NEW.jenis_kendaraan_kurir));
				column_change_name := array_append(column_change_name, 'jenis_kendaraan_informasi_kendaraan_kurir');
			END IF;

			IF OLD.nama_kendaraan IS DISTINCT FROM NEW.nama_kendaraan THEN
				changed_columns := jsonb_set(changed_columns, '{nama_kendaraan_informasi_kendaraan_kurir}', to_jsonb(NEW.nama_kendaraan));
				column_change_name := array_append(column_change_name, 'nama_kendaraan_informasi_kendaraan_kurir');
			END IF;

			IF OLD.roda_kendaraan IS DISTINCT FROM NEW.roda_kendaraan THEN
				changed_columns := jsonb_set(changed_columns, '{roda_kendaraan_informasi_kendaraan_kurir}', to_jsonb(NEW.roda_kendaraan));
				column_change_name := array_append(column_change_name, 'roda_kendaraan_informasi_kendaraan_kurir');
			END IF;

			IF OLD.informasi_stnk IS DISTINCT FROM NEW.informasi_stnk THEN
				changed_columns := jsonb_set(changed_columns, '{informasi_stnk_informasi_kendaraan_kurir}', to_jsonb(NEW.informasi_stnk));
				column_change_name := array_append(column_change_name, 'informasi_stnk_informasi_kendaraan_kurir');
			END IF;

			IF OLD.informasi_bpkb IS DISTINCT FROM NEW.informasi_bpkb THEN
				changed_columns := jsonb_set(changed_columns, '{informasi_bpkb_informasi_kendaraan_kurir}', to_jsonb(NEW.informasi_bpkb));
				column_change_name := array_append(column_change_name, 'informasi_bpkb_informasi_kendaraan_kurir');
			END IF;

			IF OLD.status IS DISTINCT FROM NEW.status THEN
				changed_columns := jsonb_set(changed_columns, '{status_informasi_kendaraan_kurir}', to_jsonb(NEW.status));
				column_change_name := array_append(column_change_name, 'status_informasi_kendaraan_kurir');
			END IF;

			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_informasi_kendaraan_kurir', NEW.id,
				'id_kurir_informasi_kendaraan_kurir', NEW.id_kurir,
				'jenis_kendaraan_informasi_kendaraan_kurir', NEW.jenis_kendaraan_kurir,
				'nama_kendaraan_informasi_kendaraan_kurir', NEW.nama_kendaraan,
				'roda_kendaraan_informasi_kendaraan_kurir', NEW.roda_kendaraan,
				'informasi_stnk_informasi_kendaraan_kurir', NEW.informasi_stnk,
				'informasi_bpkb_informasi_kendaraan_kurir', NEW.informasi_bpkb,
				'status_informasi_kendaraan_kurir', NEW.status,
				'created_at', OLD.created_at,
				'updated_at', NOW(),
				'deleted_at', OLD.deleted_at,
				'changed_columns', changed_columns,
				'column_change_name', column_change_name
			);
			PERFORM pg_notify('informasi_kendaraan_kurir_channel', payload::text);
			RETURN NEW;
		END IF;

		-- ======================================================
		-- HANDLE DELETE
		-- ======================================================
		IF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_informasi_kendaraan_kurir', OLD.id,
				'id_kurir_informasi_kendaraan_kurir', OLD.id_kurir,
				'jenis_kendaraan_informasi_kendaraan_kurir', OLD.jenis_kendaraan_kurir,
				'nama_kendaraan_informasi_kendaraan_kurir', OLD.nama_kendaraan,
				'roda_kendaraan_informasi_kendaraan_kurir', OLD.roda_kendaraan,
				'informasi_stnk_informasi_kendaraan_kurir', OLD.informasi_stnk,
				'informasi_bpkb_informasi_kendaraan_kurir', OLD.informasi_bpkb,
				'status_informasi_kendaraan_kurir', OLD.status,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', NOW(),
				'changed_columns', changed_columns,
				'column_change_name', column_change_name
			);
			PERFORM pg_notify('informasi_kendaraan_kurir_channel', payload::text);
			RETURN OLD;
		END IF;

		RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER informasi_kendaraan_trigger
	AFTER INSERT OR UPDATE OR DELETE ON informasi_kendaraan_kurir
	FOR EACH ROW
	EXECUTE FUNCTION notify_informasi_kendaraan_change();
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk keranjang
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func KeranjangDropper() string {
	return `
	DROP TRIGGER IF EXISTS keranjang_trigger ON keranjang;
	DROP FUNCTION IF EXISTS notify_keranjang_change();
	`
}

func KeranjangTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_keranjang_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
		changed_columns JSONB := '{}'::jsonb;
		column_change_name TEXT[] := ARRAY[]::TEXT[];
	BEGIN
		-- ============================================================
		-- HANDLE UPDATE
		-- ============================================================
		IF TG_OP = 'UPDATE' THEN
			IF OLD.id_seller IS DISTINCT FROM NEW.id_seller THEN
				changed_columns := jsonb_set(changed_columns, '{id_seller}', to_jsonb(NEW.id_seller));
				column_change_name := array_append(column_change_name, 'id_seller');
			END IF;
			IF OLD.id_barang_induk IS DISTINCT FROM NEW.id_barang_induk THEN
				changed_columns := jsonb_set(changed_columns, '{id_barang_induk}', to_jsonb(NEW.id_barang_induk));
				column_change_name := array_append(column_change_name, 'id_barang_induk');
			END IF;
			IF OLD.id_kategori_barang IS DISTINCT FROM NEW.id_kategori_barang THEN
				changed_columns := jsonb_set(changed_columns, '{id_kategori_barang}', to_jsonb(NEW.id_kategori_barang));
				column_change_name := array_append(column_change_name, 'id_kategori_barang');
			END IF;
			IF OLD.count IS DISTINCT FROM NEW.count THEN
				changed_columns := jsonb_set(changed_columns, '{count}', to_jsonb(NEW.count));
				column_change_name := array_append(column_change_name, 'count');
			END IF;
			IF OLD.status IS DISTINCT FROM NEW.status THEN
				changed_columns := jsonb_set(changed_columns, '{status}', to_jsonb(NEW.status));
				column_change_name := array_append(column_change_name, 'status');
			END IF;

			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_pengguna_keranjang', OLD.id,
				'id_seller_barang_induk_keranjang', OLD.id_seller,
				'id_barang_induk_keranjang', OLD.id_barang_induk,
				'id_kategori_barang_keranjang', OLD.id_kategori_barang,
				'count_keranjang', OLD.count,
				'status_keranjang', OLD.status,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,
				'changed_columns', changed_columns,
				'column_change_name', column_change_name
			);

		-- ============================================================
		-- HANDLE INSERT
		-- ============================================================
		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_pengguna_keranjang', NEW.id,
				'id_seller_barang_induk_keranjang', NEW.id_seller,
				'id_barang_induk_keranjang', NEW.id_barang_induk,
				'id_kategori_barang_keranjang', NEW.id_kategori_barang,
				'count_keranjang', NEW.count,
				'status_keranjang', NEW.status,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);

		-- ============================================================
		-- HANDLE DELETE
		-- ============================================================
		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_pengguna_keranjang', OLD.id,
				'id_seller_barang_induk_keranjang', OLD.id_seller,
				'id_barang_induk_keranjang', OLD.id_barang_induk,
				'id_kategori_barang_keranjang', OLD.id_kategori_barang,
				'count_keranjang', OLD.count,
				'status_keranjang', OLD.status,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
		END IF;

		PERFORM pg_notify('keranjang_channel', payload::text);
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER keranjang_trigger
	AFTER INSERT OR UPDATE OR DELETE ON keranjang
	FOR EACH ROW
	EXECUTE FUNCTION notify_keranjang_change();
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk barang_disukai
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func BarangDisukaiDropper() string {
	return `
	DROP TRIGGER IF EXISTS barang_disukai_trigger ON barang_disukai;
	DROP FUNCTION IF EXISTS notify_barang_disukai_change();
	`
}

func BarangDisukaiTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_barang_disukai_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
		changed_columns JSONB := '{}'::jsonb;
		column_change_name TEXT[] := ARRAY[]::TEXT[];
	BEGIN
		IF TG_OP = 'UPDATE' THEN
			-- Deteksi kolom yang berubah
			IF OLD.id_pengguna IS DISTINCT FROM NEW.id_pengguna THEN
				changed_columns := jsonb_set(changed_columns, '{id_pengguna_barang_disukai}', to_jsonb(NEW.id_pengguna));
				column_change_name := array_append(column_change_name, 'id_pengguna_barang_disukai');
			END IF;

			IF OLD.id_barang_induk IS DISTINCT FROM NEW.id_barang_induk THEN
				changed_columns := jsonb_set(changed_columns, '{id_barang_induk_disukai}', to_jsonb(NEW.id_barang_induk));
				column_change_name := array_append(column_change_name, 'id_barang_induk_disukai');
			END IF;

			IF OLD.updated_at IS DISTINCT FROM NEW.updated_at THEN
				changed_columns := jsonb_set(changed_columns, '{updated_at}', to_jsonb(NEW.updated_at));
				column_change_name := array_append(column_change_name, 'updated_at');
			END IF;

			-- Bangun payload untuk UPDATE
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'changed_columns', changed_columns,
				'column_change_name', column_change_name,
				'id_pengguna_barang_disukai', NEW.id_pengguna,
				'id_barang_induk_disukai', NEW.id_barang_induk,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at
			);
			PERFORM pg_notify('barang_disukai_channel', payload::text);

		ELSIF TG_OP = 'INSERT' THEN
			-- Bangun payload untuk INSERT
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_pengguna_barang_disukai', NEW.id_pengguna,
				'id_barang_induk_disukai', NEW.id_barang_induk,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('barang_disukai_channel', payload::text);

		ELSIF TG_OP = 'DELETE' THEN
			-- Bangun payload untuk DELETE
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_pengguna_barang_disukai', OLD.id_pengguna,
				'id_barang_induk_disukai', OLD.id_barang_induk,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('barang_disukai_channel', payload::text);
		END IF;

		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	DROP TRIGGER IF EXISTS barang_disukai_trigger ON barang_disukai;
	CREATE TRIGGER barang_disukai_trigger
	AFTER INSERT OR UPDATE OR DELETE ON barang_disukai
	FOR EACH ROW
	EXECUTE FUNCTION notify_barang_disukai_change();
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk entity_social_media
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func EntitySocialMediaDropper() string {
	return `	
	DROP TRIGGER IF EXISTS entity_social_media_trigger ON entity_social_media;
	DROP FUNCTION IF EXISTS notify_entity_social_media_change();
	`
}

func EntitySocialMediaTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_entity_social_media_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
		changed_columns JSONB := '{}'::jsonb;
		column_change_name TEXT[] := ARRAY[]::TEXT[];
	BEGIN
		-- =====================================================================
		-- HANDLE UPDATE: deteksi kolom yang berubah
		-- =====================================================================
		IF TG_OP = 'UPDATE' THEN
			IF OLD.entity_id IS DISTINCT FROM NEW.entity_id THEN
				changed_columns := jsonb_set(changed_columns, '{entity_id_social_media}', to_jsonb(NEW.entity_id));
				column_change_name := array_append(column_change_name, 'entity_id_social_media');
			END IF;

			IF OLD.whatsapp IS DISTINCT FROM NEW.whatsapp THEN
				changed_columns := jsonb_set(changed_columns, '{whatsapp_social_media}', to_jsonb(NEW.whatsapp));
				column_change_name := array_append(column_change_name, 'whatsapp_social_media');
			END IF;

			IF OLD.facebook IS DISTINCT FROM NEW.facebook THEN
				changed_columns := jsonb_set(changed_columns, '{facebook_social_media}', to_jsonb(NEW.facebook));
				column_change_name := array_append(column_change_name, 'facebook_social_media');
			END IF;

			IF OLD.tiktok IS DISTINCT FROM NEW.tiktok THEN
				changed_columns := jsonb_set(changed_columns, '{tiktok_social_media}', to_jsonb(NEW.tiktok));
				column_change_name := array_append(column_change_name, 'tiktok_social_media');
			END IF;

			IF OLD.instagram IS DISTINCT FROM NEW.instagram THEN
				changed_columns := jsonb_set(changed_columns, '{instagram_social_media}', to_jsonb(NEW.instagram));
				column_change_name := array_append(column_change_name, 'instagram_social_media');
			END IF;

			IF OLD.entity_type IS DISTINCT FROM NEW.entity_type THEN
				changed_columns := jsonb_set(changed_columns, '{entity_type_social_media}', to_jsonb(NEW.entity_type));
				column_change_name := array_append(column_change_name, 'entity_type_social_media');
			END IF;

			IF OLD.updated_at IS DISTINCT FROM NEW.updated_at THEN
				changed_columns := jsonb_set(changed_columns, '{updated_at}', to_jsonb(NEW.updated_at));
				column_change_name := array_append(column_change_name, 'updated_at');
			END IF;

			-- Bangun payload JSON untuk UPDATE
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'changed_columns', changed_columns,
				'column_change_name', column_change_name,
				'id_social_media', NEW.id,
				'entity_id_social_media', NEW.entity_id,
				'whatsapp_social_media', NEW.whatsapp,
				'facebook_social_media', NEW.facebook,
				'tiktok_social_media', NEW.tiktok,
				'instagram_social_media', NEW.instagram,
				'entity_type_social_media', NEW.entity_type,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at
			);
			PERFORM pg_notify('entity_social_media_channel', payload::text);

		-- =====================================================================
		-- HANDLE INSERT
		-- =====================================================================
		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_social_media', NEW.id,
				'entity_id_social_media', NEW.entity_id,
				'whatsapp_social_media', NEW.whatsapp,
				'facebook_social_media', NEW.facebook,
				'tiktok_social_media', NEW.tiktok,
				'instagram_social_media', NEW.instagram,
				'entity_type_social_media', NEW.entity_type,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('entity_social_media_channel', payload::text);

		-- =====================================================================
		-- HANDLE DELETE
		-- =====================================================================
		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_social_media', OLD.id,
				'entity_id_social_media', OLD.entity_id,
				'whatsapp_social_media', OLD.whatsapp,
				'facebook_social_media', OLD.facebook,
				'tiktok_social_media', OLD.tiktok,
				'instagram_social_media', OLD.instagram,
				'entity_type_social_media', OLD.entity_type,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('entity_social_media_channel', payload::text);
		END IF;

		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	DROP TRIGGER IF EXISTS entity_social_media_trigger ON entity_social_media;
	CREATE TRIGGER entity_social_media_trigger
	AFTER INSERT OR UPDATE OR DELETE ON entity_social_media
	FOR EACH ROW
	EXECUTE FUNCTION notify_entity_social_media_change();
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk aktivitas_pengguna
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func AktivitasPenggunaDropper() string {
	return `
	DROP TRIGGER IF EXISTS aktivitas_pengguna_trigger ON aktivitas_pengguna;
	DROP FUNCTION IF EXISTS notify_aktivitas_pengguna_change();
	`
}

func AktivitasPenggunaTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_aktivitas_pengguna_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
		changed_columns JSONB := '{}'::jsonb;
		column_change_name TEXT[] := ARRAY[]::TEXT[];
	BEGIN
		-- =====================================================================
		-- HANDLE UPDATE: deteksi perubahan kolom
		-- =====================================================================
		IF TG_OP = 'UPDATE' THEN
			IF OLD.id_pengguna IS DISTINCT FROM NEW.id_pengguna THEN
				changed_columns := jsonb_set(changed_columns, '{id_pengguna_aktivitas_pengguna}', to_jsonb(NEW.id_pengguna));
				column_change_name := array_append(column_change_name, 'id_pengguna_aktivitas_pengguna');
			END IF;

			IF OLD.waktu_dilakukan IS DISTINCT FROM NEW.waktu_dilakukan THEN
				changed_columns := jsonb_set(changed_columns, '{waktu_dilakukan_aktivitas_pengguna}', to_jsonb(NEW.waktu_dilakukan));
				column_change_name := array_append(column_change_name, 'waktu_dilakukan_aktivitas_pengguna');
			END IF;

			IF OLD.aksi IS DISTINCT FROM NEW.aksi THEN
				changed_columns := jsonb_set(changed_columns, '{aksi_aktivitas_pengguna}', to_jsonb(NEW.aksi));
				column_change_name := array_append(column_change_name, 'aksi_aktivitas_pengguna');
			END IF;

			IF OLD.updated_at IS DISTINCT FROM NEW.updated_at THEN
				changed_columns := jsonb_set(changed_columns, '{updated_at}', to_jsonb(NEW.updated_at));
				column_change_name := array_append(column_change_name, 'updated_at');
			END IF;

			-- Payload JSON untuk UPDATE
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'changed_columns', changed_columns,
				'column_change_name', column_change_name,
				'id_aktivitas_pengguna', NEW.id,
				'id_pengguna_aktivitas_pengguna', NEW.id_pengguna,
				'waktu_dilakukan_aktivitas_pengguna', NEW.waktu_dilakukan,
				'aksi_aktivitas_pengguna', NEW.aksi,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at
			);
			PERFORM pg_notify('aktivitas_pengguna_channel', payload::text);

		-- =====================================================================
		-- HANDLE INSERT
		-- =====================================================================
		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_aktivitas_pengguna', NEW.id,
				'id_pengguna_aktivitas_pengguna', NEW.id_pengguna,
				'waktu_dilakukan_aktivitas_pengguna', NEW.waktu_dilakukan,
				'aksi_aktivitas_pengguna', NEW.aksi,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('aktivitas_pengguna_channel', payload::text);

		-- =====================================================================
		-- HANDLE DELETE
		-- =====================================================================
		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_aktivitas_pengguna', OLD.id,
				'id_pengguna_aktivitas_pengguna', OLD.id_pengguna,
				'waktu_dilakukan_aktivitas_pengguna', OLD.waktu_dilakukan,
				'aksi_aktivitas_pengguna', OLD.aksi,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('aktivitas_pengguna_channel', payload::text);
		END IF;

		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	DROP TRIGGER IF EXISTS aktivitas_pengguna_trigger ON aktivitas_pengguna;
	CREATE TRIGGER aktivitas_pengguna_trigger
	AFTER INSERT OR UPDATE OR DELETE ON aktivitas_pengguna
	FOR EACH ROW
	EXECUTE FUNCTION notify_aktivitas_pengguna_change();
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk aktivitas_seller
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func AktivitasSellerDropper() string {
	return `
	DROP TRIGGER IF EXISTS aktivitas_seller_trigger ON aktivitas_seller;
	DROP FUNCTION IF EXISTS notify_aktivitas_seller_change();
	`
}

func AktivitasSellerTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_aktivitas_seller_change()
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
			IF OLD.id_seller IS DISTINCT FROM NEW.id_seller THEN
				changed_columns := jsonb_set(changed_columns, '{id_seller}', to_jsonb(NEW.id_seller));
				column_change_name := array_append(column_change_name, 'id_seller');
			END IF;

			IF OLD.waktu_dilakukan IS DISTINCT FROM NEW.waktu_dilakukan THEN
				changed_columns := jsonb_set(changed_columns, '{waktu_dilakukan}', to_jsonb(NEW.waktu_dilakukan));
				column_change_name := array_append(column_change_name, 'waktu_dilakukan');
			END IF;

			IF OLD.aksi IS DISTINCT FROM NEW.aksi THEN
				changed_columns := jsonb_set(changed_columns, '{aksi}', to_jsonb(NEW.aksi));
				column_change_name := array_append(column_change_name, 'aksi');
			END IF;

			IF OLD.updated_at IS DISTINCT FROM NEW.updated_at THEN
				changed_columns := jsonb_set(changed_columns, '{updated_at}', to_jsonb(NEW.updated_at));
				column_change_name := array_append(column_change_name, 'updated_at');
			END IF;

			IF OLD.deleted_at IS DISTINCT FROM NEW.deleted_at THEN
				changed_columns := jsonb_set(changed_columns, '{deleted_at}', to_jsonb(NEW.deleted_at));
				column_change_name := array_append(column_change_name, 'deleted_at');
			END IF;

			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_aktivitas_seller', NEW.id,
				'id_seller_aktivitas_seller', NEW.id_seller,
				'waktu_dilakukan_aktivitas_seller', NEW.waktu_dilakukan,
				'aksi_aktivitas_seller', NEW.aksi,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', changed_columns,
				'column_change_name', column_change_name
			);

			PERFORM pg_notify('aktivitas_seller_channel', payload::text);
			RETURN NEW;

		-- =====================================================================
		-- HANDLE INSERT
		-- =====================================================================
		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_aktivitas_seller', NEW.id,
				'id_seller_aktivitas_seller', NEW.id_seller,
				'waktu_dilakukan_aktivitas_seller', NEW.waktu_dilakukan,
				'aksi_aktivitas_seller', NEW.aksi,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);

			PERFORM pg_notify('aktivitas_seller_channel', payload::text);
			RETURN NEW;

		-- =====================================================================
		-- HANDLE DELETE
		-- =====================================================================
		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_aktivitas_seller', OLD.id,
				'id_seller_aktivitas_seller', OLD.id_seller,
				'waktu_dilakukan_aktivitas_seller', OLD.waktu_dilakukan,
				'aksi_aktivitas_seller', OLD.aksi,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);

			PERFORM pg_notify('aktivitas_seller_channel', payload::text);
			RETURN OLD;
		END IF;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER aktivitas_seller_trigger
	AFTER INSERT OR UPDATE OR DELETE ON aktivitas_seller
	FOR EACH ROW
	EXECUTE FUNCTION notify_aktivitas_seller_change();
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk alamat_pengguna
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func AlamatPenggunaDropper() string {
	return `
	DROP TRIGGER IF EXISTS alamat_pengguna_trigger ON alamat_pengguna;
	DROP FUNCTION IF EXISTS notify_alamat_pengguna_change();
	`
}

func AlamatPenggunaTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_alamat_pengguna_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
		changed_columns JSONB := '{}'::jsonb;
		column_change_name TEXT[] := ARRAY[]::TEXT[];
	BEGIN
		-- =========================================================
		-- HANDLE UPDATE
		-- =========================================================
		IF TG_OP = 'UPDATE' THEN
			IF OLD.id_pengguna IS DISTINCT FROM NEW.id_pengguna THEN
				changed_columns := jsonb_set(changed_columns, '{id_pengguna}', to_jsonb(NEW.id_pengguna));
				column_change_name := array_append(column_change_name, 'id_pengguna');
			END IF;

			IF OLD.panggilan_alamat IS DISTINCT FROM NEW.panggilan_alamat THEN
				changed_columns := jsonb_set(changed_columns, '{panggilan_alamat}', to_jsonb(NEW.panggilan_alamat));
				column_change_name := array_append(column_change_name, 'panggilan_alamat');
			END IF;

			IF OLD.nomor_telefon IS DISTINCT FROM NEW.nomor_telefon THEN
				changed_columns := jsonb_set(changed_columns, '{nomor_telefon}', to_jsonb(NEW.nomor_telefon));
				column_change_name := array_append(column_change_name, 'nomor_telefon');
			END IF;

			IF OLD.nama_alamat IS DISTINCT FROM NEW.nama_alamat THEN
				changed_columns := jsonb_set(changed_columns, '{nama_alamat}', to_jsonb(NEW.nama_alamat));
				column_change_name := array_append(column_change_name, 'nama_alamat');
			END IF;

			IF OLD.kota IS DISTINCT FROM NEW.kota THEN
				changed_columns := jsonb_set(changed_columns, '{kota}', to_jsonb(NEW.kota));
				column_change_name := array_append(column_change_name, 'kota');
			END IF;

			IF OLD.kode_pos IS DISTINCT FROM NEW.kode_pos THEN
				changed_columns := jsonb_set(changed_columns, '{kode_pos}', to_jsonb(NEW.kode_pos));
				column_change_name := array_append(column_change_name, 'kode_pos');
			END IF;

			IF OLD.kode_negara IS DISTINCT FROM NEW.kode_negara THEN
				changed_columns := jsonb_set(changed_columns, '{kode_negara}', to_jsonb(NEW.kode_negara));
				column_change_name := array_append(column_change_name, 'kode_negara');
			END IF;

			IF OLD.deskripsi IS DISTINCT FROM NEW.deskripsi THEN
				changed_columns := jsonb_set(changed_columns, '{deskripsi}', to_jsonb(NEW.deskripsi));
				column_change_name := array_append(column_change_name, 'deskripsi');
			END IF;

			IF OLD.longitude IS DISTINCT FROM NEW.longitude THEN
				changed_columns := jsonb_set(changed_columns, '{longitude}', to_jsonb(NEW.longitude));
				column_change_name := array_append(column_change_name, 'longitude');
			END IF;

			IF OLD.latitude IS DISTINCT FROM NEW.latitude THEN
				changed_columns := jsonb_set(changed_columns, '{latitude}', to_jsonb(NEW.latitude));
				column_change_name := array_append(column_change_name, 'latitude');
			END IF;

			IF OLD.updated_at IS DISTINCT FROM NEW.updated_at THEN
				changed_columns := jsonb_set(changed_columns, '{updated_at}', to_jsonb(NEW.updated_at));
				column_change_name := array_append(column_change_name, 'updated_at');
			END IF;

			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_alamat_user', NEW.id,
				'id_pengguna_alamat_user', NEW.id_pengguna,
				'panggilan_alamat_user', NEW.panggilan_alamat,
				'nomor_telfon_alamat_user', NEW.nomor_telefon,
				'nama_alamat_user', NEW.nama_alamat,
				'kota_alamat_user', NEW.kota,
				'kode_pos_alamat_user', NEW.kode_pos,
				'kode_negara_alamat_user', NEW.kode_negara,
				'deskripsi_alamat_user', NEW.deskripsi,
				'longitude_alamat_user', NEW.longitude,
				'latitude_alamat_user', NEW.latitude,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', changed_columns,
				'column_change_name', column_change_name
			);

			PERFORM pg_notify('alamat_pengguna_channel', payload::text);
			RETURN NEW;

		-- =========================================================
		-- HANDLE INSERT
		-- =========================================================
		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_alamat_user', NEW.id,
				'id_pengguna_alamat_user', NEW.id_pengguna,
				'panggilan_alamat_user', NEW.panggilan_alamat,
				'nomor_telfon_alamat_user', NEW.nomor_telefon,
				'nama_alamat_user', NEW.nama_alamat,
				'kota_alamat_user', NEW.kota,
				'kode_pos_alamat_user', NEW.kode_pos,
				'kode_negara_alamat_user', NEW.kode_negara,
				'deskripsi_alamat_user', NEW.deskripsi,
				'longitude_alamat_user', NEW.longitude,
				'latitude_alamat_user', NEW.latitude,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('alamat_pengguna_channel', payload::text);
			RETURN NEW;

		-- =========================================================
		-- HANDLE DELETE
		-- =========================================================
		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_alamat_user', OLD.id,
				'id_pengguna_alamat_user', OLD.id_pengguna,
				'panggilan_alamat_user', OLD.panggilan_alamat,
				'nomor_telfon_alamat_user', OLD.nomor_telefon,
				'nama_alamat_user', OLD.nama_alamat,
				'kota_alamat_user', OLD.kota,
				'kode_pos_alamat_user', OLD.kode_pos,
				'kode_negara_alamat_user', OLD.kode_negara,
				'deskripsi_alamat_user', OLD.deskripsi,
				'longitude_alamat_user', OLD.longitude,
				'latitude_alamat_user', OLD.latitude,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('alamat_pengguna_channel', payload::text);
			RETURN OLD;
		END IF;

		RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER alamat_pengguna_trigger
	AFTER INSERT OR UPDATE OR DELETE ON alamat_pengguna
	FOR EACH ROW
	EXECUTE FUNCTION notify_alamat_pengguna_change();
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk jenis_seller_validation
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func JenisSellerValidationDropper() string {
	return `
	DROP TRIGGER IF EXISTS jenis_seller_validation_trigger ON jenis_seller_validation;
	DROP FUNCTION IF EXISTS notify_jenis_seller_validation_change();
	`
}

func JenisSellerValidationTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_jenis_seller_validation_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
	BEGIN
		IF TG_OP = 'UPDATE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_jenis_seller', NEW.id,
				'id_seller_jenis_seller', NEW.id_seller,
				'validation_status_jenis_seller', NEW.validation_status,
				'alasan_seller_jenis_seller', NEW.alasan_seller,
				'alasan_admin_jenis_seller', NEW.alasan_admin,
				'target_jenis_seller', NEW.target_jenis,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('jenis_seller_validation_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_jenis_seller', NEW.id,
				'id_seller_jenis_seller', NEW.id_seller,
				'validation_status_jenis_seller', NEW.validation_status,
				'alasan_seller_jenis_seller', NEW.alasan_seller,
				'alasan_admin_jenis_seller', NEW.alasan_admin,
				'target_jenis_seller', NEW.target_jenis,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('jenis_seller_validation_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_jenis_seller', OLD.id,
				'id_seller_jenis_seller', OLD.id_seller,
				'validation_status_jenis_seller', OLD.validation_status,
				'alasan_seller_jenis_seller', OLD.alasan_seller,
				'alasan_admin_jenis_seller', OLD.alasan_admin,
				'target_jenis_seller', OLD.target_jenis,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('jenis_seller_validation_channel', payload::text);
			RETURN OLD;
		END IF;

		RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER jenis_seller_validation_trigger
	AFTER INSERT OR UPDATE OR DELETE ON jenis_seller_validation
	FOR EACH ROW
	EXECUTE FUNCTION notify_jenis_seller_validation_change();
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk alamat_seller
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func AlamatSellerDropper() string {
	return `
	DROP TRIGGER IF EXISTS alamat_seller_trigger ON alamat_seller;
	DROP FUNCTION IF EXISTS notify_alamat_seller_change();
	`
}

func AlamatSellerTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_alamat_seller_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
	BEGIN
		IF TG_OP = 'UPDATE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_alamat_seller', NEW.id,
				'id_seller_alamat_seller', NEW.id_seller,
				'panggilan_alamat_seller', NEW.panggilan_alamat,
				'nomor_telfon_alamat_seller', NEW.nomor_telefon,
				'nama_alamat_seller', NEW.nama_alamat,
				'deskripsi_alamat_seller', NEW.deskripsi,
				'longitude_alamat_seller', NEW.longitude,
				'latitude_alamat_seller', NEW.latitude,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('alamat_seller_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_alamat_seller', NEW.id,
				'id_seller_alamat_seller', NEW.id_seller,
				'panggilan_alamat_seller', NEW.panggilan_alamat,
				'nomor_telfon_alamat_seller', NEW.nomor_telefon,
				'nama_alamat_seller', NEW.nama_alamat,
				'deskripsi_alamat_seller', NEW.deskripsi,
				'longitude_alamat_seller', NEW.longitude,
				'latitude_alamat_seller', NEW.latitude,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('alamat_seller_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_alamat_seller', OLD.id,
				'id_seller_alamat_seller', OLD.id_seller,
				'panggilan_alamat_seller', OLD.panggilan_alamat,
				'nomor_telfon_alamat_seller', OLD.nomor_telefon,
				'nama_alamat_seller', OLD.nama_alamat,
				'deskripsi_alamat_seller', OLD.deskripsi,
				'longitude_alamat_seller', OLD.longitude,
				'latitude_alamat_seller', OLD.latitude,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('alamat_seller_channel', payload::text);
			RETURN OLD;
		END IF;

		RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER alamat_seller_trigger
	AFTER INSERT OR UPDATE OR DELETE ON alamat_seller
	FOR EACH ROW
	EXECUTE FUNCTION notify_alamat_seller_change();
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk batal_transaksi
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func BatalTransaksiDropper() string {
	return `
	DROP TRIGGER IF EXISTS batal_transaksi_trigger ON batal_transaksi;
	DROP FUNCTION IF EXISTS notify_batal_transaksi_change();
	`
}

func BatalTransaksiTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_batal_transaksi_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
	BEGIN
		IF TG_OP = 'UPDATE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_batal_transaksi', NEW.id,
				'id_transaksi_batal_transaksi', NEW.id_transaksi,
				'transaksi_dibatalkan_oleh', NEW.dibatalkan_oleh,
				'alasan_batal_transaksi', NEW.alasan,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('batal_transaksi_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_batal_transaksi', NEW.id,
				'id_transaksi_batal_transaksi', NEW.id_transaksi,
				'transaksi_dibatalkan_oleh', NEW.dibatalkan_oleh,
				'alasan_batal_transaksi', NEW.alasan,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('batal_transaksi_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_batal_transaksi', OLD.id,
				'id_transaksi_batal_transaksi', OLD.id_transaksi,
				'transaksi_dibatalkan_oleh', OLD.dibatalkan_oleh,
				'alasan_batal_transaksi', OLD.alasan,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('batal_transaksi_channel', payload::text);
			RETURN OLD;
		END IF;

		RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER batal_transaksi_trigger
	AFTER INSERT OR UPDATE OR DELETE ON batal_transaksi
	FOR EACH ROW
	EXECUTE FUNCTION notify_batal_transaksi_change();
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk diskon
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func DiskonDropper() string {
	return `
	DROP TRIGGER IF EXISTS diskon_trigger ON diskon;
	DROP FUNCTION IF EXISTS notify_diskon_change();
	`
}

func DiskonTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_diskon_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
	BEGIN
		IF TG_OP = 'UPDATE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_barang_induk_diskon', NEW.id_barang_induk,
				'deskripsi_diskon', NEW.berlaku,
				'expired_diskon', NEW.expired,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('diskon_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_barang_induk_diskon', NEW.id_barang_induk,
				'deskripsi_diskon', NEW.berlaku,
				'expired_diskon', NEW.expired,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('diskon_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_barang_induk_diskon', OLD.id_barang_induk,
				'deskripsi_diskon', OLD.berlaku,
				'expired_diskon', OLD.expired,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('diskon_channel', payload::text);
			RETURN OLD;
		END IF;

		RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER diskon_trigger
	AFTER INSERT OR UPDATE OR DELETE ON diskon
	FOR EACH ROW
	EXECUTE FUNCTION notify_diskon_change();
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk rekening_seller
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func RekeningSellerDropper() string {
	return `
	DROP TRIGGER IF EXISTS rekening_seller_trigger ON rekening_seller;
	DROP FUNCTION IF EXISTS notify_rekening_seller_change();
	`
}

func RekeningSellerTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_rekening_seller_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
	BEGIN
		IF TG_OP = 'UPDATE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_rekening_seller', NEW.id,
				'id_seller', NEW.id_seller,
				'nama_bank_rekening_seller', NEW.nama_bank,
				'nomor_rekening_seller', NEW.nomor_rekening,
				'pemilik_rekening_seller', NEW.pemilik_rekening,
				'is_default_rekening_seller', NEW.id_default,
				'status_rekening_seller', NEW.status,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('rekening_seller_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_rekening_seller', NEW.id,
				'id_seller', NEW.id_seller,
				'nama_bank_rekening_seller', NEW.nama_bank,
				'nomor_rekening_seller', NEW.nomor_rekening,
				'pemilik_rekening_seller', NEW.pemilik_rekening,
				'is_default_rekening_seller', NEW.id_default,
				'status_rekening_seller', NEW.status,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('rekening_seller_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_rekening_seller', OLD.id,
				'id_seller', OLD.id_seller,
				'nama_bank_rekening_seller', OLD.nama_bank,
				'nomor_rekening_seller', OLD.nomor_rekening,
				'pemilik_rekening_seller', OLD.pemilik_rekening,
				'is_default_rekening_seller', OLD.id_default,
				'status_rekening_seller', OLD.status,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('rekening_seller_channel', payload::text);
			RETURN OLD;
		END IF;

		RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER rekening_seller_trigger
	AFTER INSERT OR UPDATE OR DELETE ON rekening_seller
	FOR EACH ROW
	EXECUTE FUNCTION notify_rekening_seller_change();
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk balance_kurir_log
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func BalanceKurirLogDropper() string {
	return `
	DROP TRIGGER IF EXISTS balance_kurir_log_trigger ON balance_kurir_log;
	DROP FUNCTION IF EXISTS notify_balance_kurir_log_change();
	`
}

func BalanceKurirLogTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_balance_kurir_log_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
	BEGIN
		IF TG_OP = 'UPDATE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_balance_kurir', NEW.id,
				'kurir_id', NEW.kurir_id,
				'amount_balance_kurir', NEW.amount,
				'type_balance_kurir', NEW.type,
				'catatan_balance_kurir', NEW.catatan,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('balance_kurir_log_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_balance_kurir', NEW.id,
				'kurir_id', NEW.kurir_id,
				'amount_balance_kurir', NEW.amount,
				'type_balance_kurir', NEW.type,
				'catatan_balance_kurir', NEW.catatan,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('balance_kurir_log_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_balance_kurir', OLD.id,
				'kurir_id', OLD.kurir_id,
				'amount_balance_kurir', OLD.amount,
				'type_balance_kurir', OLD.type,
				'catatan_balance_kurir', OLD.catatan,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('balance_kurir_log_channel', payload::text);
			RETURN OLD;
		END IF;

		RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER balance_kurir_log_trigger
	AFTER INSERT OR UPDATE OR DELETE ON balance_kurir_log
	FOR EACH ROW
	EXECUTE FUNCTION notify_balance_kurir_log_change();
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk alamat_gudang
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func AlamatGudangDropper() string {
	return `
	DROP TRIGGER IF EXISTS alamat_gudang_trigger ON alamat_gudang;
	DROP FUNCTION IF EXISTS notify_alamat_gudang_change();
	`
}

func AlamatGudangTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_alamat_gudang_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
	BEGIN
		IF TG_OP = 'UPDATE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_alamat_gudang', NEW.id,
				'id_seller_alamat_gudang', NEW.id_seller,
				'panggilan_alamat_gudang', NEW.panggilan_alamat,
				'nomor_telefon_alamat_gudang', NEW.nomor_telefon,
				'nama_alamat_gudang', NEW.nama_alamat,
				'kota_alamat_gudang', NEW.kota,
				'kode_pos_alamat_gudang', NEW.kode_pos,
				'kode_negara_alamat_gudang', NEW.kode_negara,
				'deskripsi_alamat_gudang', NEW.deskripsi,
				'longitude_alamat_gudang', NEW.longitude,
				'latitude_alamat_gudang', NEW.latitude,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('alamat_gudang_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_alamat_gudang', NEW.id,
				'id_seller_alamat_gudang', NEW.id_seller,
				'panggilan_alamat_gudang', NEW.panggilan_alamat,
				'nomor_telefon_alamat_gudang', NEW.nomor_telefon,
				'nama_alamat_gudang', NEW.nama_alamat,
				'kota_alamat_gudang', NEW.kota,
				'kode_pos_alamat_gudang', NEW.kode_pos,
				'kode_negara_alamat_gudang', NEW.kode_negara,
				'deskripsi_alamat_gudang', NEW.deskripsi,
				'longitude_alamat_gudang', NEW.longitude,
				'latitude_alamat_gudang', NEW.latitude,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('alamat_gudang_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,

				'id_alamat_gudang', OLD.id,
				'id_seller_alamat_gudang', OLD.id_seller,
				'panggilan_alamat_gudang', OLD.panggilan_alamat,
				'nomor_telefon_alamat_gudang', OLD.nomor_telefon,
				'nama_alamat_gudang', OLD.nama_alamat,
				'kota_alamat_gudang', OLD.kota,
				'kode_pos_alamat_gudang', OLD.kode_pos,
				'kode_negara_alamat_gudang', OLD.kode_negara,
				'deskripsi_alamat_gudang', OLD.deskripsi,
				'longitude_alamat_gudang', OLD.longitude,
				'latitude_alamat_gudang', OLD.latitude,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,

				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('alamat_gudang_channel', payload::text);
			RETURN OLD;
		END IF;

		RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER alamat_gudang_trigger
	AFTER INSERT OR UPDATE OR DELETE ON alamat_gudang
	FOR EACH ROW
	EXECUTE FUNCTION notify_alamat_gudang_change();
	`
}

func SetupEngagementEntityTriggers(db *gorm.DB) error {
	// 1️⃣ Drop dulu semuanya (best-effort)
	drops := []string{
		FollowerDropper(),
		KomentarDropper(),
		InformasiKurirDropper(),
		InformasiKendaraanKurirDropper(),
		KeranjangDropper(),
		BarangDisukaiDropper(),
		EntitySocialMediaDropper(),
		AktivitasPenggunaDropper(),
		AktivitasSellerDropper(),
		AlamatPenggunaDropper(),
		JenisSellerValidationDropper(),
		AlamatSellerDropper(),
		BatalTransaksiDropper(),
		DiskonDropper(),
		RekeningSellerDropper(),
		BalanceKurirLogDropper(),
		AlamatGudangDropper(),
	}

	for _, drop := range drops {
		if err := db.Exec(drop).Error; err != nil {
			log.Printf("[WARN] Gagal drop trigger/function: %v", err)
			continue // lanjut aja
		}
	}

	// 2️⃣ Buat ulang semuanya (best-effort juga)
	triggers := []string{
		FollowerTrigger(),
		KomentarTrigger(),
		InformasiKurirTrigger(),
		InformasiKendaraanKurirTrigger(),
		KeranjangTrigger(),
		BarangDisukaiTrigger(),
		EntitySocialMediaTrigger(),
		AktivitasPenggunaTrigger(),
		AktivitasSellerTrigger(),
		AlamatPenggunaTrigger(),
		JenisSellerValidationTrigger(),
		AlamatSellerTrigger(),
		BatalTransaksiTrigger(),
		DiskonTrigger(),
		RekeningSellerTrigger(),
		BalanceKurirLogTrigger(),
		AlamatGudangTrigger(),
	}

	for _, trig := range triggers {
		if err := db.Exec(trig).Error; err != nil {
			log.Printf("[ERROR] Gagal buat trigger/function: %v", err)
			continue // tetap lanjut ke berikutnya
		}
	}

	log.Println("[DONE] Semua trigger entity engagement diproses (best-effort mode).")
	return nil

}

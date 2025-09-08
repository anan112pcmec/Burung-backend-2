package trigger

import (
	"fmt"

	"gorm.io/gorm"
)

func SetupBarangTriggers(db *gorm.DB) error {
	drops := []string{
		`DROP TRIGGER IF EXISTS barang_induk_trigger ON barang_induk;
		 DROP FUNCTION IF EXISTS notify_barang_induk_change();`,
	}

	for _, drop := range drops {
		if err := db.Exec(drop).Error; err != nil {
			return fmt.Errorf("gagal hapus trigger/function: %w", err)
		}
	}

	triggerbarang := [...]string{
		`
	CREATE OR REPLACE FUNCTION notify_barang_induk_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
		changed_columns JSON := '{}';
	BEGIN
		IF TG_OP = 'UPDATE' THEN
			IF OLD.id_seller IS DISTINCT FROM NEW.id_seller THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{id_seller}', to_jsonb(NEW.id_seller));
			END IF;
			IF OLD.nama_barang IS DISTINCT FROM NEW.nama_barang THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{nama_barang}', to_jsonb(NEW.nama_barang));
			END IF;
			IF OLD.jenis_barang IS DISTINCT FROM NEW.jenis_barang THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{jenis_barang}', to_jsonb(NEW.jenis_barang));
			END IF;
			IF OLD.original_kategori IS DISTINCT FROM NEW.original_kategori THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{original_kategori}', to_jsonb(NEW.original_kategori));
			END IF;
			IF OLD.deskripsi IS DISTINCT FROM NEW.deskripsi THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{deskripsi}', to_jsonb(NEW.deskripsi));
			END IF;
			IF OLD.tanggal_rilis IS DISTINCT FROM NEW.tanggal_rilis THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{tanggal_rilis}', to_jsonb(NEW.tanggal_rilis));
			END IF;
			IF OLD.viewed IS DISTINCT FROM NEW.viewed THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{viewed}', to_jsonb(NEW.viewed));
			END IF;
			IF OLD.likes IS DISTINCT FROM NEW.likes THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{likes}', to_jsonb(NEW.likes));
			END IF;
			IF OLD.total_komentar IS DISTINCT FROM NEW.total_komentar THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{total_komentar}', to_jsonb(NEW.total_komentar));
			END IF;
		END IF;

		-- Payload sesuai action
		IF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_barang_induk', NEW.id,
				'id_seller_barang_induk', NEW.id_seller,
				'nama_barang_induk', NEW.nama_barang,
				'jenis_barang_induk', NEW.jenis_barang,
				'original_kategori', NEW.original_kategori,
				'deskripsi_barang_induk', NEW.deskripsi,
				'tanggal_rilis_barang_induk', NEW.tanggal_rilis,
				'viewed_barang_induk', NEW.viewed,
				'likes_barang_induk', NEW.likes,
				'total_komentar_barang_induk', NEW.total_komentar,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at
			);
		ELSIF TG_OP = 'UPDATE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_barang_induk', NEW.id,
				'changed_columns', changed_columns,
				'nama_barang_induk', NEW.nama_barang,
				'jenis_barang_induk', NEW.jenis_barang,
				'viewed_barang_induk', NEW.viewed,
				'likes_barang_induk', NEW.likes,
				'total_komentar_barang_induk', NEW.total_komentar,
				'updated_at', NEW.updated_at
			);
		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_barang_induk', OLD.id,
				'nama_barang_induk', OLD.nama_barang,
				'jenis_barang_induk', OLD.jenis_barang,
				'viewed_barang_induk', OLD.viewed,
				'likes_barang_induk', OLD.likes,
				'total_komentar_barang_induk', OLD.total_komentar,
				'deleted_at', OLD.deleted_at
			);
		END IF;

		PERFORM pg_notify('barang_induk_channel', payload::text);
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER barang_induk_trigger
	AFTER INSERT OR UPDATE OR DELETE ON barang_induk
	FOR EACH ROW
	EXECUTE FUNCTION notify_barang_induk_change();
	`,
	}

	for _, trig := range triggerbarang {
		if err := db.Exec(trig).Error; err != nil {
			return fmt.Errorf("gagal buat trigger/function: %w", err)
		}
	}

	fmt.Println("Semua trigger berhasil dibuat ulang!")
	return nil
}

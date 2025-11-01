package trigger

import (
	"fmt"

	"gorm.io/gorm"
)

func PenggunaDropper() string {
	return `
	DROP TRIGGER IF EXISTS pengguna_trigger ON pengguna;
	DROP FUNCTION IF EXISTS notify_pengguna_change();
	`
}

func PenggunaTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_pengguna_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
		changed_columns JSONB := '{}'::jsonb;
		column_change_name TEXT[] := ARRAY[]::TEXT[];
	BEGIN
		IF TG_OP = 'UPDATE' THEN
			-- deteksi semua kolom yang berubah
			IF OLD.username IS DISTINCT FROM NEW.username THEN
				changed_columns := changed_columns || jsonb_build_object('username', NEW.username);
				column_change_name := array_append(column_change_name, 'username');
			END IF;

			IF OLD.email IS DISTINCT FROM NEW.email THEN
				changed_columns := changed_columns || jsonb_build_object('email', NEW.email);
				column_change_name := array_append(column_change_name, 'email');
			END IF;

			IF OLD.pin_hash IS DISTINCT FROM NEW.pin_hash THEN
				changed_columns := changed_columns || jsonb_build_object('pin_hash', 'diubah');
				column_change_name := array_append(column_change_name, 'pin_hash');
			END IF;

			IF OLD.password_hash IS DISTINCT FROM NEW.password_hash THEN
				changed_columns := changed_columns || jsonb_build_object('password_hash', 'diubah');
				column_change_name := array_append(column_change_name, 'password_hash');
			END IF;

			IF OLD.status IS DISTINCT FROM NEW.status THEN
				changed_columns := changed_columns || jsonb_build_object('status', NEW.status);
				column_change_name := array_append(column_change_name, 'status');
			END IF;

			IF OLD.nama IS DISTINCT FROM NEW.nama THEN
				changed_columns := changed_columns || jsonb_build_object('nama', NEW.nama);
				column_change_name := array_append(column_change_name, 'nama');
			END IF;

			IF OLD.updated_at IS DISTINCT FROM NEW.updated_at THEN
				changed_columns := changed_columns || jsonb_build_object('updated_at', NEW.updated_at);
				column_change_name := array_append(column_change_name, 'updated_at');
			END IF;

			IF OLD.deleted_at IS DISTINCT FROM NEW.deleted_at THEN
				changed_columns := changed_columns || jsonb_build_object('deleted_at', NEW.deleted_at);
				column_change_name := array_append(column_change_name, 'deleted_at');
			END IF;

			-- payload notifikasi UPDATE
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_user', OLD.id,
				'username_user', OLD.username,
				'nama_user', OLD.nama,
				'email_user', OLD.email,
				'status_user', OLD.status,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,
				'changed_columns', changed_columns,
				'column_change_name', column_change_name
			);
			PERFORM pg_notify('pengguna_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_user', NEW.id,
				'username_user', NEW.username,
				'nama_user', NEW.nama,
				'email_user', NEW.email,
				'status_user', NEW.status,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('pengguna_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_user', OLD.id,
				'username_user', OLD.username,
				'nama_user', OLD.nama,
				'email_user', OLD.email,
				'status_user', OLD.status,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('pengguna_channel', payload::text);
			RETURN OLD;
		END IF;

		RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER pengguna_trigger
	AFTER INSERT OR UPDATE OR DELETE ON pengguna
	FOR EACH ROW
	EXECUTE FUNCTION notify_pengguna_change();
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk seller
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func SellerDropper() string {
	return `
	DROP TRIGGER IF EXISTS seller_trigger ON seller;
	DROP FUNCTION IF EXISTS notify_seller_change();
	`
}

func SellerTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_seller_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
		changed_columns JSONB := '{}'::jsonb;
		column_change_name TEXT[] := ARRAY[]::TEXT[];
	BEGIN
		IF TG_OP = 'UPDATE' THEN
			-- Deteksi perubahan kolom
			IF OLD.username IS DISTINCT FROM NEW.username THEN
				changed_columns := jsonb_set(changed_columns, '{username}', to_jsonb(NEW.username));
				column_change_name := array_append(column_change_name, 'username');
			END IF;

			IF OLD.email IS DISTINCT FROM NEW.email THEN
				changed_columns := jsonb_set(changed_columns, '{email}', to_jsonb(NEW.email));
				column_change_name := array_append(column_change_name, 'email');
			END IF;

			IF OLD.jenis IS DISTINCT FROM NEW.jenis THEN
				changed_columns := jsonb_set(changed_columns, '{jenis}', to_jsonb(NEW.jenis));
				column_change_name := array_append(column_change_name, 'jenis');
			END IF;

			

			IF OLD.seller_dedication IS DISTINCT FROM NEW.seller_dedication THEN
				changed_columns := jsonb_set(changed_columns, '{seller_dedication}', to_jsonb(NEW.seller_dedication));
				column_change_name := array_append(column_change_name, 'seller_dedication');
			END IF;

			IF OLD.jam_operasional IS DISTINCT FROM NEW.jam_operasional THEN
				changed_columns := jsonb_set(changed_columns, '{jam_operasional}', to_jsonb(NEW.jam_operasional));
				column_change_name := array_append(column_change_name, 'jam_operasional');
			END IF;

			IF OLD.password_hash IS DISTINCT FROM NEW.password_hash THEN
				changed_columns := jsonb_set(changed_columns, '{password_hash}', to_jsonb('diubah'));
				column_change_name := array_append(column_change_name, 'password_hash');
			END IF;

			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_seller', OLD.id,
				'username_seller', OLD.username,
				'nama_seller', OLD.nama,
				'email_seller', OLD.email,
				'jenis_seller', OLD.jenis,
				'jam_operasional_seller', OLD.jam_operasional,
				'punchline_seller', OLD.punchline,
				'deskripsi_seller', OLD.deskripsi,
				'seller_dedication', OLD.seller_dedication,
				'follower_total_seller', OLD.follower_total,
				'status_seller', OLD.status,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', changed_columns,
				'column_change_name', column_change_name
			);
			PERFORM pg_notify('seller_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_seller', NEW.id,
				'username_seller', NEW.username,
				'nama_seller', NEW.nama,
				'email_seller', NEW.email,
				'jenis_seller', NEW.jenis,
				'jam_operasional_seller', NEW.jam_operasional,
				'punchline_seller', NEW.punchline,
				'deskripsi_seller', NEW.deskripsi,
				'seller_dedication', NEW.seller_dedication,
				'follower_total_seller', NEW.follower_total,
				'status_seller', NEW.status,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('seller_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_seller', OLD.id,
				'username_seller', OLD.username,
				'nama_seller', OLD.nama,
				'email_seller', OLD.email,
				'jenis_seller', OLD.jenis,
				'jam_operasional_seller', OLD.jam_operasional,
				'punchline_seller', OLD.punchline,
				'deskripsi_seller', OLD.deskripsi,
				'seller_dedication', OLD.seller_dedication,
				'follower_total_seller', OLD.follower_total,
				'status_seller', OLD.status,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('seller_channel', payload::text);
			RETURN OLD;
		END IF;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER seller_trigger
	AFTER INSERT OR UPDATE OR DELETE
	ON seller
	FOR EACH ROW
	EXECUTE FUNCTION notify_seller_change();
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk kurir
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func KurirDropper() string {
	return `
	DROP TRIGGER IF EXISTS kurir_trigger ON kurir;
	DROP FUNCTION IF EXISTS notify_kurir_change();
	`
}

func KurirTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_kurir_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
		changed_columns JSONB := '{}'::jsonb;
		column_change_name TEXT[] := ARRAY[]::TEXT[];
	BEGIN
		IF TG_OP = 'UPDATE' THEN
			-- Deteksi perubahan kolom
			IF OLD.username IS DISTINCT FROM NEW.username THEN
				changed_columns := jsonb_set(changed_columns, '{username}', to_jsonb(NEW.username));
				column_change_name := array_append(column_change_name, 'username');
			END IF;

			IF OLD.email IS DISTINCT FROM NEW.email THEN
				changed_columns := jsonb_set(changed_columns, '{email}', to_jsonb(NEW.email));
				column_change_name := array_append(column_change_name, 'email');
			END IF;

			IF OLD.jenis IS DISTINCT FROM NEW.jenis THEN
				changed_columns := jsonb_set(changed_columns, '{jenis}', to_jsonb(NEW.jenis));
				column_change_name := array_append(column_change_name, 'jenis');
			END IF;

			IF OLD.password_hash IS DISTINCT FROM NEW.password_hash THEN
				changed_columns := jsonb_set(changed_columns, '{password_hash}', to_jsonb(NEW.password_hash));
				column_change_name := array_append(column_change_name, 'password_hash');
			END IF;

			IF OLD.verified IS DISTINCT FROM NEW.verified THEN
				changed_columns := jsonb_set(changed_columns, '{verified}', to_jsonb(NEW.verified));
				column_change_name := array_append(column_change_name, 'verified');
			END IF;

			IF OLD.balance_kurir IS DISTINCT FROM NEW.balance_kurir THEN
				changed_columns := jsonb_set(changed_columns, '{balance_kurir}', to_jsonb(NEW.balance_kurir));
				column_change_name := array_append(column_change_name, 'balance_kurir');
			END IF;

			IF OLD.tipe_kendaraan IS DISTINCT FROM NEW.tipe_kendaraan THEN
				changed_columns := jsonb_set(changed_columns, '{tipe_kendaraan}', to_jsonb(NEW.tipe_kendaraan));
				column_change_name := array_append(column_change_name, 'tipe_kendaraan');
			END IF;

			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_kurir', OLD.id,
				'nama_kurir', OLD.nama,
				'username_kurir', OLD.username,
				'email_kurir', OLD.email,
				'jenis_kurir', OLD.jenis,
				'deskripsi_kurir', OLD.deskripsi,
				'status_kurir', OLD.status,
				'status_narik_kurir', OLD.status_narik,
				'verified_kurir', OLD.verified,
				'jumlah_pengiriman_kurir', OLD.jumlah_pengiriman,
				'balance_kurir', OLD.balance_kurir,
				'rating_kurir', OLD.rating,
				'jumlah_rating_kurir', OLD.jumlah_rating,
				'tipe_kendaraan_kurir', OLD.tipe_kendaraan,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', changed_columns,
				'column_change_name', column_change_name
			);
			PERFORM pg_notify('kurir_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_kurir', NEW.id,
				'nama_kurir', NEW.nama,
				'username_kurir', NEW.username,
				'email_kurir', NEW.email,
				'jenis_kurir', NEW.jenis,
				'deskripsi_kurir', NEW.deskripsi,
				'status_kurir', NEW.status,
				'status_narik_kurir', NEW.status_narik,
				'verified_kurir', NEW.verified,
				'jumlah_pengiriman_kurir', NEW.jumlah_pengiriman,
				'balance_kurir', NEW.balance_kurir,
				'rating_kurir', NEW.rating,
				'jumlah_rating_kurir', NEW.jumlah_rating,
				'tipe_kendaraan_kurir', NEW.tipe_kendaraan,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('kurir_channel', payload::text);
			RETURN NEW;

		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_kurir', OLD.id,
				'nama_kurir', OLD.nama,
				'username_kurir', OLD.username,
				'email_kurir', OLD.email,
				'jenis_kurir', OLD.jenis,
				'deskripsi_kurir', OLD.deskripsi,
				'status_kurir', OLD.status,
				'status_narik_kurir', OLD.status_narik,
				'verified_kurir', OLD.verified,
				'jumlah_pengiriman_kurir', OLD.jumlah_pengiriman,
				'balance_kurir', OLD.balance_kurir,
				'rating_kurir', OLD.rating,
				'jumlah_rating_kurir', OLD.jumlah_rating,
				'tipe_kendaraan_kurir', OLD.tipe_kendaraan,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('kurir_channel', payload::text);
			RETURN OLD;
		END IF;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER kurir_trigger
	AFTER INSERT OR UPDATE OR DELETE
	ON kurir
	FOR EACH ROW
	EXECUTE FUNCTION notify_kurir_change();
	`
}

func SetupEntityTriggers(db *gorm.DB) error {
	drops := []string{
		PenggunaDropper(),
		SellerDropper(),
		KurirDropper(),
	}

	for _, drop := range drops {
		if err := db.Exec(drop).Error; err != nil {
			fmt.Errorf("gagal hapus trigger/function: %w", err)
		}
	}

	triggers := [...]string{
		PenggunaTrigger(),
		SellerTrigger(),
		KurirTrigger(),
	}

	for _, trig := range triggers {
		if err := db.Exec(trig).Error; err != nil {
			fmt.Errorf("gagal buat trigger/function: %w", err)
		}
	}

	fmt.Println("Semua trigger untuk entity berhasil dibuat ulang!")
	return nil
}

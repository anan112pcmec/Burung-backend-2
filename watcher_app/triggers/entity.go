package trigger

import (
	"fmt"

	"gorm.io/gorm"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Fungsi Setup Entity Triggers
// Fungsi Ini melakukan setup trigger dan channel guna nanti akan di watch oleh watcher
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func SetupEntityTriggers(db *gorm.DB) error {
	drops := []string{
		// =================================================
		// Akan Drop Trigger dan Channel Yang sudah ada dan membuat yang baru
		// =================================================
		`DROP TRIGGER IF EXISTS pengguna_trigger ON pengguna;
		 DROP FUNCTION IF EXISTS notify_pengguna_change();`,

		`DROP TRIGGER IF EXISTS seller_trigger ON seller;
		 DROP FUNCTION IF EXISTS notify_seller_change();`,

		`DROP TRIGGER IF EXISTS kurir_trigger ON kurir;
		 DROP FUNCTION IF EXISTS notify_kurir_change();`,
	}

	for _, drop := range drops {
		if err := db.Exec(drop).Error; err != nil {
			return fmt.Errorf("gagal hapus trigger/function: %w", err)
		}
	}

	triggers := [...]string{
		// =================================================
		// Fungsi Untuk Pengguna
		// =================================================
		`
	CREATE OR REPLACE FUNCTION notify_pengguna_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
		changed_columns JSON := '{}';
		column_change_name TEXT[] := ARRAY[]::TEXT[];
	BEGIN
		IF TG_OP = 'UPDATE' THEN
			IF OLD.username IS DISTINCT FROM NEW.username THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{username}', to_jsonb(NEW.username));
				column_change_name := array_append(column_change_name, 'username');
			END IF;
			IF OLD.email IS DISTINCT FROM NEW.email THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{email}', to_jsonb(NEW.email));
				column_change_name := array_append(column_change_name, 'email');
			END IF;
			IF OLD.pin_hash IS DISTINCT FROM NEW.pin_hash THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{pin_hash}', to_jsonb(NEW.pin_hash));
				column_change_name := array_append(column_change_name, 'pin_hash');
			END IF;
			IF OLD.password_hash IS DISTINCT FROM NEW.password_hash THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{password_hash}', to_jsonb(NEW.password_hash));
				column_change_name := array_append(column_change_name, 'password_hash');
			END IF;
		END IF;

		payload := json_build_object(
			'table', TG_TABLE_NAME,
			'action', TG_OP,
			'id_user', NEW.id,
			'username_user', NEW.username,
			'nama_user', NEW.nama,
			'email_user', NEW.email,
			'changed_columns_pengguna', changed_columns,
			'column_change_name', column_change_name
		);

		PERFORM pg_notify('pengguna_channel', payload::text);
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER pengguna_trigger
	AFTER INSERT OR UPDATE OR DELETE
	ON pengguna
	FOR EACH ROW
	EXECUTE FUNCTION notify_pengguna_change();`,

		// =================================================
		// Fungsi Untuk Seller
		// =================================================

		`-- Fungsi untuk seller
	CREATE OR REPLACE FUNCTION notify_seller_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
		changed_columns JSON := '{}';
		column_change_name TEXT[] := ARRAY[]::TEXT[];
	BEGIN
		IF TG_OP = 'UPDATE' THEN
			IF OLD.username IS DISTINCT FROM NEW.username THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{username}', to_jsonb(NEW.username));
				column_change_name := array_append(column_change_name, 'username');
			END IF;
			IF OLD.email IS DISTINCT FROM NEW.email THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{email}', to_jsonb(NEW.email));
				column_change_name := array_append(column_change_name, 'email');
			END IF;
			IF OLD.jenis IS DISTINCT FROM NEW.jenis THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{jenis}', to_jsonb(NEW.jenis));
				column_change_name := array_append(column_change_name, 'jenis');
			END IF;
			IF OLD.norek IS DISTINCT FROM NEW.norek THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{norek}', to_jsonb(NEW.norek));
				column_change_name := array_append(column_change_name, 'norek');
			END IF;
			IF OLD.seller_dedication IS DISTINCT FROM NEW.seller_dedication THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{seller_dedication}', to_jsonb(NEW.seller_dedication));
				column_change_name := array_append(column_change_name, 'seller_dedication');
			END IF;
			IF OLD.jam_operasional IS DISTINCT FROM NEW.jam_operasional THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{jam_operasional}', to_jsonb(NEW.jam_operasional));
				column_change_name := array_append(column_change_name, 'jam_operasional');
			END IF;
			IF OLD.password_hash IS DISTINCT FROM NEW.password_hash THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{password_hash}', to_jsonb(NEW.password_hash));
				column_change_name := array_append(column_change_name, 'password_hash');
			END IF;
		END IF;

		payload := json_build_object(
			'table', TG_TABLE_NAME,
			'action', TG_OP,
			'id_seller', NEW.id,
			'username_seller', NEW.username,
			'nama_seller', NEW.nama,
			'email_seller', NEW.email,
			'jenis_seller', NEW.jenis,
			'norek_seller', NEW.norek,
			'seller_dedication', NEW.seller_dedication,
			'follower_total_seller', NEW.follower_total,
			'changed_columns_seller', changed_columns,
			'column_change_name', column_change_name
		);

		PERFORM pg_notify('seller_channel', payload::text);
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER seller_trigger
	AFTER INSERT OR UPDATE OR DELETE
	ON seller
	FOR EACH ROW
	EXECUTE FUNCTION notify_seller_change();`,

		// =================================================
		// Fungsi Untuk Kurir
		// =================================================

		`-- Fungsi untuk kurir
	CREATE OR REPLACE FUNCTION notify_kurir_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
		changed_columns JSON := '{}';
		column_change_name TEXT[] := ARRAY[]::TEXT[];
	BEGIN
		IF TG_OP = 'UPDATE' THEN
			IF OLD.username IS DISTINCT FROM NEW.username THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{username}', to_jsonb(NEW.username));
				column_change_name := array_append(column_change_name, 'username');
			END IF;
			IF OLD.email IS DISTINCT FROM NEW.email THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{email}', to_jsonb(NEW.email));
				column_change_name := array_append(column_change_name, 'email');
			END IF;
			IF OLD.jenis IS DISTINCT FROM NEW.jenis THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{jenis}', to_jsonb(NEW.jenis));
				column_change_name := array_append(column_change_name, 'jenis');
			END IF;
			IF OLD.password_hash IS DISTINCT FROM NEW.password_hash THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{password_hash}', to_jsonb(NEW.password_hash));
				column_change_name := array_append(column_change_name, 'password_hash');
			END IF;
			IF OLD.verified IS DISTINCT FROM NEW.verified THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{verified}', to_jsonb(NEW.verified));
				column_change_name := array_append(column_change_name, 'verified');
			END IF;
			IF OLD.balance_kurir IS DISTINCT FROM NEW.balance_kurir THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{balance_kurir}', to_jsonb(NEW.balance_kurir));
				column_change_name := array_append(column_change_name, 'balance_kurir');
			END IF;
			IF OLD.tipe_kendaraan IS DISTINCT FROM NEW.tipe_kendaraan THEN
				changed_columns := jsonb_set(changed_columns::jsonb, '{tipe_kendaraan}', to_jsonb(NEW.tipe_kendaraan));
				column_change_name := array_append(column_change_name, 'tipe_kendaraan');
			END IF;
		END IF;

		payload := json_build_object(
			'table', TG_TABLE_NAME,
			'action', TG_OP,
			'id_kurir', NEW.id,
			'nama_kurir', NEW.nama,
			'username_kurir', NEW.username,
			'email_kurir', NEW.email,
			'jenis_kurir', NEW.jenis,
			'verified_kurir', NEW.verified,
			'status_narik_kurir', NEW.status_narik,
			'balance_kurir', NEW.balance_kurir,
			'tipe_kendaraan_kurir', NEW.tipe_kendaraan,
			'changed_columns_kurir', changed_columns,
			'column_change_name', column_change_name
		);

		PERFORM pg_notify('kurir_channel', payload::text);
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER kurir_trigger
	AFTER INSERT OR UPDATE OR DELETE
	ON kurir
	FOR EACH ROW
	EXECUTE FUNCTION notify_kurir_change();`,
	}

	for _, trig := range triggers {
		if err := db.Exec(trig).Error; err != nil {
			return fmt.Errorf("gagal buat trigger/function: %w", err)
		}
	}

	fmt.Println("Semua trigger untuk entity berhasil dibuat ulang!")
	return nil
}

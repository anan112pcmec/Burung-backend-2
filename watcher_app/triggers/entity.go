package trigger

import (
	"fmt"

	"gorm.io/gorm"

)

func SetupEntityTriggers(db *gorm.DB) error {
	drops := []string{
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
		// Pengguna (sudah lengkap)
		`-- Fungsi untuk pengguna
		CREATE OR REPLACE FUNCTION notify_pengguna_change()
		RETURNS trigger AS $$
		DECLARE
			payload JSON;
			changed_columns JSON := '{}';
		BEGIN
    -- Cek update untuk kolom tertentu
    IF TG_OP = 'UPDATE' THEN
        IF OLD.username IS DISTINCT FROM NEW.username THEN
            changed_columns := jsonb_set(changed_columns::jsonb, '{username}', to_jsonb(NEW.username));
        END IF;
        IF OLD.nama IS DISTINCT FROM NEW.nama THEN
            changed_columns := jsonb_set(changed_columns::jsonb, '{nama}', to_jsonb(NEW.nama));
        END IF;
        IF OLD.email IS DISTINCT FROM NEW.email THEN
            changed_columns := jsonb_set(changed_columns::jsonb, '{email}', to_jsonb(NEW.email));
        END IF;
        IF OLD.password_hash IS DISTINCT FROM NEW.password_hash THEN
            changed_columns := jsonb_set(changed_columns::jsonb, '{password_hash}', to_jsonb(NEW.password_hash));
        END IF;
        IF OLD.pin_hash IS DISTINCT FROM NEW.pin_hash THEN
            changed_columns := jsonb_set(changed_columns::jsonb, '{pin_hash}', to_jsonb(NEW.pin_hash));
        END IF;
        IF OLD.status IS DISTINCT FROM NEW.status THEN
            changed_columns := jsonb_set(changed_columns::jsonb, '{status}', to_jsonb(NEW.status));
        END IF;
    END IF;

    payload := json_build_object(
        'table', TG_TABLE_NAME,
        'action', TG_OP,
        'id_user', NEW.id,
        'username_user', NEW.username,
        'nama_user', NEW.nama,
        'email_user', NEW.email,
        'changed_columns_pengguna', changed_columns
    );

    PERFORM pg_notify('pengguna_channel', payload::text);
    RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;

		CREATE TRIGGER pengguna_trigger
		AFTER INSERT OR UPDATE OR DELETE
		ON pengguna
		FOR EACH ROW
		EXECUTE FUNCTION notify_pengguna_change();

		-- Fungsi untuk seller dengan changed_columns_seller
		CREATE OR REPLACE FUNCTION notify_seller_change()
		RETURNS trigger AS $$
		DECLARE
			payload JSON;
			changed_columns JSON := '{}';
		BEGIN
			-- Cek update untuk kolom tertentu, misalnya status
			IF TG_OP = 'UPDATE' THEN
				IF OLD.status IS DISTINCT FROM NEW.status THEN
					changed_columns := jsonb_set(changed_columns::jsonb, '{status}', to_jsonb(NEW.status));
				END IF;
				-- Bisa ditambahkan kolom lain jika ingin track perubahan
			END IF;

			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_seller', NEW.id,
				'username_seller', NEW.username,
				'nama_seller', NEW.nama,
				'email_seller', NEW.email,
				'jenis_seller', NEW.jenis,
				'seller_dedication', NEW.seller_dedication,
				'follower_total_seller', NEW.follower_total,
				'changed_columns_seller', changed_columns
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

		// Kurir (tidak ada username)
		`CREATE OR REPLACE FUNCTION notify_kurir_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
	BEGIN
		payload := json_build_object(
			'table', TG_TABLE_NAME,
			'action', TG_OP,
			'id', NEW.id,
			'nama', NEW.nama,
			'email', NEW.email,
			'no_hp', NEW.no_hp
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

	fmt.Println("Semua trigger berhasil dibuat ulang!")
	return nil
}

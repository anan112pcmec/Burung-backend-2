package trigger

import (
	"fmt"

	"gorm.io/gorm"
)

func SetupFollowerTriggers(db *gorm.DB) error {
	drops := []string{
		`DROP TRIGGER IF EXISTS follower_trigger ON follower;
		 DROP FUNCTION IF EXISTS notify_follower_change();`,
	}

	for _, drop := range drops {
		if err := db.Exec(drop).Error; err != nil {
			return fmt.Errorf("gagal hapus trigger/function: %w", err)
		}
	}

	triggerfollower := [...]string{
		`
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
		`,
	}

	for _, trig := range triggerfollower {
		if err := db.Exec(trig).Error; err != nil {
			return fmt.Errorf("gagal buat trigger/function: %w", err)
		}
	}

	fmt.Println("Semua trigger berhasil dibuat ulang!")
	return nil
}

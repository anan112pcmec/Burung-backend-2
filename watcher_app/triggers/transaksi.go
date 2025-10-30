package trigger

import (
	"fmt"

	"gorm.io/gorm"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk transaksi
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TransaksiDropper() string {
	return `
	DROP TRIGGER IF EXISTS transaksi_trigger ON transaksi;
	DROP FUNCTION IF EXISTS notify_transaksi_change();
	`
}

func TransaksiTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_transaksi_change()
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
			-- Deteksi perubahan setiap kolom penting
			IF OLD.status IS DISTINCT FROM NEW.status THEN
				changed_columns := jsonb_set(changed_columns, '{status}', to_jsonb(NEW.status));
				column_change_name := array_append(column_change_name, 'status');
			END IF;

			IF OLD.metode IS DISTINCT FROM NEW.metode THEN
				changed_columns := jsonb_set(changed_columns, '{metode}', to_jsonb(NEW.metode));
				column_change_name := array_append(column_change_name, 'metode');
			END IF;

			IF OLD.catatan IS DISTINCT FROM NEW.catatan THEN
				changed_columns := jsonb_set(changed_columns, '{catatan}', to_jsonb(NEW.catatan));
				column_change_name := array_append(column_change_name, 'catatan');
			END IF;

			IF OLD.kuantitas_barang IS DISTINCT FROM NEW.kuantitas_barang THEN
				changed_columns := jsonb_set(changed_columns, '{kuantitas_barang}', to_jsonb(NEW.kuantitas_barang));
				column_change_name := array_append(column_change_name, 'kuantitas_barang');
			END IF;

			IF OLD.total IS DISTINCT FROM NEW.total THEN
				changed_columns := jsonb_set(changed_columns, '{total}', to_jsonb(NEW.total));
				column_change_name := array_append(column_change_name, 'total');
			END IF;

			-- Payload untuk update
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_transaksi', OLD.id,
				'id_pengguna_transaksi', OLD.id_pengguna,
				'id_seller_transaksi', OLD.id_seller,
				'id_barang_induk_transaksi', OLD.id_barang_induk,
				'id_alamat_pengguna_transaksi', OLD.id_alamat_pengguna_transaksi,
				'id_pembayaran_transaksi', OLD.id_pembayaran,
				'kode_order_transaksi', OLD.kode_order,
				'jenis_pengiriman_transaksi', OLD.jenis_pengiriman,
				'status_transaksi', OLD.status,
				'metode_transaksi', OLD.metode,
				'catatan_transaksi', OLD.catatan,
				'kuantitas_barang_transaksi', OLD.kuantitas_barang,
				'total_transaksi', OLD.total,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,
				'changed_columns', changed_columns,
				'column_change_name', column_change_name
			);
			PERFORM pg_notify('transaksi_channel', payload::text);
			RETURN NEW;

		-- =====================================================================
		-- HANDLE INSERT
		-- =====================================================================
		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_transaksi', NEW.id,
				'id_pengguna_transaksi', NEW.id_pengguna,
				'id_seller_transaksi', NEW.id_seller,
				'id_barang_induk_transaksi', NEW.id_barang_induk,
				'id_alamat_pengguna_transaksi', NEW.id_alamat_pengguna_transaksi,
				'id_pembayaran_transaksi', NEW.id_pembayaran,
				'kode_order_transaksi', NEW.kode_order,
				'jenis_pengiriman_transaksi', NEW.jenis_pengiriman,
				'status_transaksi', NEW.status,
				'metode_transaksi', NEW.metode,
				'catatan_transaksi', NEW.catatan,
				'kuantitas_barang_transaksi', NEW.kuantitas_barang,
				'total_transaksi', NEW.total,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('transaksi_channel', payload::text);
			RETURN NEW;

		-- =====================================================================
		-- HANDLE DELETE
		-- =====================================================================
		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_transaksi', OLD.id,
				'id_pengguna_transaksi', OLD.id_pengguna,
				'id_seller_transaksi', OLD.id_seller,
				'id_barang_induk_transaksi', OLD.id_barang_induk,
				'id_alamat_pengguna_transaksi', OLD.id_alamat_pengguna_transaksi,
				'id_pembayaran_transaksi', OLD.id_pembayaran,
				'kode_order_transaksi', OLD.kode_order,
				'jenis_pengiriman_transaksi', OLD.jenis_pengiriman,
				'status_transaksi', OLD.status,
				'metode_transaksi', OLD.metode,
				'catatan_transaksi', OLD.catatan,
				'kuantitas_barang_transaksi', OLD.kuantitas_barang,
				'total_transaksi', OLD.total,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('transaksi_channel', payload::text);
			RETURN OLD;
		END IF;

		-- Default safety return
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER transaksi_trigger
	AFTER INSERT OR UPDATE OR DELETE
	ON transaksi
	FOR EACH ROW
	EXECUTE FUNCTION notify_transaksi_change();
	`
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Dropper Dan Trigger untuk pembayaran
// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func PembayaranDropper() string {
	return `
	DROP TRIGGER IF EXISTS pembayaran_trigger ON pembayaran;
	DROP FUNCTION IF EXISTS notify_pembayaran_change();
	`
}

func PembayaranTrigger() string {
	return `
	CREATE OR REPLACE FUNCTION notify_pembayaran_change()
	RETURNS trigger AS $$
	DECLARE
		payload JSON;
		changed_columns JSONB := '{}'::jsonb;
		column_change_name TEXT[] := ARRAY[]::TEXT[];
	BEGIN
		IF TG_OP = 'UPDATE' THEN
			-- deteksi kolom yang berubah
			IF NEW.kode_transaksi IS DISTINCT FROM OLD.kode_transaksi THEN
				changed_columns := changed_columns || jsonb_build_object('kode_transaksi_pembayaran', jsonb_build_object('old', OLD.kode_transaksi, 'new', NEW.kode_transaksi));
				column_change_name := array_append(column_change_name, 'kode_transaksi_pembayaran');
			END IF;
			IF NEW.kode_order IS DISTINCT FROM OLD.kode_order THEN
				changed_columns := changed_columns || jsonb_build_object('kode_order_pembayaran', jsonb_build_object('old', OLD.kode_order, 'new', NEW.kode_order));
				column_change_name := array_append(column_change_name, 'kode_order_pembayaran');
			END IF;
			IF NEW.provider IS DISTINCT FROM OLD.provider THEN
				changed_columns := changed_columns || jsonb_build_object('provider_pembayaran', jsonb_build_object('old', OLD.provider, 'new', NEW.provider));
				column_change_name := array_append(column_change_name, 'provider_pembayaran');
			END IF;
			IF NEW.amount IS DISTINCT FROM OLD.amount THEN
				changed_columns := changed_columns || jsonb_build_object('amount_pembayaran', jsonb_build_object('old', OLD.amount, 'new', NEW.amount));
				column_change_name := array_append(column_change_name, 'amount_pembayaran');
			END IF;
			IF NEW.payment_type IS DISTINCT FROM OLD.payment_type THEN
				changed_columns := changed_columns || jsonb_build_object('payment_type_pembayaran', jsonb_build_object('old', OLD.payment_type, 'new', NEW.payment_type));
				column_change_name := array_append(column_change_name, 'payment_type_pembayaran');
			END IF;
			IF NEW.paid_at IS DISTINCT FROM OLD.paid_at THEN
				changed_columns := changed_columns || jsonb_build_object('paid_at_pembayaran', jsonb_build_object('old', OLD.paid_at, 'new', NEW.paid_at));
				column_change_name := array_append(column_change_name, 'paid_at_pembayaran');
			END IF;

			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_pembayaran', OLD.id,
				'kode_transaksi_pembayaran', OLD.kode_transaksi,
				'kode_order_pembayaran', OLD.kode_order,
				'provider_pembayaran', OLD.provider,
				'amount_pembayaran', OLD.amount,
				'payment_type_pembayaran', OLD.payment_type,
				'paid_at_pembayaran', OLD.paid_at,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,
				'changed_columns', changed_columns,
				'column_change_name', column_change_name
			);
			PERFORM pg_notify('pembayaran_channel', payload::text);

		ELSIF TG_OP = 'INSERT' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_pembayaran', NEW.id,
				'kode_transaksi_pembayaran', NEW.kode_transaksi,
				'kode_order_pembayaran', NEW.kode_order,
				'provider_pembayaran', NEW.provider,
				'amount_pembayaran', NEW.amount,
				'payment_type_pembayaran', NEW.payment_type,
				'paid_at_pembayaran', NEW.paid_at,
				'created_at', NEW.created_at,
				'updated_at', NEW.updated_at,
				'deleted_at', NEW.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('pembayaran_channel', payload::text);

		ELSIF TG_OP = 'DELETE' THEN
			payload := json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'id_pembayaran', OLD.id,
				'kode_transaksi_pembayaran', OLD.kode_transaksi,
				'kode_order_pembayaran', OLD.kode_order,
				'provider_pembayaran', OLD.provider,
				'amount_pembayaran', OLD.amount,
				'payment_type_pembayaran', OLD.payment_type,
				'paid_at_pembayaran', OLD.paid_at,
				'created_at', OLD.created_at,
				'updated_at', OLD.updated_at,
				'deleted_at', OLD.deleted_at,
				'changed_columns', '{}'::jsonb,
				'column_change_name', ARRAY[]::TEXT[]
			);
			PERFORM pg_notify('pembayaran_channel', payload::text);
		END IF;

		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER pembayaran_trigger
	AFTER INSERT OR UPDATE OR DELETE ON pembayaran
	FOR EACH ROW
	EXECUTE FUNCTION notify_pembayaran_change();
	`
}

func SetupTransaksiTriggers(db *gorm.DB) error {
	drops := []string{
		TransaksiDropper(),
	}

	for _, drop := range drops {
		if err := db.Exec(drop).Error; err != nil {
			return fmt.Errorf("gagal hapus trigger/function: %w", err)
		}
	}

	triggertransaksi := [...]string{
		TransaksiTrigger(),
	}

	for _, trig := range triggertransaksi {
		if err := db.Exec(trig).Error; err != nil {
			return fmt.Errorf("gagal buat trigger/function: %w", err)
		}
	}

	fmt.Println("Semua trigger berhasil dibuat ulang!")
	return nil
}

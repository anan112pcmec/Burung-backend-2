package notify_payload

import "github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"

type NotifyResponsePayloadTransaksi struct {
	TableAndAction
	models.Transaksi
	ColumnChangeInfo
}

type NotifyResponsePayloadPembayaran struct {
	TableAndAction
	models.Pembayaran
	ColumnChangeInfo
}

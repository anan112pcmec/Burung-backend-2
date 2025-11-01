package notify_payload

import "github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"

type NotifyResponsePayloadPengiriman struct {
	TableAndAction
	models.Pengiriman
	ColumnChangeInfo
}

type NotifyResponsePayloadJejakPengiriman struct {
	TableAndAction
	models.JejakPengiriman
	ColumnChangeInfo
}

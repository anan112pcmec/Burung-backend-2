package notify_payload

import "github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"

type NotifyResponsePayloadFollower struct {
	TableAndAction
	models.Follower
	ColumnChangeInfo
}

type NotifyPayloadResponseInformasiKendaraanKurir struct {
	TableAndAction
	models.InformasiKendaraanKurir
	ColumnChangeInfo
}

package notify_payload

import "github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"

type TableAndAction struct {
	Table  string `json:"table"`
	Action string `json:"action"`
}

type ChangedColumns struct {
	Status string `json:"status"`
}

// ////////////////////////////////////////////////////////////////////////////
// BARANG PAYLOAD
// ////////////////////////////////////////////////////////////////////////////

type NotifyResponsesPayloadBarang struct {
	TableAndAction
	models.BarangInduk
	ChangedColumns ChangedColumns `json:"changed_columns"`
}

type NotifyResponsesPayloadKategoriBarang struct {
	TableAndAction
	models.KategoriBarang
	ChangedColumns ChangedColumns `json:"changed_columns"`
}

type NotifyResponsePayloadVarianBarang struct {
	TableAndAction
	OldData models.VarianBarang `json:"old_data"`
	NewData models.VarianBarang `json:"new_data"`
}

// ////////////////////////////////////////////////////////////////////////////
// ENGAGEMENT PAYLOAD
// ////////////////////////////////////////////////////////////////////////////

type NotifyResponsePayloadKomentar struct {
	TableAndAction
	models.Komentar
}

// ////////////////////////////////////////////////////////////////////////////
// TRANSAKSI PAYLOAD
// ////////////////////////////////////////////////////////////////////////////

type NotifyResponseTransaksi struct {
	TableAndAction
	models.Transaksi
}

// ////////////////////////////////////////////////////////////////////////////
// INFORMASI KURIR PAYLOAD
// ////////////////////////////////////////////////////////////////////////////

type NotifyResponseInformasiKurir struct {
	TableAndAction
	IdKurir         int64  `json:"informasi_id_kurir"`
	StatusPerizinan string `json:"status_perizinan_kurir"`
	JenisKendaraan  string `json:"jenis_kendaraan"`
}

type NotifyResponsePengiriman struct {
	TableAndAction
	models.Pengiriman
}

// ////////////////////////////////////////////////////////////////////////////
// SOCIAL MEDIA PAYLOAD
// ////////////////////////////////////////////////////////////////////////////

type NotifyResponseFollower struct {
	TableAndAction
	models.Follower
}

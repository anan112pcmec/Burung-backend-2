package notify_payload

import "github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"

type TableAndAction struct {
	Table  string `json:"table"`
	Action string `json:"action"`
}

type SellerNotifyPayload struct {
	TableAndAction
	Id               int32  `json:"id_seller"`
	Nama             string `json:"nama_seller"`
	Email            string `json:"email_seller"`
	Jenis            string `json:"jenis_seller"`
	SellerDedication string `json:"seller_dedication"`
	FollowerTotal    string `json:"follower_total_seller"`
}

func (s *SellerNotifyPayload) Validate() (Message string, Action string) {
	if s.Table != "seller" {
		return "Bukan Dari Table Seller", s.Action
	}

	if s.Id == 0 || s.Nama == "" || s.Email == "" {
		return "Payload Ada Yang Kosong ID/NAMA/EMAIL", s.Action
	}

	return "Memenuhi_Syarat", s.Action
}

type KurirNotifyPayload struct {
	TableAndAction
	Id    int64  `json:"id"`
	Nama  string `json:"nama"`
	Email string `json:"email"`
}

func (k *KurirNotifyPayload) Validate() (Message string, Action string) {
	if k.Table != "kurir" {
		return "Bukan Dari Table Kurir", k.Action
	}

	if k.Id == 0 || k.Nama == "" || k.Email == "" {
		return "Payload Ada Yang Kosong ID/NAMA/EMAIL/NOHP", k.Action
	}

	return "Memenuhi_Syarat", k.Action
}

type ChangedColumns struct {
	Status string `json:"status"`
}

// ////////////////////////////////////////////////////////////////////////////
// ENTITY PAYLOAD
// ////////////////////////////////////////////////////////////////////////////

type NotifyResponsesPayloadPengguna struct {
	TableAndAction
	models.Pengguna
	ChangedColumns ChangedColumns `json:"changed_columns_pengguna"`
}

type NotifyResponsePayloadSeller struct {
	TableAndAction
	models.Seller
	ChangedColumns ChangedColumns `json:"changed_columns_seller"`
}

// ////////////////////////////////////////////////////////////////////////////
// BARANG PAYLOAD
// ////////////////////////////////////////////////////////////////////////////

type NotifyResponsesPayloadBarang struct {
	TableAndAction
	models.BarangInduk
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

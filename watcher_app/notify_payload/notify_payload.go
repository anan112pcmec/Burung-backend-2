package notify_payload

// --- Pengguna ---
type PenggunaNotifyPayload struct {
	Table    string `json:"table_pengguna"`
	Action   string `json:"action_pengguna"`
	Id       int64  `json:"id_pengguna"`
	Username string `json:"username_pengguna"`
	Nama     string `json:"nama_pengguna"`
	Email    string `json:"email_pengguna"`
}

func (p *PenggunaNotifyPayload) Validate() (Message string, Action string) {
	if p.Table != "pengguna" {
		return "Bukan Dari Table Pengguna", p.Action
	}

	if p.Id == 0 || p.Email == "" || p.Username == "" || p.Nama == "" {
		return "Payload Ada Yang Kosong ID/NAMA/USERNAME/EMAIL", p.Action
	}

	return "Memenuhi_Syarat", p.Action
}

type SellerNotifyPayload struct {
	Table            string `json:"table_seller"`
	Action           string `json:"action_seller"`
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
	Table  string `json:"table"`
	Action string `json:"action"`
	Id     int64  `json:"id"`
	Nama   string `json:"nama"`
	Email  string `json:"email"`
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

type NotifyResponsesPayloadPengguna struct {
	PenggunaNotifyPayload
	ChangedColumns ChangedColumns `json:"changed_columns"`
}

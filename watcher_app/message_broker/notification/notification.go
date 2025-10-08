package notification

import "time"

const (
	LevelCritical = "critical"
	LevelMedium   = "medium"
	LevelLow      = "low"
)

const (
	User   = "user"
	Seller = "seller"
	Kurir  = "kurir"
)

type Notification struct {
	ID        string    `json:"id_notification"`        // UUID untuk tracking
	Level     string    `json:"level_notification"`     // critical | medium | low
	Entity    string    `json:"entity_notification"`    // user | seller | courier | system
	Type      string    `json:"type_notification"`      // transaction | security | social | marketing | system
	EventName string    `json:"event_notification"`     // ex: order.new, payment.success
	Message   string    `json:"pesan_notification"`     // pesan human-readable
	Timestamp time.Time `json:"timestamp_notification"` // waktu dibuat
	Metadata  any       `json:"metadata"`               // payload fleksibel
}

// Penjelasan Notification Struct
// ID -> Tracking Message
// Level -> Critical Level Notification (Critical / Medium / Low) memebri tahu seberapa penting sebuah notif
// Entity -> Jenis Entity Penerima (pengguna / seller / kurir)
// Type -> Jenis Misal Transaksi / Account dll
// EventName -> irisan lebih bawah lagi dari Type misal membayar barang
// Timestamp -> data waktu
// Metadata -> file bisa pdf dll jika dibutuhkan

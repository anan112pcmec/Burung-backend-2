package notification

import "time"

const (
	KurirAccount  = "account"
	KurirDelivery = "delivery"
)

// //////////////////////////////////////////
// Mendefinisikan Semua Method Critical Kurir
// //////////////////////////////////////////

func (n *Notification) KurirAccount(Aksi, Pesan string, File any) {
	n.Level = LevelCritical
	n.Entity = Kurir
	n.Type = KurirAccount
	n.EventName = Aksi
	n.Message = Pesan
	n.Timestamp = time.Now()
	if File != nil {
		n.Metadata = File
	}
}

func (n *Notification) KurirDelivery(Aksi, Pesan string, File any) {
	n.Level = LevelCritical
	n.Entity = Kurir
	n.Type = KurirDelivery
	n.EventName = Aksi
	n.Message = Pesan
	n.Timestamp = time.Now()
	if File != nil {
		n.Metadata = File
	}
}

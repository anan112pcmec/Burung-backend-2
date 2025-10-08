package notification

import "time"

// //////////////////
// Static Type message
// /////////////////

const (
	UserTransaction = "transaksi"
	UserAccount     = "account"
)

// //////////////////////////////////////////
// Mendefinisikan Semua Method Critical User
// //////////////////////////////////////////

func (n *Notification) UserTransaksi(Aksi, Pesan string, File any) {
	n.Level = LevelCritical
	n.Entity = User
	n.Type = UserTransaction
	n.EventName = Aksi
	n.Message = Pesan
	n.Timestamp = time.Now()
	if File != nil {
		n.Metadata = File
	}
}

func (n *Notification) UserAccount(Aksi, Pesan string, File any) {
	n.Level = LevelCritical
	n.Entity = User
	n.Type = UserAccount
	n.EventName = Aksi
	n.Message = Pesan
	n.Timestamp = time.Now()
	if File != nil {
		n.Metadata = File
	}
}

package notification

import "time"

const (
	SellerProduct = "product"
	SellerAccount = "account"
	SellerOrder   = "order"
)

// //////////////////////////////////////////
// Mendefinisikan Semua Method Critical Seller
// //////////////////////////////////////////

func (n *Notification) SellerProduct(Aksi, Pesan string, File any) {
	n.Level = LevelCritical
	n.Entity = Seller
	n.Type = SellerProduct
	n.EventName = Aksi
	n.Message = Pesan
	n.Timestamp = time.Now()
	if File != nil {
		n.Metadata = File
	}
}

func (n *Notification) SellerAccount(Aksi, Pesan string, File any) {
	n.Level = LevelCritical
	n.Entity = Seller
	n.Type = SellerAccount
	n.EventName = Aksi
	n.Message = Pesan
	n.Timestamp = time.Now()
	if File != nil {
		n.Metadata = File
	}
}

func (n *Notification) SellerOrder(Aksi, Pesan string, File any) {
	n.Level = LevelCritical
	n.Entity = Seller
	n.Type = SellerOrder
	n.EventName = Aksi
	n.Message = Pesan
	n.Timestamp = time.Now()
	if File != nil {
		n.Metadata = File
	}
}

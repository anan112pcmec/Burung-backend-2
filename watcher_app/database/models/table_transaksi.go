package models

import "time"

type StatusTransaksi string

const (
	Dibayar             StatusTransaksi = "Dibayar"
	DiprosesTransaction StatusTransaksi = "Diproses"
	Waiting             StatusTransaksi = "Waiting"
	Dikirim             StatusTransaksi = "Dikirim"
	Selesai             StatusTransaksi = "Selesai"
	Dibatalkan          StatusTransaksi = "Dibatalkan"
)

type MetodePembayaran string

const (
	Transfer_Bank MetodePembayaran = "Transfer Bank"
	Kartu_Kredit  MetodePembayaran = "Kartu Kredit"
	E_Wallet      MetodePembayaran = "E Wallet"
	COD           MetodePembayaran = "COD"
)

type Ongkir int16

const (
	Via_Reguler Ongkir = 13000
	Via_Express Ongkir = 17000
	Via_Ekonomi Ongkir = 10000
	Via_Sameday Ongkir = 31000
	Via_Nextday Ongkir = 25000
	Via_Cargo   Ongkir = 7000
)

type Transaksi struct {
	ID             int64             `gorm:"primaryKey;autoIncrement" json:"id_transaksi"`
	IdPengguna     int64             `gorm:"column:id_pengguna;not null" json:"id_pengguna_transaksi"`
	Pengguna       Pengguna          `gorm:"foreignKey:IdPengguna;references:ID"`
	IdSeller       int32             `gorm:"column:id_seller;not null" json:"id_seller_transaksi"`
	Seller         Seller            `gorm:"foreignKey:IdSeller;references:ID"`
	IdBarangInduk  int64             `gorm:"column:id_barang_induk;not null" json:"id_barang_induk_transaksi"`
	BarangInduk    BarangInduk       `gorm:"foreignKey:IdBarangInduk;references:ID"`
	KodeTransaksi  string            `gorm:"column:kode_transaksi;type:varchar(250);not null" json:"kode_transaksi"`
	Status         StatusTransaksi   `gorm:"column:status;type:varchar(250);not null;default:'Dibayar'" json:"status_transaksi"`
	Metode         MetodePembayaran  `gorm:"column:metode;type:varchar(250);not null;default: 'Transfer Bank'" json:"metode_transaksi"`
	AlamatPengirim string            `gorm:"column:alamat_pengiriman;type:text;not null;default:''" json:"alamat_pengiriman_transaksi"`
	Catatan        string            `gorm:"column:catatan;type:text" json:"catatan_transaksi"`
	Jumlah         int16             `gorm:"column:jumlah;type:int2;not null;default:0" json:"jumlah_barang_transaksi"`
	Layanan        JenisLayananKurir `gorm:"column:layanan;type:varchar(250);not null;default:'Reguler'" json:"layanan_pengiriman_transaksi"`
	Ongkir         Ongkir            `gorm:"column:ongkir;type:int2;not null;default:13000" json:"ongkir_transaksi"`
	Total          int32             `gorm:"column:total;type:int4; not null;default:0" json:"total_transaksi"`
	CreatedAt      time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      *time.Time        `gorm:"index" json:"deleted_at,omitempty"`
}

func (Transaksi) TableName() string {
	return "transaksi"
}

type Pembayaran struct {
	ID          int64      `gorm:"primaryKey;autoIncrement" json:"id_pembayaran"`
	IdTransaksi int64      `gorm:"column:id_transaksi;not null" json:"id_transaksi_pembayaran"`
	Provider    string     `gorm:"column:provider;type:text;not null;default:''" json:"provider_pembayaran"`
	Amount      int32      `gorm:"column:amount;type:int4;not null,default:0" json:"amount_pembayaran"`
	PaidAt      string     `gorm:"column:paid_at;type:text;not null;default:''" json:"paid_at_pembayaran"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (Pembayaran) TableName() string {
	return "pembayaran"
}

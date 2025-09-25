package models

import (
	"time"
)

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
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id_transaksi"`
	IdPengguna    int64          `gorm:"column:id_pengguna;not null" json:"id_pengguna_transaksi"`
	Pengguna      Pengguna       `gorm:"foreignKey:IdPengguna;references:ID" json:"-"`
	IdSeller      int32          `gorm:"column:id_seller;not null" json:"id_seller_transaksi"`
	Seller        Seller         `gorm:"foreignKey:IdSeller;references:ID" json:"-"`
	IdBarangInduk int32          `gorm:"column:id_barang_induk;not null" json:"id_barang_induk_transaksi"`
	BarangInduk   BarangInduk    `gorm:"foreignKey:IdBarangInduk;references:ID" json:"-"`
	IdAlamat      int64          `gorm:"column:id_alamat_pengguna_transaksi" json:"id_alamat_pengguna_transaksi"`
	Alamat        AlamatPengguna `gorm:"foreignKey:IdAlamat;references:ID" json:"-"`
	IdPembayaran  int64          `gorm:"column:id_pembayaran;not null" json:"id_pembayaran_transaksi"`
	IPembayaran   Pembayaran     `gorm:"foreignKey:IdPembayaran;references:ID" json:"-"`
	KodeOrder     string         `gorm:"column:kode_order;not null" json:"kode_order_transaksi"`
	KPembayaran   Pembayaran     `gorm:"foreignKey:KodeOrder;references:KodeOrderTransaksi" json:"-"`
	Status        string         `gorm:"column:status;type:status_transaksi;not null;default:'Dibayar'" json:"status_transaksi"`
	Metode        string         `gorm:"column:metode;type:varchar(50);not null;" json:"metode_transaksi"`
	Catatan       string         `gorm:"column:catatan;type:text" json:"catatan_transaksi"`
	Kuantitas     int16          `gorm:"column:kuantitas_barang;type:int2;not null" json:"kuantitas_barang_transaksi"` // Isinya JumlahBarang
	Total         int32          `gorm:"column:total;type:int4; not null;default:0" json:"total_transaksi"`            // Isinyaduit
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     *time.Time     `gorm:"index" json:"deleted_at,omitempty"`
}

func (Transaksi) TableName() string {
	return "transaksi"
}

type Pembayaran struct {
	ID                 int64      `gorm:"primaryKey;autoIncrement" json:"id_pembayaran"`
	KodeTransaksi      string     `gorm:"column:kode_transaksi;not null" json:"kode_transaksi_pembayaran"`
	KodeOrderTransaksi string     `gorm:"column:kode_order;type:varchar(250);unique;not null" json:"kode_order_pembayaran"`
	Provider           string     `gorm:"column:provider;type:text;not null;default:''" json:"provider_pembayaran"`
	Amount             int32      `gorm:"column:amount;type:int4;not null,default:0" json:"amount_pembayaran"`
	PaymentType        string     `gorm:"column:payment_type;type:varchar(120);not null" json:"payment_type_pembayaran"`
	PaidAt             string     `gorm:"column:paid_at;type:text;not null;default:''" json:"paid_at_pembayaran"`
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (Pembayaran) TableName() string {
	return "pembayaran"
}

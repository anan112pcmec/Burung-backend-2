package models

import "time"

type Status string

const (
	StatusOnline  Status = "Online"
	StatusOffline Status = "Offline"
)

type Pengguna struct {
	ID             int64      `gorm:"primaryKey;autoIncrement" json:"id_user"`
	Username       string     `gorm:"column:username;type:varchar(100);not null;default:''" json:"username_user"`
	Nama           string     `gorm:"column:nama;type:text;not null;default:''" json:"nama_user"`
	Email          string     `gorm:"column:email;type:varchar(100);not null;uniqueIndex" json:"email_user"`
	PasswordHash   string     `gorm:"column:password_hash;type:varchar(250);not null;default:''" json:"pass_user"`
	PinHash        string     `gorm:"column:pin_hash;type:varchar(250);not null;default:''" json:"pin_user"`
	StatusPengguna Status     `gorm:"column:status;type:varchar(250);not null;default:'Offline'" json:"status_user"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (Pengguna) TableName() string {
	return "pengguna"
}

type JenisSeller string

const (
	Brands      JenisSeller = "Brands"
	Distributor JenisSeller = "Distributors"
	Personal    JenisSeller = "Personal"
)

type SellerType string

const (
	Pakaian_Fashion           SellerType = "Pakaian & Fashion"
	Kosmetik_Kecantikan       SellerType = "Kosmetik & Kecantikan"
	Elektronik_Gadget         SellerType = "Elektronik & Gadget"
	Buku_Media                SellerType = "Buku & Media"
	Makanan_Minuman           SellerType = "Makanan & Minuman"
	Kesehatan_Obat            SellerType = "Kesehatan & Obat"
	Ibu_Bayi                  SellerType = "Ibu & Bayi"
	Mainan_Hobi               SellerType = "Mainan & Hobi"
	Olahraga_Outdoor          SellerType = "Olahraga & Outdoor"
	Otomotis_SparePart        SellerType = "Otomotif & Sparepart"
	Rumah_Tangga_Perabotan    SellerType = "Rumah Tangga & Perabotan"
	AlatTulis_Kantor          SellerType = "Alat Tulis & Kantor"
	Perhiasan_Aksesoris       SellerType = "Perhiasan & Aksesoris"
	Produk_Digital            SellerType = "Produk Digital"
	BahanBangunan_Perkakas    SellerType = "Bahan Bangunan & Perkakas"
	ProdukPertanian_Perikanan SellerType = "Produk Pertanian & Perikanan"
	Musik_Instrumen           SellerType = "Musik & Instrumen"
	Film_Koleksi              SellerType = "Film & Koleksi"
	Semua_Barang              SellerType = "Semua Barang"
)

type Seller struct {
	ID               int32       `gorm:"primaryKey;autoIncrement" json:"id_seller"`
	Username         string      `gorm:"column:username;type:varchar(100);notnull;default:''" json:"username_seller"`
	Nama             string      `gorm:"column:nama;type:varchar(150);not null;default:''" json:"nama_seller"`
	Email            string      `gorm:"column:email;type:varchar(150);not null;default:''" json:"email_seller"`
	Jenis            JenisSeller `gorm:"column:jenis;type:varchar(250);not null;default:'Personal'" json:"jenis_seller"`
	Norek            string      `gorm:"column:norek;type:varchar(250);not null;default:''" json:"norek_seller"`
	SellerDedication SellerType  `gorm:"column:seller_dedication;type:varchar(250);not null;default:'Semua Barang'" json:"seller_dedication"`
	JamOperasional   string      `gorm:"column:jam_operasional;type:text;not null;default:''" json:"jam_operasional_seller"`
	Punchline        string      `gorm:"column:punchline;type:text;not null;default:''" json:"punchline_seller"`
	Password         string      `gorm:"column:password_hash;type:varchar(250);not null;default:''" json:"pass_seller"`
	Deskripsi        string      `gorm:"column:deskripsi;type:text;not null;default:''" json:"deskripsi_seller"`
	FollowerTotal    int32       `gorm:"column:follower_total;type:int4;not null;default:0" json:"follower_total_seller"`
	StatusSeller     Status      `gorm:"column:status;type:varchar(250);not null;default:'Offline'" json:"status_seller"`
	CreatedAt        time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        *time.Time  `gorm:"index" json:"deleted_at,omitempty"`
}

func (Seller) TableName() string {
	return "seller"
}

type JenisLayananKurir string

const (
	Reguler JenisLayananKurir = "Reguler"
	Express JenisLayananKurir = "Express"
	Ekonomi JenisLayananKurir = "Ekonomi"
	Sameday JenisLayananKurir = "Sameday"
	NextDay JenisLayananKurir = "NextDay"
	Cargo   JenisLayananKurir = "Cargo"
)

type Kurir struct {
	ID               int64             `gorm:"primaryKey;autoIncrement" json:"id_kurir"`
	Nama             string            `gorm:"column:nama;type:varchar(150);not null;default:''" json:"nama_kurir"`
	Email            string            `gom:"column:email;type:varchar(150);not null;default:''" json:"email_kurir"`
	Jenis            JenisLayananKurir `gorrm:"column:jenis;type:varchar(250);not null;default:'Reguler'" json:"jenis_kurir"`
	PasswordHash     string            `gorm:"column:password_hash;type:varchar(250);not null;default:''" json:"pass_kurir"`
	Deskripsi        string            `gorm:"column:deskripsi;type:text;not null;default:''" json:"deskripsi_kurir"`
	StatusKurir      Status            `gorm:"column:status;type:varchar(150);not null;default:'Offline'" json:"status_kurir"`
	JumlahPengiriman int32             `gorm:"column:jumlah_pengiriman;type:int4;not null;default:0"`
	CreatedAt        time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        *time.Time        `gorm:"index" json:"deleted_at,omitempty"`
}

func (Kurir) TableName() string {
	return "kurir"
}

package models

import (
	"time"
)

type Komentar struct {
	ID            int64       `gorm:"primaryKey;autoIncrement" json:"id_komentar"`
	IdBarangInduk int32       `gorm:"column:id_barang_induk;not null" json:"id_barang_induk"`
	baranginduk   BarangInduk `gorm:"foreignKey:IdBarangInduk;references:ID"`
	IdEntity      int64       `gorm:"column:id_entity;not null" json:"id_entity"`
	JenisEntity   string      `gorm:"column:jenis_entity;type:varchar(50);not null" json:"jenis_entity"`
	Komentar      string      `gorm:"column:komentar;type:text;not null" json:"isi_komentar"`
	ParentID      *int64      `gorm:"column:parent_id" json:"parent_id,omitempty"`
	CreatedAt     time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     *time.Time  `gorm:"index" json:"deleted_at,omitempty"`
}

func (Komentar) TableName() string {
	return "komentar"
}

type Keranjang struct {
	IdPengguna     int64          `gorm:"column:id_pengguna;not null" json:"id_pengguna_keranjang"`
	Pengguna       Pengguna       `gorm:"foreignKey:IdPengguna;references:ID"`
	IdSeller       int32          `gorm:"column:id_seller;not null" json:"id_seller_barang_induk_keranjang"`
	Seller         Seller         `gorm:"foreignKey:IdSeller;references:ID"`
	IdBarangInduk  int32          `gorm:"column:id_barang_induk;not null" json:"id_barang_induk_keranjang"`
	BarangInduk    BarangInduk    `gorm:"foreignKey:IdBarangInduk;references:ID"`
	IdKategori     int64          `gorm:"id_kategori_barang;not null" json:"id_kategori_barang_keranjang"`
	Kategoribarang KategoriBarang `gorm:"foreignKey:IdKategori;references:ID"`
	Count          int16          `gorm:"column:count;type:int2;not null" json:"count_keranjang"`
	Status         string         `gorm:"column:status;type:status_keranjang;not null" json:"status_keranjang"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      *time.Time     `gorm:"index" json:"deleted_at,omitempty"`
}

func (Keranjang) TableName() string {
	return "keranjang"
}

type BarangDisukai struct {
	IdPengguna    int64       `gorm:"column:id_pengguna;not null" json:"id_pengguna_barang_disukai"`
	Pengguna      Pengguna    `gorm:"foreignKey:IdPengguna;references:ID"`
	IdBarangInduk int32       `gorm:"column:id_barang_induk;not null" json:"id_barang_induk_disukai"`
	BarangInduk   BarangInduk `gorm:"foreignKey:IdBarangInduk;references:ID"`
	CreatedAt     time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     *time.Time  `gorm:"index" json:"deleted_at,omitempty"`
}

func (BarangDisukai) TableName() string {
	return "barang_disukai"
}

type Follower struct {
	IdFollower int64    `gorm:"column:id_follower;not null" json:"id_follower"`
	Pengguna   Pengguna `gorm:"foreignKey:IdFollower;references:ID"` // user yang follow
	IdFollowed int64    `gorm:"column:id_followed;not null" json:"id_followed"`
	Seller     Seller   `gorm:"foreignKey:IdFollowed;references:ID"` // seller yang di-follow
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `gorm:"index"`
}

func (Follower) TableName() string {
	return "follower"
}

type EntitySocialMedia struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id_social_media"`
	EntityId  int64      `gorm:"column:entity_id;type:int8;not null" json:"entity_id_social_media"`
	Whatsapp  string     `gorm:"column:whatsapp;type:varchar(20)" json:"whatsapp_social_media"`
	Facebook  string     `gorm:"column:facebook;type:text" json:"facebook_social_media"`
	TikTok    string     `gorm:"column:tiktok;type:text" json:"tiktok_social_media"`
	Instagram string     `gorm:"column:instagram;type:text" json:"instagram_social_media"`
	Metadata  []byte     `gorm:"column:metadata;type:bytea" json:"metadata_social_media"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"`
}

func (EntitySocialMedia) TableName() string {
	return "entity_social_media"
}

type AksiPengguna string

const (
	Registrasi       AksiPengguna = "Registrasi"
	Login            AksiPengguna = "Login"
	Logout           AksiPengguna = "Logout"
	Pembelian        AksiPengguna = "Pembelian"
	Tambah_keranjang AksiPengguna = "Tambah Keranjang"
	Hapus_keranjang  AksiPengguna = "Hapus Keranjang"
	Rating           AksiPengguna = "Rating"
	Update_profil    AksiPengguna = "Update Profil"
	Wishlist         AksiPengguna = "Wishlist"
	Pencarian        AksiPengguna = "Pencarian"
	Promo            AksiPengguna = "Promo"
)

type AktivitasPengguna struct {
	ID             int64      `gorm:"primaryKey;autoIncrement" json:"id_aktivitas_pengguna"`
	IdPengguna     int64      `gorm:"column:id_pengguna;not null" json:"id_pengguna_aktivitas_pengguna"`
	Pengguna       Pengguna   `gorm:"foreignKey:IdPengguna;references:ID"`
	WaktuDilakukan time.Time  `gorm:"column:waktu_dilakukan;autoCreateTime" json:"waktu_dilakukan_aktivitas_pengguna"`
	Aksi           string     `gorm:"column:aksi;type:aksi_pengguna" json:"aksi_aktivitas_pengguna"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime"`
	DeletedAt      *time.Time `gorm:"index"`
}

func (AktivitasPengguna) TableName() string {
	return "aktivitas_pengguna"
}

type AktivitasSeller struct {
	ID             int64      `gorm:"primaryKey;autoIncrement" json:"id_aktivitas_seller"`
	IdSeler        int32      `gorm:"column:id_seller;not null" json:"id_seller_aktivitas_seller"`
	seller         Seller     `gorm:"foreignKey:IdSeller;references:ID"`
	WaktuDilakukan time.Time  `gorm:"column:waktu_dilakukan;autoCreateTime" json:"waktu_dilakukan_aktivitas_seller"`
	Aksi           string     `gorm:"column:aksi;type:aksi_seller" json:"aksi_aktivitas_seller"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime"`
	DeletedAt      *time.Time `gorm:"index"`
}

func (AktivitasSeller) TableName() string {
	return "aktivitas_seller"
}

type Diskon struct {
	IdBarangInduk int64       `gorm:"column:id_barang_induk;not null" json:"id_barang_induk_diskon"`
	BarangInduk   BarangInduk `gorm:"foreignKey:IdBarangInduk;references:ID"`
	Deskripsi     string      `gorm:"column:deskripsi;type:text" json:"deskripsi_diskon"`
	Berlaku       time.Time   `gorm:"column:berlaku;type:date;not null" json:"berlaku_diskon"`
	Expired       time.Time   `gorm:"column:expired;type:date;not null" json:"expired_diskon"`
	CreatedAt     time.Time   `gorm:"autoCreateTime"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime"`
	DeletedAt     *time.Time  `gorm:"index"`
}

func (Diskon) TableName() string {
	return "diskon"
}

type AlamatPengguna struct {
	ID              int64    `gorm:"primaryKey;autoIncrement" json:"id_alamat_user"`
	IDPengguna      int64    `gorm:"column:id_pengguna;not null" json:"id_pengguna_alamat_user"`
	Pengguna        Pengguna `gorm:"foreignKey:IDPengguna;references:ID" json:"-"`
	PanggilanAlamat string   `gorm:"column:panggilan_alamat;type:varchar(250);not null" json:"panggilan_alamat_user"`
	NomorTelephone  string   `gorm:"column:nomor_telefon;type:varchar(20);not null" json:"nomor_telfon_alamat_user"`
	NamaAlamat      string   `gorm:"column:nama_alamat;type:text;not null" json:"nama_alamat_user"`
	Kota            string   `gorm:"column:kota;type:varchar(100);not null" json:"kota_alamat_user"`
	KodePos         string   `gorm:"column:kode_pos;type:varchar(40);not null" json:"kode_pos_alamat_user"`
	KodeNegara      string   `gorm:"column:kode_negara;default:'IDN';not null" json:"kode_negara_alamat_user"`
	Deskripsi       string   `gorm:"column:deskripsi;type:text;" json:"deskripsi_alamat_user"`
	Longitude       float64  `gorm:"column:longitude;type:decimal(10,8);" json:"longitude_alamat_user"`
	Latitude        float64  `gorm:"column:latitude;type:decimal(10,8);" json:"latitude_alamat_user"`
}

func (AlamatPengguna) TableName() string {
	return "alamat_pengguna"
}

type AlamatSeller struct {
	ID              int64   `gorm:"primaryKey;autoIncrement" json:"id_alamat_seller"`
	IDSeller        int32   `gorm:"column:id_seller;not null" json:"id_seller_alamat_seller"`
	Seller          Seller  `gorm:"foreignKey:IDSeller;references:ID"`
	PanggilanAlamat string  `gorm:"column:panggilan_alamat;type:varchar(250);not null" json:"panggilan_alamat_seller"`
	NomorTelephone  string  `gorm:"column:nomor_telefon;type:varchar(20);not null" json:"nomor_telfon_alamat_seller"`
	NamaAlamat      string  `gorm:"column:nama_alamat;type:text;not null" json:"nama_alamat_seller"`
	Deskripsi       string  `gorm:"column:deskripsi;type:text;" json:"deskripsi_alamat_seller"`
	Longitude       float64 `gorm:"column:longitude;type:decimal(10,8);" json:"longitude_alamat_seller"`
	Latitude        float64 `gorm:"column:latitude;type:decimal(10,8);" json:"latitude_alamat_seller"`
}

func (AlamatSeller) TableName() string {
	return "alamat_seller"
}

type RekeningSeller struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id_rekening_seller"`
	IDSeller        int32     `gorm:"column:id_seller;not null;index" json:"id_seller"`
	NamaBank        string    `gorm:"column:nama_bank;type:varchar(50);not null" json:"nama_bank_rekening_seller"`
	NomorRekening   string    `gorm:"column:nomor_rekening;type:varchar(50);not null" json:"nomor_rekening_seller"`
	PemilikRekening string    `gorm:"column:pemilik_rekening;type:varchar(100);not null" json:"pemilik_rekening_seller"`
	IsDefault       bool      `gorm:"column:id_default;default:false" json:"is_default_rekening_seller"`
	Status          string    `gorm:"column:status;type:varchar(20);default:'pending'" json:"status_rekening_seller"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (RekeningSeller) TableName() string {
	return "rekening_seller"
}

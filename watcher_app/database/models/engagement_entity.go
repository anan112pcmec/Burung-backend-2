package models

import (
	"time"
)

type Komentar struct {
	ID            int64       `gorm:"primaryKey;autoIncrement" json:"id_komentar"`
	IdBarangInduk int64       `gorm:"column:id_barang_induk;not null" json:"id_barang_induk_komentar"`
	BarangInduk   BarangInduk `gorm:"foreignKey:IdBarangInduk;references:ID"`
	IdPengguna    int64       `gorm:"column:id_pengguna;not null" json:"id_pengguna_komentar"`
	Pengguna      Pengguna    `gorm:"foreignKey:IdPengguna;references:ID"`
	Komentar      string      `gorm:"column:komentar;type:text;not null;default:''" json:"isi_komentar"`
	Rating        int16       `gorm:"column:rating;type:int2;default:0" json:"rating_komentar"`
	CreatedAt     time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     *time.Time  `gorm:"index" json:"deleted_at,omitempty"`
}

func (Komentar) TableName() string {
	return "komentar"
}

type Keranjang struct {
	IdPengguna    int64       `gorm:"column:id_pengguna;not null" json:"id_pengguna_keranjang"`
	Pengguna      Pengguna    `gorm:"foreignKey:IdPengguna;references:ID"`
	IdBarangInduk int64       `gorm:"column:id_barang_induk;not null" json:"id_barang_induk_keranjang"`
	BarangInduk   BarangInduk `gorm:"foreignKey:IdBarangInduk;references:ID"`
	Count         int16       `gorm:"column:count;type:int2;not null" json:"count_keranjang"`
	CreatedAt     time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     *time.Time  `gorm:"index" json:"deleted_at,omitempty"`
}

func (Keranjang) TableName() string {
	return "keranjang"
}

type BarangDisukai struct {
	IdPengguna    int64       `gorm:"column:id_pengguna;not null" json:"id_pengguna_barang_disukai"`
	Pengguna      Pengguna    `gorm:"foreignKey:IdPengguna;references:ID"`
	IdBarangInduk int64       `gorm:"column:id_barang_induk;not null" json:"id_barang_induk_disukai"`
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
	Aksi           string     `gorm:"column:aksi;type:aksi_pengguna;not null" json:"aksi_aktivitas_pengguna"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime"`
	DeletedAt      *time.Time `gorm:"index"`
}

func (AktivitasPengguna) TableName() string {
	return "aktivitas_pengguna"
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

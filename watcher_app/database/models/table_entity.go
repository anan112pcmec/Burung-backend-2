package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Pengguna struct {
	ID             int64          `gorm:"primaryKey;autoIncrement" json:"id_user"`
	Username       string         `gorm:"column:username;type:varchar(100);not null;default:''" json:"username_user"`
	Nama           string         `gorm:"column:nama;type:text;not null;default:''" json:"nama_user"`
	Email          string         `gorm:"column:email;type:varchar(100);not null;uniqueIndex" json:"email_user"`
	PasswordHash   string         `gorm:"column:password_hash;type:varchar(250);not null;default:''" json:"pass_user"`
	PinHash        string         `gorm:"column:pin_hash;type:varchar(250);not null;default:''" json:"pin_user"`
	StatusPengguna string         `gorm:"column:status;type:status;not null;default:'Offline'" json:"status_user"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Pengguna) TableName() string {
	return "pengguna"
}

type Seller struct {
	ID               int32          `gorm:"primaryKey;autoIncrement" json:"id_seller"`
	Username         string         `gorm:"column:username;type:varchar(100);notnull;default:''" json:"username_seller"`
	Nama             string         `gorm:"column:nama;type:varchar(150);not null;default:''" json:"nama_seller"`
	Email            string         `gorm:"column:email;type:varchar(150);not null;default:''" json:"email_seller"`
	Jenis            string         `gorm:"column:jenis;type:jenis_seller;not null;default:'Personal'" json:"jenis_seller"`
	SellerDedication string         `gorm:"column:seller_dedication;type:seller_dedication;not null;default:'Semua Barang'" json:"seller_dedication"`
	JamOperasional   string         `gorm:"column:jam_operasional;type:text;not null;default:''" json:"jam_operasional_seller"`
	Punchline        string         `gorm:"column:punchline;type:text;not null;default:''" json:"punchline_seller"`
	Password         string         `gorm:"column:password_hash;type:varchar(250);not null;default:''" json:"pass_seller"`
	Deskripsi        string         `gorm:"column:deskripsi;type:text;not null;default:''" json:"deskripsi_seller"`
	FollowerTotal    int32          `gorm:"column:follower_total;type:int4;not null;default:0" json:"follower_total_seller"`
	StatusSeller     string         `gorm:"column:status;type:status;not null;default:'Offline'" json:"status_seller"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (s *Seller) Validating() error {
	if s.ID == 0 {
		return fmt.Errorf("id tidak valid")
	}
	if s.Username == "" {
		return fmt.Errorf("username tidak valid")
	}
	if s.Email == "" {
		return fmt.Errorf("email tidak valid")
	}
	return nil
}

func (Seller) TableName() string {
	return "seller"
}

type JenisLayananKurir string

type Kurir struct {
	ID               int64          `gorm:"primaryKey;autoIncrement" json:"id_kurir"`
	Nama             string         `gorm:"column:nama;type:varchar(150);not null;default:''" json:"nama_kurir"`
	Username         string         `gorm:"column:username;type:text;not null" json:"username_kurir"`
	Email            string         `gorm:"column:email;type:varchar(150);not null;default:''" json:"email_kurir"`
	Jenis            string         `gorm:"column:jenis;type:jenis_layanan_kurir;not null;default:'Reguler'" json:"jenis_kurir"`
	PasswordHash     string         `gorm:"column:password_hash;type:varchar(250);not null;default:''" json:"pass_kurir"`
	Deskripsi        string         `gorm:"column:deskripsi;type:text;not null;default:''" json:"deskripsi_kurir"`
	StatusKurir      string         `gorm:"column:status;type:status;not null;default:'Offline'" json:"status_kurir"`
	StatusNarik      string         `gorm:"column:status_narik;type:status_kurir_narik;not null; default:'Off'" json:"status_narik_kurir"`
	VerifiedKurir    bool           `gorm:"column:verified;type:boolean;not null;default:false" json:"verified_kurir"`
	JumlahPengiriman int32          `gorm:"column:jumlah_pengiriman;type:int4;not null;default:0" json:"jumlah_pengiriman_kurir"`
	Balance          int64          `gorm:"column:balance_kurir;type:int8;default:0" json:"balance_kurir"`
	Rating           float32        `gorm:"column:rating;type:float;default:0" json:"rating_kurir"`
	JumlahRating     int32          `gorm:"column:jumlah_rating;type:int4;default:0" json:"jumlah_rating_kurir"`
	TipeKendaraan    string         `gorm:"column:tipe_kendaraan;type:varchar(50);default:''" json:"tipe_kendaraan_kurir"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Kurir) TableName() string {
	return "kurir"
}

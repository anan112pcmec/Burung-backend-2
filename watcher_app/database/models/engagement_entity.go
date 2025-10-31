package models

import (
	"time"

	"gorm.io/gorm"
)

// ///////////////////////////////////////////////////////////////////////////////////////////
// ENGAGEMENT PENGGUNA
// ///////////////////////////////////////////////////////////////////////////////////////////

type Komentar struct {
	ID            int64       `gorm:"primaryKey;autoIncrement" json:"id_komentar"`
	IdBarangInduk int32       `gorm:"column:id_barang_induk;not null" json:"id_barang_induk_komentar"`
	baranginduk   BarangInduk `gorm:"foreignKey:IdBarangInduk;references:ID"`
	IdEntity      int64       `gorm:"column:id_entity;not null" json:"id_entity_komentar"`
	JenisEntity   string      `gorm:"column:jenis_entity;type:varchar(50);not null" json:"jenis_entity_komentar"`
	Komentar      string      `gorm:"column:komentar;type:text;not null" json:"isi_komentar"`
	ParentID      *int64      `gorm:"column:parent_id" json:"parent_id_komentar,omitempty"`
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

type EntitySocialMedia struct {
	ID         int64      `gorm:"primaryKey;autoIncrement" json:"id_social_media"`
	EntityId   int64      `gorm:"column:entity_id;type:int8;not null" json:"entity_id_social_media"`
	Whatsapp   string     `gorm:"column:whatsapp;type:varchar(20)" json:"whatsapp_social_media"`
	Facebook   string     `gorm:"column:facebook;type:text" json:"facebook_social_media"`
	TikTok     string     `gorm:"column:tiktok;type:text" json:"tiktok_social_media"`
	Instagram  string     `gorm:"column:instagram;type:text" json:"instagram_social_media"`
	Metadata   []byte     `gorm:"column:metadata;type:bytea" json:"metadata_social_media"`
	EntityType string     `gorm:"column:entity_type;type:varchar(20);not null" json:"entity_type_social_media"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (EntitySocialMedia) TableName() string {
	return "entity_social_media"
}

type AktivitasPengguna struct {
	ID             int64      `gorm:"primaryKey;autoIncrement" json:"id_aktivitas_pengguna"`
	IdPengguna     int64      `gorm:"column:id_pengguna;not null" json:"id_pengguna_aktivitas_pengguna"`
	Pengguna       Pengguna   `gorm:"foreignKey:IdPengguna;references:ID"`
	WaktuDilakukan time.Time  `gorm:"column:waktu_dilakukan;autoCreateTime" json:"waktu_dilakukan_aktivitas_pengguna"`
	Aksi           string     `gorm:"column:aksi;type:aksi_pengguna" json:"aksi_aktivitas_pengguna"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (AktivitasPengguna) TableName() string {
	return "aktivitas_pengguna"
}

type AktivitasSeller struct {
	ID             int64      `gorm:"primaryKey;autoIncrement" json:"id_aktivitas_seller"`
	IdSeler        int32      `gorm:"column:id_seller;not null" json:"id_seller_aktivitas_seller"`
	Seller         Seller     `gorm:"foreignKey:IdSeller;references:ID" json:"-"`
	WaktuDilakukan time.Time  `gorm:"column:waktu_dilakukan;autoCreateTime" json:"waktu_dilakukan_aktivitas_seller"`
	Aksi           string     `gorm:"column:aksi;type:aksi_seller" json:"aksi_aktivitas_seller"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (AktivitasSeller) TableName() string {
	return "aktivitas_seller"
}

type AlamatPengguna struct {
	ID              int64          `gorm:"primaryKey;autoIncrement" json:"id_alamat_user"`
	IDPengguna      int64          `gorm:"column:id_pengguna;not null" json:"id_pengguna_alamat_user"`
	Pengguna        Pengguna       `gorm:"foreignKey:IDPengguna;references:ID" json:"-"`
	PanggilanAlamat string         `gorm:"column:panggilan_alamat;type:varchar(250);not null" json:"panggilan_alamat_user"`
	NomorTelephone  string         `gorm:"column:nomor_telefon;type:varchar(20);not null" json:"nomor_telfon_alamat_user"`
	NamaAlamat      string         `gorm:"column:nama_alamat;type:text;not null" json:"nama_alamat_user"`
	Kota            string         `gorm:"column:kota;type:varchar(100);not null" json:"kota_alamat_user"`
	KodePos         string         `gorm:"column:kode_pos;type:varchar(40);not null" json:"kode_pos_alamat_user"`
	KodeNegara      string         `gorm:"column:kode_negara;default:'IDN';not null" json:"kode_negara_alamat_user"`
	Deskripsi       string         `gorm:"column:deskripsi;type:text;" json:"deskripsi_alamat_user"`
	Longitude       float64        `gorm:"column:longitude;type:decimal(10,8);" json:"longitude_alamat_user"`
	Latitude        float64        `gorm:"column:latitude;type:decimal(10,8);" json:"latitude_alamat_user"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (AlamatPengguna) TableName() string {
	return "alamat_pengguna"
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// ENGAGEMENT SELLER
// ///////////////////////////////////////////////////////////////////////////////////////////

type Jenis_Seller struct {
	ID               int64      `gorm:"primaryKey;autoIncrement" json:"id_jenis_seller"`
	IdSeller         int32      `gorm:"column:id_seller;not null" json:"id_seller_jenis_seller"`
	Seller           Seller     `gorm:"foreignKey:IdSeller;references:ID" json:"-"`
	ValidationStatus string     `gorm:"column:validation_status; not null; default:'Pending'" json:"validation_status_jenis_seller"`
	Alasan           string     `gorm:"alasan_seller;type:text" json:"alasan_seller_jenis_seller"`
	AlasanAdmin      string     `gorm:"alasan_admin;type:text" json:"alasan_admin_jenis_seller"`
	TargetJenis      string     `gorm:"column:target_jenis;type:jenis_seller" json:"target_jenis_seller"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (Jenis_Seller) TableName() string {
	return "jenis_seller_validation"
}

type AlamatSeller struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id_alamat_seller"`
	IDSeller        int32      `gorm:"column:id_seller;not null" json:"id_seller_alamat_seller"`
	Seller          Seller     `gorm:"foreignKey:IDSeller;references:ID"`
	PanggilanAlamat string     `gorm:"column:panggilan_alamat;type:varchar(250);not null" json:"panggilan_alamat_seller"`
	NomorTelephone  string     `gorm:"column:nomor_telefon;type:varchar(20);not null" json:"nomor_telfon_alamat_seller"`
	NamaAlamat      string     `gorm:"column:nama_alamat;type:text;not null" json:"nama_alamat_seller"`
	Deskripsi       string     `gorm:"column:deskripsi;type:text;" json:"deskripsi_alamat_seller"`
	Longitude       float64    `gorm:"column:longitude;type:decimal(10,8);" json:"longitude_alamat_seller"`
	Latitude        float64    `gorm:"column:latitude;type:decimal(10,8);" json:"latitude_alamat_seller"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (AlamatSeller) TableName() string {
	return "alamat_seller"
}

type BatalTransaksi struct {
	ID             int64      `gorm:"primaryKey;autoIncrement" json:"id_batal_transaksi"`
	IdTransaksi    int64      `gorm:"column:id_transaksi;not null" json:"id_transaksi_batal_transaksi"`
	ITransaksi     Transaksi  `gorm:"foreignKey:IdTransaksi;references:ID" json:"-"`
	DibatalkanOleh string     `gorm:"column:dibatalkan_oleh;type:varchar(20);not null" json:"transaksi_dibatalkan_oleh"`
	Alasan         string     `gorm:"column:alasan;type:text;not null" json:"alasan_batal_transaksi"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (BatalTransaksi) TableName() string {
	return "batal_transaksi"
}

type Follower struct {
	IdFollower int64      `gorm:"column:id_follower;not null" json:"id_follower"`
	Pengguna   Pengguna   `gorm:"foreignKey:IdFollower;references:ID"` // user yang follow
	IdFollowed int64      `gorm:"column:id_followed;not null" json:"id_followed"`
	Seller     Seller     `gorm:"foreignKey:IdFollowed;references:ID"` // seller yang di-follow
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (Follower) TableName() string {
	return "follower"
}

type Diskon struct {
	IdBarangInduk int64       `gorm:"column:id_barang_induk;not null" json:"id_barang_induk_diskon"`
	BarangInduk   BarangInduk `gorm:"foreignKey:IdBarangInduk;references:ID"`
	Deskripsi     string      `gorm:"column:deskripsi;type:text" json:"deskripsi_diskon"`
	Berlaku       time.Time   `gorm:"column:berlaku;type:date;not null" json:"berlaku_diskon"`
	Expired       time.Time   `gorm:"column:expired;type:date;not null" json:"expired_diskon"`
	CreatedAt     time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     *time.Time  `gorm:"index" json:"deleted_at,omitempty"`
}

func (Diskon) TableName() string {
	return "diskon"
}

type RekeningSeller struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id_rekening_seller"`
	IDSeller        int32      `gorm:"column:id_seller;not null;index" json:"id_seller"`
	NamaBank        string     `gorm:"column:nama_bank;type:varchar(50);not null" json:"nama_bank_rekening_seller"`
	NomorRekening   string     `gorm:"column:nomor_rekening;type:varchar(50);not null" json:"nomor_rekening_seller"`
	PemilikRekening string     `gorm:"column:pemilik_rekening;type:varchar(100);not null" json:"pemilik_rekening_seller"`
	IsDefault       bool       `gorm:"column:id_default;default:false" json:"is_default_rekening_seller"`
	Status          string     `gorm:"column:status;type:varchar(20);default:'pending'" json:"status_rekening_seller"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (RekeningSeller) TableName() string {
	return "rekening_seller"
}

// ///////////////////////////////////////////////////////////////////////////////////////////
// ENGAGEMENT KURIR
// ///////////////////////////////////////////////////////////////////////////////////////////

// Pertimbangkan Ulang

type BalanceKurirLog struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id_balance_kurir"`
	KurirID   int64      `gorm:"column:kurir_id;not null" json:"kurir_id"`
	Kurir     Kurir      `gorm:"foreignKey:KurirID;references:ID" json:"-"`
	Amount    int64      `gorm:"column:amount;type:bigint;default:0" json:"amount_balance_kurir"`
	Type      string     `gorm:"column:type;type:varchar(10);default:'credit'" json:"type_balance_kurir"`
	Catatan   string     `gorm:"column:catatan;type:text" json:"catatan_balance_kurir"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (BalanceKurirLog) TableName() string {
	return "balance_kurir_log"
}

type InformasiKendaraanKurir struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id_informasi_kendaraan_kurir"`
	IDkurir         int64      `gorm:"column:id_kurir;not null" json:"id_kurir_informasi_kendaraan_kurir"`
	Kurir           Kurir      `gorm:"foreignKey:IDkurir;references:ID" json:"-"`
	JenisKendaraan  string     `gorm:"column:jenis_kendaraan;type:jenis_kendaraan_kurir;not null; default:'Motor'" json:"jenis_kendaraan_informasi_kendaraan_kurir"`
	NamaKendaraan   string     `gorm:"column:nama_kendaraan;type:text;not null" json:"nama_kendaraan_informasi_kendaraan_kurir"`
	RodaKendaraan   string     `gorm:"column:roda_kendaraan;type:roda_kendaraan;not null" json:"roda_kendaraan_informasi_kendaraan_kurir"`
	STNK            bool       `gorm:"column:informasi_stnk;type:boolean;not null; default:false" json:"informasi_stnk_informasi_kendaraan_kurir"`
	BPKB            bool       `gorm:"column:informasi_bpkb;type:boolean;not null; default:false" json:"informasi_bpkb_informasi_kendaraan_kurir"`
	StatusPerizinan string     `gorm:"column:status;type:status_perizinan_kendaraan;not null; default:'Pending'" json:"status_informasi_kendaraan_kurir"`
	CreatedAt       time.Time  `gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime"`
	DeletedAt       *time.Time `gorm:"index"`
}

func (InformasiKendaraanKurir) TableName() string {
	return "informasi_kendaraan_kurir"
}

type InformasiKurir struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id_informasi_kurir"`
	IDkurir         int64      `gorm:"column:id_kurir;not null" json:"id_kurir_informasi_kurir"`
	Kurir           Kurir      `gorm:"foreignKey:IDkurir;references:ID" json:"-"`
	Umur            int8       `gorm:"column:umur;type:int;not null" json:"umur_informasi_kurir"`
	Alasan          string     `gorm:"column:alasan;type:text" json:"alasan_informasi_kurir"`
	Ktp             bool       `gorm:"column:informasi_ktp;type:boolean;not null;default:false" json:"informasi_ktp_informasi_kurir"`
	Alamat          string     `gorm:"column:alamat;type:text" json:"alamat_informasi_kurir"`
	StatusPerizinan string     `gorm:"column:status;type:status_perizinan_kendaraan;not null; default:'Pending'" json:"status_informasi_kurir"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (InformasiKurir) TableName() string {
	return "informasi_kurir"
}

type AlamatGudang struct {
	ID              int64          `gorm:"primaryKey;autoIncrement" json:"id_alamat_gudang"`
	IDSeller        int32          `gorm:"column:id_seller;not null" json:"id_seller_alamat_gudang"`
	Seller          Seller         `gorm:"foreignKey:IDSeller;references:ID" json:"-"`
	PanggilanAlamat string         `gorm:"column:panggilan_alamat;type:varchar(250);not null" json:"panggilan_alamat_gudang"`
	NomorTelephone  string         `gorm:"column:nomor_telefon;type:varchar(20);not null" json:"nomor_telfon_alamat_gudang"`
	NamaAlamat      string         `gorm:"column:nama_alamat;type:text;not null" json:"nama_alamat_gudang"`
	Kota            string         `gorm:"column:kota;type:varchar(100);not null" json:"kota_alamat_gudang"`
	KodePos         string         `gorm:"column:kode_pos;type:varchar(40);not null" json:"kode_pos_alamat_gudang"`
	KodeNegara      string         `gorm:"column:kode_negara;default:'IDN';not null" json:"kode_negara_alamat_gudang"`
	Deskripsi       string         `gorm:"column:deskripsi;type:text;" json:"deskripsi_alamat_gudang"`
	Longitude       float64        `gorm:"column:longitude;type:decimal(10,8);" json:"longitude_alamat_gudang"`
	Latitude        float64        `gorm:"column:latitude;type:decimal(10,8);" json:"latitude_alamat_gudang"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (AlamatGudang) TableName() string {
	return "alamat_gudang"
}

package models

import "time"

type BarangContract interface {
	Validating() string
}

type BarangInduk struct {
	ID               int32      `gorm:"primaryKey;autoIncrement" json:"id_barang_induk"`
	SellerID         int32      `gorm:"column:id_seller;not null" json:"id_seller_barang_induk"`
	seller           Seller     `gorm:"foreignKey:SellerID;references:ID"`
	NamaBarang       string     `gorm:"column:nama_barang;type:varchar(200);not null" json:"nama_barang_induk"`
	JenisBarang      string     `gorm:"column:jenis_barang;type:seller_dedication;not null;default:'Semua Barang'" json:"jenis_barang_induk,omitempty"`
	OriginalKategori string     `gorm:"column:original_kategori;type:varchar(250)" json:"original_kategori,omitempty"`
	Deskripsi        string     `gorm:"column:deskripsi;type:text" json:"deskripsi_barang_induk,omitempty"`
	TanggalRilis     string     `gorm:"column:tanggal_rilis;type:date;not null" json:"tanggal_rilis_barang_induk,omitempty"`
	Viewed           int32      `gorm:"column:viewed;type:int4;not null;default:0" json:"viewed_barang_induk,omitempty"`
	Likes            int32      `gorm:"column:likes;type:int4;not null;default:0" json:"likes_barang_induk,omitempty"`
	TotalKomentar    int32      `gorm:"column:total_komentar;type:int4;not null;default:0" json:"total_komentar_barang_induk,omitempty"`
	HargaKategoris   int32      `gorm:"-" json:"harga_kategori_barang"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt        *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (b *BarangInduk) Validating() string {
	if b.SellerID == 0 {
		return "Gagal: IdSeller kosong"
	}
	if b.NamaBarang == "" {
		return "Gagal: Nama Barang kosong"
	}
	if b.JenisBarang == "" {
		return "Gagal: Jenis Barang kosong"
	}
	if b.Deskripsi == "" {
		return "Gagal: Deskripsi Barang kosong"
	}
	if b.TanggalRilis == "" {
		return "Gagal: Tanggal Rilis kosong"
	}

	return "Data Lengkap"
}

func (BarangInduk) TableName() string {
	return "barang_induk"
}

type KategoriBarang struct {
	ID             int64       `gorm:"primaryKey;autoIncrement" json:"id_kategori_barang"`
	IdBarangInduk  int32       `gorm:"column:id_barang_induk;not null" json:"id_barang_induk_kategori"`
	barangInduk    BarangInduk `gorm:"foreignKey:IdBarangInduk;references:ID"`
	Nama           string      `gorm:"column:nama;type:varchar(120);not null" json:"nama_kategori_barang"`
	Deskripsi      string      `gorm:"column:deskripsi;type:text" json:"deskripsi_kategori_barang"`
	Warna          string      `gorm:"column:warna;type:varchar(50)" json:"warna_kategori_barang"`
	Stok           int32       `gorm:"column:stok;type:int4;not null" json:"stok_kategori_barang"`
	Harga          int32       `gorm:"column:harga;type:int4;not null" json:"harga_kategori_barang"`
	BeratGram      int16       `gorm:"column:berat_gram;type:int2" json:"berat_gram_kategori_barang"`
	DimensiPanjang int16       `gorm:"column:dimensi_panjang_cm;type:int2" json:"dimensi_panjang_cm_kategori_barang"`
	DimensiLebar   int16       `gorm:"column:dimensi_lebar_cm;type:int2" json:"dimensi_tinggi_cm_kategori_barang"`
	Sku            string      `json:"sku_kategori"`
	CreatedAt      time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      *time.Time  `gorm:"index" json:"deleted_at,omitempty"`
}

func (KategoriBarang) TableName() string {
	return "kategori_barang"
}

type StatusVarianBarang string

const (
	Ready    StatusVarianBarang = "Ready"
	Dipesan  StatusVarianBarang = "Dipesan"
	Diproses StatusVarianBarang = "Diproses"
	Terjual  StatusVarianBarang = "Terjual"
)

type VarianBarang struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id_varian_barang"`
	IdBarangInduk int32          `gorm:"column:id_barang_induk;not null" json:"id_barang_induk_varian_barang"`
	barangInduk   BarangInduk    `gorm:"foreignKey:IdBarangInduk;references:ID"`
	IdKategori    int64          `gorm:"column:id_kategori;not null" json:"id_kategori_varian_barang"`
	Kategori      KategoriBarang `gorm:"foreignKey:IdKategori;references:ID"`
	IdTransaksi   int64          `gorm:"column:id_transaksi;type:int8" json:"id_transksi_varian_barang,omitempty"`
	Sku           string         `gorm:"column:sku;type:varchar(100);not null" json:"Sku_varian_barang,omitempty"`
	Status        string         `gorm:"column:status;type:status_varian;not null;default:'Ready'" json:"status_varian_barang,omitempty"`
	HoldBy        int64          `gorm:"column:hold_by;type:int8;default:0" json:"hold_by_varian_barang"`
	HolderEntity  string         `gorm:"column:holder_entity;type:varchar(30)" json:"holder_entity_varian_barang"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt     *time.Time     `gorm:"index" json:"deleted_at,omitempty"`
}

func (VarianBarang) TableName() string {
	return "varian_barang"
}

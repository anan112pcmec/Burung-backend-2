package watcher_app

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/meilisearch/meilisearch-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/maintain"
	trigger "github.com/anan112pcmec/Burung-backend-2/watcher_app/triggers"
)

type Databases struct {
	DB *gorm.DB // query biasa via GORM
}

type PostgreSettings struct {
	Host, User, Pass, Port, DBName string
}

func (data *Databases) InitializeWatcher(psg *PostgreSettings, ctx context.Context, wg *sync.WaitGroup) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		psg.Host, psg.User, psg.Pass, psg.DBName, psg.Port,
	)

	var err error
	data.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Sprintf("❌ Gagal koneksi GORM: %v", err))
	}
	fmt.Println("✅ Berhasil koneksi GORM")

	redisEntityCache := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})
	redisBarangCache := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       2,
	})
	redisEngagementCache := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       3,
	})

	keysnya := "SuperSecureKey1234567890"

	SearchEngine := meilisearch.New("http://localhost:7700", meilisearch.WithAPIKey(keysnya))

	barangIndukIndex := SearchEngine.Index("barang_induk_all")

	attrs := []interface{}{"jenis_barang_induk", "nama_barang_induk", "id_seller_barang_induk", "tanggal_rilis_barang_induk"}
	task2, err2 := barangIndukIndex.UpdateFilterableAttributes(&attrs)
	if err2 != nil {
		log.Fatal("❌ Gagal update filterable attributes:", err2)
	}
	log.Println("✅ Task UID:", task2.TaskUID)

	var currentDB string
	data.DB.Raw("SELECT current_database();").Scan(&currentDB)
	fmt.Println("Database aktif:", currentDB)

	if err := trigger.SetupEntityTriggers(data.DB); err != nil {
		fmt.Println(" Gagal Membuat Trigger Entity", err)
	} else {
		fmt.Println(" Berhasil Membuat Trigger Entity")
	}

	if err := trigger.SetupBarangTriggers(data.DB); err != nil {
		fmt.Println(" Gagal Membuat Trigger Barang", err)
	} else {
		fmt.Println(" Berhasil Membuat Trigger Barang")
	}

	if err := trigger.SetupKomentarTriggers(data.DB); err != nil {
		fmt.Println(" Gagal Membuat Trigger Komentar", err)
	} else {
		fmt.Println(" Berhasil Membuat Trigger Komentar")
	}

	if err := trigger.SetupTransaksiTriggers(data.DB); err != nil {
		fmt.Println(" Gagal Membuat Transaksi Trigger")
	} else {
		fmt.Println(" Berhasil Membuat Trigger Transaksi")
	}

	if err := trigger.SetupInformasiKurirTriggers(data.DB); err != nil {
		fmt.Println(" Gagal Membuat Informasi Kurir Trigger")
	} else {
		fmt.Println(" Berhasil Membuat Trigger Informasi Kurir")
	}

	if err := trigger.SetupPengirimanTriggers(data.DB); err != nil {
		fmt.Println(" Gagal Membuat Pengiriman Trigger")
	} else {
		fmt.Println(" Berhasil Membuat Trigger Pengiriman")
	}

	wg.Add(10)
	go func() {
		defer wg.Done()
		fmt.Println("Maintain Barang Jalan")
		maintain.BarangMaintainLoop(ctx, data.DB, redisBarangCache, SearchEngine)
	}()
	go func() {
		defer wg.Done()
		fmt.Println("Maintain Engagement Jalan")
		maintain.EngagementMaintainLoop(ctx, data.DB, redisEngagementCache)
	}()
	go func() {
		defer wg.Done()
		Pengguna_Watcher(ctx, dsn, data.DB, redisEntityCache)
	}()
	go func() {
		defer wg.Done()
		Seller_Watcher(ctx, dsn, data.DB, redisEntityCache)
	}()
	go func() {
		defer wg.Done()
		Barang_Induk_Watcher(ctx, dsn, data.DB, redisBarangCache, SearchEngine)
	}()
	go func() {
		defer wg.Done()
		Varian_Barang_Watcher(ctx, dsn, data.DB)
	}()
	go func() {
		defer wg.Done()
		Komentar_Barang_Watcher(ctx, dsn, redisEngagementCache)
	}()
	go func() {
		defer wg.Done()
		Transaksi_Watcher(ctx, dsn, data.DB)
	}()
	go func() {
		defer wg.Done()
		Informasi_Kurir_Watcher(ctx, dsn, data.DB)
	}()
	go func() {
		defer wg.Done()
		Informasi_Pengiriman_Watcher(ctx, dsn, data.DB)
	}()
}

func Getenvi(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("❌ Tidak ada file .env")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	var watcher = Databases{}

	var postgreconfig = PostgreSettings{
		Host:   Getenvi("DBHOST", ""),
		User:   Getenvi("DBUSER", ""),
		Pass:   Getenvi("DBPASS", ""),
		Port:   Getenvi("DBPORT", ""),
		DBName: Getenvi("DBNAME", ""),
	}

	watcher.InitializeWatcher(&postgreconfig, ctx, &wg)

	fmt.Println("🟢 Watcher berjalan... tekan CTRL+C untuk exit")
	wg.Wait()
}

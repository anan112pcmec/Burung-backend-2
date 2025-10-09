package config

import (
	"fmt"
	"log"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	ENVFILE = "env"
	YAML    = "yaml"
	JSON    = "json"
)

type Environment struct {
	DBHOST, DBUSER, DBPASS, DBNAME, DBPORT           string
	RDSHOST, RDSPORT                                 string
	RDSENTITYDB, RDSBARANGDB, RDSENGAGEMENTDB        int
	MEILIHOST, MEILIKEY, MEILIPORT                   string
	RMQ_HOST, RMQ_USER, RMQ_PASS, EXCHANGE, RMQ_PORT string
	RMQ_NOTIF_EXCHANGE                               string
}

func (e *Environment) RunConnectionEnvironment() (
	db *gorm.DB,
	redis_entity *redis.Client,
	redis_barang *redis.Client,
	redis_engagement *redis.Client,
	search_engine meilisearch.ServiceManager,
	notification *amqp091.Connection,
) {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		e.DBHOST, e.DBUSER, e.DBPASS, e.DBNAME, e.DBPORT,
	)

	log.Println("🔍 Mencoba koneksi ke PostgreSQL...")
	log.Println("🔗 DSN:", dsn)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // pakai level Warn agar log tidak terlalu ramai
	})
	if err != nil {
		log.Fatalf("❌ Gagal konek ke PostgreSQL: %v", err)
	}

	// Coba koneksi langsung
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Gagal mendapatkan *sql.DB dari GORM: %v", err)
	}

	// Coba ping database untuk memastikan koneksi aktif
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Gagal ping ke PostgreSQL: %v", err)
	}

	// Atur pool koneksi
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	var currentDB string
	if err := db.Raw("SELECT current_database();").Scan(&currentDB).Error; err != nil {
		log.Printf("⚠️ Tidak bisa membaca nama database: %v", err)
	} else {
		log.Println("✅ Berhasil terkoneksi ke database:", currentDB)
	}

	redis_entity = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", e.RDSHOST, e.RDSPORT),
		Password: "",
		DB:       e.RDSENTITYDB,
	})

	redis_barang = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", e.RDSHOST, e.RDSPORT),
		Password: "",
		DB:       e.RDSBARANGDB,
	})

	redis_engagement = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", e.RDSHOST, e.RDSPORT),
		Password: "",
		DB:       e.RDSENGAGEMENTDB,
	})

	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", e.RMQ_USER, e.RMQ_PASS, e.RMQ_HOST, e.RMQ_PORT)
	notification, _ = amqp091.Dial(connStr)

	search_engine = meilisearch.New(fmt.Sprintf("http://%s:%s", e.MEILIHOST, e.MEILIPORT), meilisearch.WithAPIKey(e.MEILIKEY))

	return
}

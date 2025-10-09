package watcher_app

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
	"github.com/meilisearch/meilisearch-go"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/config"
)

type Connection struct {
	DB            *gorm.DB
	RDSENTITY     *redis.Client
	RDSBARANG     *redis.Client
	RDSENGAGEMENT *redis.Client
	SE            meilisearch.ServiceManager
	NOTIFICATION  *amqp091.Connection
}

func Getenvi(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Run() {
	var conn Connection
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("❌ Tidak ada file .env")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	rdsentity, _ := strconv.Atoi(Getenvi("RDSENTITY", "0"))
	rdsbarang, _ := strconv.Atoi(Getenvi("RDSBARANG", "0"))
	rdsengagement, _ := strconv.Atoi(Getenvi("RDSENGAGEMET", "0"))

	env := config.Environment{
		DBHOST:             Getenvi("DBHOST", "NIL"),
		DBUSER:             Getenvi("DBUSER", "NIL"),
		DBPASS:             Getenvi("DBPASS", "NIL"),
		DBNAME:             Getenvi("DBNAME", "NIL"),
		DBPORT:             Getenvi("DBPORT", "NIL"),
		RDSHOST:            Getenvi("RDSHOST", "NIL"),
		RDSPORT:            Getenvi("RDSPORT", "NIL"),
		RDSENTITYDB:        rdsentity,
		RDSBARANGDB:        rdsbarang,
		RDSENGAGEMENTDB:    rdsengagement,
		MEILIHOST:          Getenvi("MEILIHOST", "NIL"),
		MEILIPORT:          Getenvi("MEILIPORT", "NIL"),
		MEILIKEY:           Getenvi("MEILIKEY", "NIL"),
		RMQ_HOST:           Getenvi("RMQ_HOST", "NIL"),
		RMQ_USER:           Getenvi("RMQ_USER", "NIL"),
		RMQ_PASS:           Getenvi("RMQ_PASS", "NIL"),
		RMQ_PORT:           Getenvi("RMQ_PORT", "NIL"),
		RMQ_NOTIF_EXCHANGE: Getenvi("RMQ_NOTIF_EXCHANGE", "NIL"),
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		env.DBHOST, env.DBUSER, env.DBPASS, env.DBNAME, env.DBPORT,
	)

	conn.DB, conn.RDSENTITY, conn.RDSBARANG, conn.RDSENGAGEMENT, conn.SE, conn.NOTIFICATION = env.RunConnectionEnvironment()

	Watcher(&conn, ctx, &wg, dsn, env.RMQ_NOTIF_EXCHANGE)

	fmt.Println("🟢 Watcher berjalan... tekan CTRL+C untuk exit")
	wg.Wait()
	defer conn.NOTIFICATION.Close()
}

package watcher_app

import (
	"context"
	"fmt"
	"sync"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/dispatcher"
	"github.com/anan112pcmec/Burung-backend-2/watcher_app/maintain"
	maintain_mb "github.com/anan112pcmec/Burung-backend-2/watcher_app/message_broker/maintain"
	producer_mb "github.com/anan112pcmec/Burung-backend-2/watcher_app/message_broker/producer"
	trigger "github.com/anan112pcmec/Burung-backend-2/watcher_app/triggers"

)

func Watcher(connection *Connection, ctx context.Context, wg *sync.WaitGroup, dsn string, Exchange string) {

	if err := trigger.SetupEntityTriggers(connection.DB); err != nil {
		fmt.Println(" Gagal Membuat Trigger Entity", err)
	} else {
		fmt.Println(" Berhasil Membuat Trigger Entity")
	}

	if err := trigger.SetupBarangTriggers(connection.DB); err != nil {
		fmt.Println(" Gagal Membuat Trigger Barang", err)
	} else {
		fmt.Println(" Berhasil Membuat Trigger Barang")
	}

	if err := trigger.SetupEngagementEntityTriggers(connection.DB); err != nil {
		fmt.Println(" Gagal Membuat Trigger Komentar", err)
	} else {
		fmt.Println(" Berhasil Membuat Trigger Komentar")
	}

	if err := trigger.SetupTransaksiTriggers(connection.DB); err != nil {
		fmt.Println(" Gagal Membuat Transaksi Trigger")
	} else {
		fmt.Println(" Berhasil Membuat Trigger Transaksi")
	}

	if err := trigger.SetupPengirimanTriggers(connection.DB); err != nil {
		fmt.Println(" Gagal Membuat Pengiriman Trigger")
	} else {
		fmt.Println(" Berhasil Membuat Trigger Pengiriman")
	}

	if err := producer_mb.UpConnectionDefaults(Exchange, connection.NOTIFICATION); err != nil {
		fmt.Println(err)
	}

	wg.Add(13)
	go func() {
		defer wg.Done()
		fmt.Println("Maintain Barang Jalan")
		maintain.BarangMaintainLoop(ctx, connection.DB, connection.RDSBARANG, connection.SE)
	}()
	go func() {
		defer wg.Done()
		fmt.Println("Maintain Engagement Jalan")
		maintain.EngagementMaintainLoop(ctx, connection.DB, connection.RDSENGAGEMENT)
	}()
	go func() {
		defer wg.Done()
		fmt.Println("Maintain Entity Jalan")
		maintain.EntityMaintainLoop(ctx, connection.DB, connection.RDSENTITY, connection.SE)
	}()
	go func() {
		defer wg.Done()
		dispatcher.Pengguna_Watcher(ctx, dsn, connection.DB, connection.RDSENTITY, connection.NOTIFICATION)
	}()
	go func() {
		defer wg.Done()
		dispatcher.Seller_Watcher(ctx, dsn, connection.DB, connection.RDSENTITY, connection.NOTIFICATION)
	}()
	go func() {
		defer wg.Done()
		dispatcher.Barang_Induk_Watcher(ctx, dsn, connection.DB, connection.RDSBARANG, connection.SE)
	}()
	go func() {
		dispatcher.Kategori_Barang_Watcher(ctx, dsn, connection.DB)
	}()
	go func() {
		defer wg.Done()
		dispatcher.Varian_Barang_Watcher(ctx, dsn, connection.DB)
	}()
	go func() {
		defer wg.Done()
		dispatcher.Transaksi_Watcher(ctx, dsn, connection.DB, connection.NOTIFICATION)
	}()
	go func() {
		defer wg.Done()
		dispatcher.Informasi_Kurir_Watcher(ctx, dsn, connection.DB)
	}()
	go func() {
		defer wg.Done()
		dispatcher.Pengiriman_Watcher(ctx, dsn, connection.DB)
	}()
	go func() {
		defer wg.Done()
		maintain_mb.NotificationMaintainLoop(ctx, connection.DB, connection.NOTIFICATION, Exchange)
	}()
	go func() {
		defer wg.Done()
		dispatcher.Follower_Watcher(ctx, dsn, connection.DB, connection.RDSENTITY)
	}()

}

package services

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/anan112pcmec/Burung-backend-2/watcher_app/database/models"
	"github.com/anan112pcmec/Burung-backend-2/watcher_app/notify_payload"
)

func SellerFollowed(ctx context.Context, data notify_payload.NotifyResponseFollower, db *gorm.DB, rds *redis.Client) {

	key := fmt.Sprintf("seller_data:%v", data.IdFollowed)

	// Update follower_total di Redis
	if errds := rds.HIncrBy(ctx, key, "follower_total_seller", 1).Err(); errds != nil {
		if err := db.WithContext(ctx).
			Model(&models.Seller{}).
			Where(&models.Seller{ID: int32(data.IdFollowed)}).
			Update("follower_total", gorm.Expr("follower_total + ?", 1)).Error; err != nil {
			fmt.Printf("[DB ERROR] Gagal menambah follower seller ID %v di database: %v\n", data.IdFollowed, err)
		} else {
			fmt.Printf("[DB INFO] Berhasil menambah follower untuk seller ID %v di database\n", data.IdFollowed)
		}
	} else {
		fmt.Printf("[REDIS INFO] Berhasil menambah follower seller ID %v di Redis (key: %s)\n", data.IdFollowed, key)
	}
}

func SellerUnfollowed(ctx context.Context, data notify_payload.NotifyResponseFollower, db *gorm.DB, rds *redis.Client) {
	// Kunci Redis untuk seller terkait
	key := fmt.Sprintf("seller_data:%v", data.IdFollowed)

	// Kurangi follower_total di Redis
	if errds := rds.HIncrBy(ctx, key, "follower_total_seller", -1).Err(); errds != nil {
		if err := db.WithContext(ctx).
			Model(&models.Seller{}).
			Where(&models.Seller{ID: int32(data.IdFollowed)}).
			Update("follower_total", gorm.Expr("GREATEST(follower_total - 1, 0)")).Error; err != nil {
			fmt.Printf("[DB ERROR] Gagal mengurangi follower seller ID %v di database: %v\n", data.IdFollowed, err)
		} else {
			fmt.Printf("[DB INFO] Berhasil mengurangi follower untuk seller ID %v di database\n", data.IdFollowed)
		}
	} else {
		fmt.Printf("[REDIS INFO] Berhasil mengurangi follower seller ID %v di Redis (key: %s)\n", data.IdFollowed, key)
	}
}

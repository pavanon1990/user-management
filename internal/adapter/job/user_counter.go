package job

import (
	"context"
	"log"
	"time"

	"user-management-api/internal/core/port"
)

func StartUserCountLogger(ctx context.Context, userService port.UserService, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Println("start background job user counter every 10 sec")

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logUserCount(ctx, userService)
		}
	}
}

func logUserCount(ctx context.Context, userService port.UserService) {
	count, err := userService.CountUsers(ctx)
	if err != nil {
		log.Printf("count users error: %v", err)
		return
	}
	log.Printf("total users: %d", count)
}

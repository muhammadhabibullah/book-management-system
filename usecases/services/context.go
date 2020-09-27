package services

import (
	"context"
	"time"
)

const contextTimeout = 5

func setContextTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, time.Duration(contextTimeout)*time.Second)
}

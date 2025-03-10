package utils

import (
	"context"
	"os/signal"
	"syscall"
)

func Wait() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
}

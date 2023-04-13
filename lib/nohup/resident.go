package nohup

import (
	"context"
)

type IResident interface {
	Run(ctx context.Context, stopFn context.CancelFunc)
	Close()
}

func NewResident(parentCtx context.Context, programs ...IResident) {
	ctx, stop := context.WithCancel(parentCtx)
	defer stop()
	for _, program := range programs {
		program.Run(ctx, stop)
	}
	<-ctx.Done()
	stop()

	for _, program := range programs {
		program.Close()
	}
}

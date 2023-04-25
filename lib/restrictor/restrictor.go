package restrictor

import (
	"anto/lib/log"
	"context"
	"fmt"
	"github.com/golang-module/carbon"
	"golang.org/x/time/rate"
	"sync"
)

var (
	apiRestrictor  *Restrictor
	onceRestrictor sync.Once
)

func Singleton() *Restrictor {
	onceRestrictor.Do(func() {
		apiRestrictor = new(Restrictor)
	})
	return apiRestrictor
}

type Restrictor struct {
	instances sync.Map
}

func (o *Restrictor) Get(key string) *rate.Limiter {
	currentLimiter, isOk := o.instances.Load(key)
	if !isOk || currentLimiter == nil {
		currentLimiter = rate.NewLimiter(1, 1)
	}
	return currentLimiter.(*rate.Limiter)
}

func (o *Restrictor) Set(key string, limiter *rate.Limiter) {
	o.instances.Store(key, limiter)
}

func (o *Restrictor) Allow(key string) bool {
	return o.Get(key).Allow()
}

func (o *Restrictor) Wait(key string, ctx context.Context) error {
	tmpLimiter := o.Get(key)
	fmt.Println(carbon.Now(), key, tmpLimiter.Limit(), tmpLimiter.Burst())
	err := o.Get(key).Wait(ctx)
	if err != nil {
		log.Singleton().ErrorF("等待令牌异常(关键字: %s), 错误: %s", key, err.Error())
	}

	return err
}

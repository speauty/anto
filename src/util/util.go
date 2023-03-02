package util

import (
	"context"
)

// IsCtxDone 判断ctx是否结束
func IsCtxDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

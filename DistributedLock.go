package distributedlock

import (
	"time"
)

type IDistributedLock interface {
	TryGetLock(key string) (bool, Lock, error)
	GetLock(key string, timeout time.Duration) (Lock, error)
}

type Lock interface {
	Lock()
	Release()
}

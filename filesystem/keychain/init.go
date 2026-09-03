package keychain

import (
	"sync"
	"time"

	"github.com/pardnchiu/go-pkg/utils"
)

const (
	secretToolTimeout   = 5 * time.Second
	secretToolWaitDelay = 1 * time.Second
)

var (
	once         sync.Once
	service      string
	fallbackPath string
)

func Init(svc, fbPath string) {
	once.Do(func() {
		service = svc
		fallbackPath = utils.AbsPath("", fbPath)
	})
}

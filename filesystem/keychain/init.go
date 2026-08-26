package keychain

import (
	"sync"

	"github.com/pardnchiu/go-pkg/utils"
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

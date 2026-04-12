package keychain

import "sync"

var (
	once         sync.Once
	service      string
	fallbackPath string
)

func Init(svc, fbPath string) {
	once.Do(func() {
		service = svc
		fallbackPath = fbPath
	})
}

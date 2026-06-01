package dummy

import (
	"crypto/md5"
)

func anotherVulnerableFunction() {
	_ = md5.New()
}

package dummy

import (
	"crypto/rsa"
)

func myVulnerableFunction() {
	_, _ = rsa.GenerateKey(nil, 2048)
}

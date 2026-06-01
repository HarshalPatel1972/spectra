package main

import (
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha1"
)

// This is a test file to trigger the Spectra CI GitHub App
func generateLegacyKeys() {
	// CRITICAL RISK
	_ = rsa.GenerateKey(nil, 1024)

	// HIGH RISK
	_ = sha1.New()

	// MEDIUM RISK
	_ = md5.New()
}

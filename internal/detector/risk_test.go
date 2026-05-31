package detector

import (
	"testing"
)

func TestComputeQRS(t *testing.T) {
	tests := []struct {
		name      string
		algo      string
		keySize   int
		occur     int
		wantQRS   int
	}{
		{"CRITICAL - RSA default", "RSA", 0, 1, 90},
		{"CRITICAL - RSA small key", "RSA", 512, 1, 100}, // 90 + 10 = 100
		{"CRITICAL - RSA large key", "RSA", 4096, 1, 60}, // 90 - 10 - 20 = 60
		{"CRITICAL - RSA freq boost", "RSA", 0, 21, 100}, // 90 + 10 = 100
		{"SAFE - ML-KEM", "ML-KEM", 0, 1, 0},
		{"SAFE - AES256", "AES256", 256, 1, 0},
		{"EDGE - unknown algo empty", "", 0, 1, 0},
		{"EDGE - zero key size", "ECDSA", 0, 1, 90},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, found := LookupAlgorithm(tt.algo)
			if !found && tt.algo != "" {
				t.Fatalf("algorithm not found: %s", tt.algo)
			}
			got := ComputeQRS(info, tt.keySize, tt.occur)
			if got != tt.wantQRS {
				t.Errorf("ComputeQRS() = %v, want %v", got, tt.wantQRS)
			}
		})
	}
}

func TestQRSToBand(t *testing.T) {
	tests := []struct {
		qrs  int
		want RiskBand
	}{
		{100, BandCritical},
		{80, BandCritical},
		{79, BandHigh},
		{60, BandHigh},
		{59, BandMedium},
		{40, BandMedium},
		{39, BandLow},
		{20, BandLow},
		{19, BandSafe},
		{0, BandSafe},
	}
	for _, tt := range tests {
		if got := QRSToBand(tt.qrs); got != tt.want {
			t.Errorf("QRSToBand(%d) = %v, want %v", tt.qrs, got, tt.want)
		}
	}
}

func TestAggregateQRS(t *testing.T) {
	tests := []struct {
		name   string
		qrs    []int
		want   int
	}{
		{"empty", []int{}, 0},
		{"all safe", []int{0, 0, 0}, 0},
		{"one critical", []int{100}, 100},
		{"mixed", []int{100, 0}, 99}, // 100 weight = 1.01, 0 weight = 0.01 -> (101) / 1.02 = 99.01 -> 99
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AggregateQRS(tt.qrs); got != tt.want {
				t.Errorf("AggregateQRS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClassifyEffort(t *testing.T) {
	tests := []struct {
		name       string
		source     SourceType
		algo       string
		wantEffort EffortLevel
	}{
		{"DEPS source", SourceDeps, "RSA", EffortBlocked},
		{"CERT source RSA", SourceCert, "RSA", EffortMedium},
		{"CERT source other", SourceCert, "SHA1", EffortEasy},
		{"CONFIG source", SourceConfig, "DES", EffortEasy},
		{"CODE source hash", SourceCode, "MD5", EffortEasy},
		{"CODE source AES128", SourceCode, "AES128", EffortEasy},
		{"CODE source symmetric", SourceCode, "DES", EffortMedium},
		{"CODE source asymmetric", SourceCode, "RSA", EffortHard},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := ClassifyEffort(tt.source, tt.algo)
			if got != tt.wantEffort {
				t.Errorf("ClassifyEffort() = %v, want %v", got, tt.wantEffort)
			}
		})
	}
}

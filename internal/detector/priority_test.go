package detector

import (
	"math"
	"testing"
)

func TestPriorityScore(t *testing.T) {
	// Base scores for 1 occurrence:
	// EASY (1.0), MEDIUM (0.6), HARD (0.3), BLOCKED (0.1)
	
	// QRS=100, EASY, 1 occur -> 100 * 1.0 * (1 + ln(1+1)*0.05) = 103.4657
	score1 := PriorityScore(100, EffortEasy, 1)
	if math.Abs(score1-103.465) > 0.01 {
		t.Errorf("PriorityScore(100, EASY, 1) = %v, want ~103.465", score1)
	}

	// QRS=100, HARD, 1 occur -> 100 * 0.3 * (1 + ln(1+1)*0.05) = 31.039
	score2 := PriorityScore(100, EffortHard, 1)
	if math.Abs(score2-31.039) > 0.01 {
		t.Errorf("PriorityScore(100, HARD, 1) = %v, want ~31.039", score2)
	}
}

func TestBuildActionPlan(t *testing.T) {
	summaries := []FindingSummary{
		{Algorithm: "RSA", QRS: 90, MigrationEffort: EffortHard},
		{Algorithm: "RSA", QRS: 100, MigrationEffort: EffortHard}, // Max QRS 100
		{Algorithm: "SHA1", QRS: 70, MigrationEffort: EffortEasy},
		{Algorithm: "AES128", QRS: 25, MigrationEffort: EffortEasy},
	}

	plan := BuildActionPlan(summaries)
	if len(plan) != 3 {
		t.Fatalf("expected 3 items in action plan, got %d", len(plan))
	}

	// SHA1 should be first (70 * 1.0 = 70)
	// RSA should be second (100 * 0.3 = 30) (approx with freq boost)
	if plan[0].Algorithm != "SHA1" {
		t.Errorf("expected SHA1 to be first priority, got %s", plan[0].Algorithm)
	}
	if plan[1].Algorithm != "RSA" {
		t.Errorf("expected RSA to be second priority, got %s", plan[1].Algorithm)
	}
	if plan[2].Algorithm != "AES128" {
		t.Errorf("expected AES128 to be third priority, got %s", plan[2].Algorithm)
	}
}

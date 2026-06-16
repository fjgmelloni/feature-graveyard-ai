package feature

import (
	"testing"
	"time"
)

func TestAggregateUsageGroupsByFeature(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 6, 16, 0, 0, 0, 0, time.UTC)
	logs := []UsageLog{
		mustLog(t, "Relatorio", "u1", "2026-01-10", 2),
		mustLog(t, "Relatorio", "u2", "2026-03-01", 4),
		mustLog(t, "Dashboard", "u3", "2026-06-15", 12),
	}

	usage := AggregateUsage(logs, now)

	if len(usage) != 2 {
		t.Fatalf("expected 2 features, got %d", len(usage))
	}
	if usage[1].Feature != "Relatorio" {
		t.Fatalf("expected sorted relatorio entry, got %q", usage[1].Feature)
	}
	if usage[1].UniqueUsers != 2 || usage[1].TotalAccess != 6 {
		t.Fatalf("unexpected aggregation: %+v", usage[1])
	}
	if usage[1].LastAccess.Format(time.DateOnly) != "2026-03-01" {
		t.Fatalf("expected latest access date, got %s", usage[1].LastAccess.Format(time.DateOnly))
	}
}

func TestClassifyDeadFeature(t *testing.T) {
	t.Parallel()

	status, risk, confidence := Classify(AggregatedFeatureUsage{
		Feature:           "Relatorio",
		UniqueUsers:       2,
		TotalAccess:       2,
		DaysSinceAccess:   180,
		FrequencyPerMonth: 0.3,
	})

	if status != StatusDeadFeature {
		t.Fatalf("expected DEAD_FEATURE, got %s", status)
	}
	if risk != RiskLow {
		t.Fatalf("expected LOW risk, got %s", risk)
	}
	if confidence < 0.9 {
		t.Fatalf("expected high confidence, got %f", confidence)
	}
}

func mustLog(t *testing.T, featureName, userID, lastAccess string, totalAccess int) UsageLog {
	t.Helper()
	log, err := NewUsageLog(UsageLogInput{
		Feature:     featureName,
		UserID:      userID,
		LastAccess:  lastAccess,
		TotalAccess: totalAccess,
	})
	if err != nil {
		t.Fatal(err)
	}
	return log
}

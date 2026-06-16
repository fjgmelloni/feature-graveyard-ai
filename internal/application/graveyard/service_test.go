package graveyard

import (
	"context"
	"testing"
	"time"

	"feature-graveyard-ai/internal/domain/feature"
	"feature-graveyard-ai/internal/infra/repository"
)

func TestServiceReportBuildsGraveyardSummary(t *testing.T) {
	t.Parallel()

	logs := []feature.UsageLog{
		mustLog(t, "RelatorioExportacaoExcel", "123", "2025-12-01", 2),
		mustLog(t, "DashboardExecutivoV1", "456", "2026-06-15", 40),
	}
	repo := repository.NewMemoryUsageRepository(logs)
	service := NewService(repo, stubAnalyzer{}).WithClock(func() time.Time {
		return time.Date(2026, 6, 16, 0, 0, 0, 0, time.UTC)
	})

	report, err := service.Report(context.Background(), 180)
	if err != nil {
		t.Fatalf("expected report without error, got %v", err)
	}

	if report.TotalFeatures != 2 {
		t.Fatalf("expected 2 features, got %d", report.TotalFeatures)
	}
	if report.DeadFeatures != 1 {
		t.Fatalf("expected 1 dead feature, got %d", report.DeadFeatures)
	}
	if report.RemovalCandidates != 1 {
		t.Fatalf("expected 1 removal candidate, got %d", report.RemovalCandidates)
	}
	if report.Analyses[0].GeneratedBy != "stub" {
		t.Fatalf("expected stub analyzer, got %q", report.Analyses[0].GeneratedBy)
	}
}

func TestServiceIngestValidatesLogs(t *testing.T) {
	t.Parallel()

	repo := repository.NewMemoryUsageRepository(nil)
	service := NewService(repo, stubAnalyzer{})

	_, err := service.Ingest(context.Background(), []feature.UsageLogInput{
		{Feature: "", UserID: "123", LastAccess: "2026-01-10", TotalAccess: 1},
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

type stubAnalyzer struct{}

func (stubAnalyzer) Analyze(_ context.Context, usage feature.AggregatedFeatureUsage, _ feature.Status, _ feature.Risk) (ExecutiveAnalysis, error) {
	return ExecutiveAnalysis{
		Summary:            usage.Feature + " summary",
		BusinessImpact:     "impact",
		SuggestedAction:    "action",
		ExecutiveRationale: "rationale",
		GeneratedBy:        "stub",
	}, nil
}

func mustLog(t *testing.T, featureName, userID, lastAccess string, totalAccess int) feature.UsageLog {
	t.Helper()
	log, err := feature.NewUsageLog(feature.UsageLogInput{
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

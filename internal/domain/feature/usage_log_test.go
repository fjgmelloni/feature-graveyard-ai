package feature

import "testing"

func TestNewUsageLogValidatesInput(t *testing.T) {
	t.Parallel()

	_, err := NewUsageLog(UsageLogInput{
		Feature:     " ",
		UserID:      "123",
		LastAccess:  "2026-01-10",
		TotalAccess: 2,
	})
	if err == nil {
		t.Fatal("expected feature validation error")
	}

	_, err = NewUsageLog(UsageLogInput{
		Feature:     "Relatorio",
		UserID:      "123",
		LastAccess:  "10/01/2026",
		TotalAccess: 2,
	})
	if err == nil {
		t.Fatal("expected date validation error")
	}
}

func TestNewUsageLogCreatesValidLog(t *testing.T) {
	t.Parallel()

	log, err := NewUsageLog(UsageLogInput{
		Feature:     " RelatorioExportacaoExcel ",
		UserID:      "123",
		LastAccess:  "2026-01-10",
		TotalAccess: 2,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if log.Feature != "RelatorioExportacaoExcel" {
		t.Fatalf("expected trimmed feature name, got %q", log.Feature)
	}
}

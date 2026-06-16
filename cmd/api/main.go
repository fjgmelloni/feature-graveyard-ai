package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"feature-graveyard-ai/internal/application/graveyard"
	"feature-graveyard-ai/internal/domain/feature"
	"feature-graveyard-ai/internal/infra/ai"
	httpserver "feature-graveyard-ai/internal/infra/http"
	"feature-graveyard-ai/internal/infra/repository"
)

func main() {
	port := env("PORT", "8080")
	apiKey := os.Getenv("GEMINI_API_KEY")
	model := env("GEMINI_MODEL", "gemini-1.5-flash")

	repo := repository.NewMemoryUsageRepository(seedUsageLogs())
	analyzer := ai.NewGeminiAnalyzer(apiKey, model)
	service := graveyard.NewService(repo, analyzer)
	server := httpserver.NewServer(service, http.Dir("web"))

	log.Printf("Feature Graveyard AI running at http://localhost:%s", port)
	if apiKey == "" {
		log.Print("GEMINI_API_KEY not set; using local rule-based executive analysis")
	}

	if err := http.ListenAndServe(":"+port, server.Handler()); err != nil {
		log.Fatal(err)
	}
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func seedUsageLogs() []feature.UsageLog {
	inputs := []feature.UsageLogInput{
		{Feature: "RelatorioExportacaoExcel", UserID: "123", LastAccess: "2026-01-10", TotalAccess: 2},
		{Feature: "RelatorioExportacaoExcel", UserID: "456", LastAccess: "2025-12-05", TotalAccess: 1},
		{Feature: "DashboardExecutivoV1", UserID: "777", LastAccess: "2026-06-14", TotalAccess: 62},
		{Feature: "DashboardExecutivoV1", UserID: "888", LastAccess: "2026-06-13", TotalAccess: 41},
		{Feature: "ImportadorCNABLegado", UserID: "321", LastAccess: "2026-02-01", TotalAccess: 6},
		{Feature: "ConsultaContratoAntigo", UserID: "222", LastAccess: "2025-09-18", TotalAccess: 1},
		{Feature: "ConsultaContratoAntigo", UserID: "333", LastAccess: "2025-09-20", TotalAccess: 1},
		{Feature: "WebhookPedidosV2", UserID: "svc-orders", LastAccess: time.Now().UTC().Format(time.DateOnly), TotalAccess: 450},
	}

	logs := make([]feature.UsageLog, 0, len(inputs))
	for _, input := range inputs {
		log, err := feature.NewUsageLog(input)
		if err != nil {
			panic(err)
		}
		logs = append(logs, log)
	}
	return logs
}

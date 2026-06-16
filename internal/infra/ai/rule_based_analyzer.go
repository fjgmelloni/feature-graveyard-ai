package ai

import (
	"context"
	"fmt"

	"feature-graveyard-ai/internal/application/graveyard"
	"feature-graveyard-ai/internal/domain/feature"
)

type RuleBasedAnalyzer struct{}

func NewRuleBasedAnalyzer() RuleBasedAnalyzer {
	return RuleBasedAnalyzer{}
}

func (RuleBasedAnalyzer) Analyze(_ context.Context, usage feature.AggregatedFeatureUsage, status feature.Status, risk feature.Risk) (graveyard.ExecutiveAnalysis, error) {
	switch status {
	case feature.StatusDeadFeature:
		return graveyard.ExecutiveAnalysis{
			Summary:            fmt.Sprintf("Funcionalidade com baixo uso nos ultimos %d dias. Pode ser candidata a remocao ou revisao.", usage.DaysSinceAccess),
			BusinessImpact:     "Reduz complexidade do sistema, custo de manutencao e ruido em discovery tecnico.",
			SuggestedAction:    "Mapear dependencias, avisar stakeholders e abrir plano de sunset controlado.",
			ExecutiveRationale: fmt.Sprintf("Apenas %d usuario(s), %d acesso(s) totais e %.2f acesso(s)/mes.", usage.UniqueUsers, usage.TotalAccess, usage.FrequencyPerMonth),
			GeneratedBy:        "rules",
		}, nil
	case feature.StatusAtRisk:
		return graveyard.ExecutiveAnalysis{
			Summary:            "Funcionalidade pouco usada, mas ainda com sinais de dependencia operacional.",
			BusinessImpact:     "Pode existir valor escondido em nichos de usuario ou rotinas periodicas.",
			SuggestedAction:    "Instrumentar eventos, entrevistar usuarios recentes e decidir entre melhorar, consolidar ou remover.",
			ExecutiveRationale: fmt.Sprintf("Ultimo acesso ha %d dias, com %d usuario(s) unico(s).", usage.DaysSinceAccess, usage.UniqueUsers),
			GeneratedBy:        "rules",
		}, nil
	default:
		return graveyard.ExecutiveAnalysis{
			Summary:            "Funcionalidade com uso suficiente para permanecer no backlog operacional.",
			BusinessImpact:     "A remocao pode gerar impacto direto em usuarios ativos e processos de negocio.",
			SuggestedAction:    "Manter e acompanhar tendencia de uso nos proximos ciclos.",
			ExecutiveRationale: fmt.Sprintf("Frequencia estimada de %.2f acesso(s)/mes.", usage.FrequencyPerMonth),
			GeneratedBy:        "rules",
		}, nil
	}
}

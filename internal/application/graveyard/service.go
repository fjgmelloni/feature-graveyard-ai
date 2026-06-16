package graveyard

import (
	"context"
	"fmt"
	"time"

	"feature-graveyard-ai/internal/domain/feature"
)

type UsageRepository interface {
	SaveMany(ctx context.Context, logs []feature.UsageLog) error
	List(ctx context.Context) ([]feature.UsageLog, error)
}

type ExecutiveAnalyzer interface {
	Analyze(ctx context.Context, usage feature.AggregatedFeatureUsage, status feature.Status, risk feature.Risk) (ExecutiveAnalysis, error)
}

type ExecutiveAnalysis struct {
	Summary            string
	BusinessImpact     string
	SuggestedAction    string
	ExecutiveRationale string
	GeneratedBy        string
}

type Service struct {
	repository UsageRepository
	analyzer   ExecutiveAnalyzer
	now        func() time.Time
}

func NewService(repository UsageRepository, analyzer ExecutiveAnalyzer) Service {
	return Service{
		repository: repository,
		analyzer:   analyzer,
		now:        time.Now,
	}
}

func (s Service) WithClock(now func() time.Time) Service {
	s.now = now
	return s
}

func (s Service) Ingest(ctx context.Context, inputs []feature.UsageLogInput) ([]feature.UsageLog, error) {
	logs := make([]feature.UsageLog, 0, len(inputs))
	for index, input := range inputs {
		log, err := feature.NewUsageLog(input)
		if err != nil {
			return nil, fmt.Errorf("usage log %d: %w", index, err)
		}
		logs = append(logs, log)
	}

	if err := s.repository.SaveMany(ctx, logs); err != nil {
		return nil, err
	}

	return logs, nil
}

func (s Service) Report(ctx context.Context, windowDays int) (feature.GraveyardReport, error) {
	if windowDays <= 0 {
		windowDays = 180
	}

	logs, err := s.repository.List(ctx)
	if err != nil {
		return feature.GraveyardReport{}, err
	}

	now := s.now().UTC()
	aggregated := feature.AggregateUsage(logs, now)
	analyses := make([]feature.FeatureAnalysis, 0, len(aggregated))

	for _, usage := range aggregated {
		status, risk, confidence := feature.Classify(usage)
		executive, err := s.analyzer.Analyze(ctx, usage, status, risk)
		if err != nil {
			executive = fallbackExecutiveAnalysis(usage, status, risk)
		}

		analysis := feature.FeatureAnalysis{
			Feature:            usage.Feature,
			Status:             status,
			Risk:               risk,
			Summary:            executive.Summary,
			BusinessImpact:     executive.BusinessImpact,
			SuggestedAction:    executive.SuggestedAction,
			LastAccess:         usage.LastAccess.Format(time.DateOnly),
			UniqueUsers:        usage.UniqueUsers,
			TotalAccess:        usage.TotalAccess,
			FrequencyPerMonth:  usage.FrequencyPerMonth,
			DaysSinceAccess:    usage.DaysSinceAccess,
			Confidence:         confidence,
			GeneratedBy:        executive.GeneratedBy,
			ExecutiveRationale: executive.ExecutiveRationale,
		}
		analyses = append(analyses, analysis)
	}

	forgotten := feature.ForgottenModules(analyses)
	for index := range analyses {
		analyses[index].ForgottenModules = forgotten
	}

	report := feature.GraveyardReport{
		GeneratedAt:   now.Format(time.RFC3339),
		WindowDays:    windowDays,
		TotalFeatures: len(analyses),
		Analyses:      analyses,
	}

	for _, analysis := range analyses {
		switch analysis.Status {
		case feature.StatusDeadFeature:
			report.DeadFeatures++
			if analysis.Risk == feature.RiskLow || analysis.Risk == feature.RiskMedium {
				report.RemovalCandidates++
			}
		case feature.StatusAtRisk:
			report.AtRiskFeatures++
		case feature.StatusActive:
			report.ActiveFeatures++
		}
	}

	return report, nil
}

func fallbackExecutiveAnalysis(usage feature.AggregatedFeatureUsage, status feature.Status, risk feature.Risk) ExecutiveAnalysis {
	switch status {
	case feature.StatusDeadFeature:
		return ExecutiveAnalysis{
			Summary:            fmt.Sprintf("%s tem baixo uso e ficou %d dias sem acesso relevante.", usage.Feature, usage.DaysSinceAccess),
			BusinessImpact:     "Pode reduzir complexidade, custo de manutenção e superficie de regressao do produto.",
			SuggestedAction:    "Validar donos de negocio, anunciar sunset e remover apos janela de observacao.",
			ExecutiveRationale: "Classificacao baseada em recencia, usuarios unicos e volume total de acessos.",
			GeneratedBy:        "rules-fallback",
		}
	case feature.StatusAtRisk:
		return ExecutiveAnalysis{
			Summary:            fmt.Sprintf("%s apresenta sinais de queda de uso e merece revisao.", usage.Feature),
			BusinessImpact:     "Manter sem revisao pode sustentar complexidade com retorno incerto.",
			SuggestedAction:    "Conversar com usuarios-chave e medir por mais um ciclo antes de remover.",
			ExecutiveRationale: "Uso recente ou base de usuarios insuficiente para recomendar remocao direta.",
			GeneratedBy:        "rules-fallback",
		}
	default:
		return ExecutiveAnalysis{
			Summary:            fmt.Sprintf("%s ainda demonstra uso recorrente.", usage.Feature),
			BusinessImpact:     "Remover agora pode afetar usuarios ativos e gerar retrabalho operacional.",
			SuggestedAction:    "Manter, instrumentar melhor e revisar novamente no proximo ciclo.",
			ExecutiveRationale: "Volume, frequencia e recencia estao acima dos limiares de risco.",
			GeneratedBy:        "rules-fallback",
		}
	}
}

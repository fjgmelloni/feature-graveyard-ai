package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"feature-graveyard-ai/internal/application/graveyard"
	"feature-graveyard-ai/internal/domain/feature"
)

type GeminiAnalyzer struct {
	apiKey     string
	model      string
	httpClient *http.Client
	fallback   RuleBasedAnalyzer
}

func NewGeminiAnalyzer(apiKey, model string) GeminiAnalyzer {
	if strings.TrimSpace(model) == "" {
		model = "gemini-1.5-flash"
	}

	return GeminiAnalyzer{
		apiKey: strings.TrimSpace(apiKey),
		model:  model,
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
		fallback: NewRuleBasedAnalyzer(),
	}
}

func (g GeminiAnalyzer) Analyze(ctx context.Context, usage feature.AggregatedFeatureUsage, status feature.Status, risk feature.Risk) (graveyard.ExecutiveAnalysis, error) {
	if g.apiKey == "" {
		return g.fallback.Analyze(ctx, usage, status, risk)
	}

	prompt := fmt.Sprintf(`Voce e um especialista executivo em sistemas legados.
Analise a feature abaixo e responda somente JSON valido, sem markdown.

Feature: %s
Status calculado: %s
Risco de remocao calculado: %s
Ultimo acesso ha dias: %d
Usuarios unicos: %d
Acessos totais: %d
Frequencia mensal estimada: %.2f

Schema obrigatorio:
{
  "summary": "maximo 180 caracteres",
  "businessImpact": "maximo 220 caracteres",
  "suggestedAction": "maximo 220 caracteres",
  "executiveRationale": "maximo 220 caracteres"
}`, usage.Feature, status, risk, usage.DaysSinceAccess, usage.UniqueUsers, usage.TotalAccess, usage.FrequencyPerMonth)

	requestBody := geminiRequest{
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: prompt}}},
		},
		GenerationConfig: geminiGenerationConfig{
			Temperature:      0.25,
			ResponseMimeType: "application/json",
		},
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return graveyard.ExecutiveAnalysis{}, err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", g.model, g.apiKey)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return graveyard.ExecutiveAnalysis{}, err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := g.httpClient.Do(request)
	if err != nil {
		return graveyard.ExecutiveAnalysis{}, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return graveyard.ExecutiveAnalysis{}, fmt.Errorf("gemini returned status %d", response.StatusCode)
	}

	var decoded geminiResponse
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return graveyard.ExecutiveAnalysis{}, err
	}
	if len(decoded.Candidates) == 0 || len(decoded.Candidates[0].Content.Parts) == 0 {
		return graveyard.ExecutiveAnalysis{}, errors.New("gemini returned no candidates")
	}

	raw := decoded.Candidates[0].Content.Parts[0].Text
	var executive graveyard.ExecutiveAnalysis
	if err := json.Unmarshal([]byte(raw), &executive); err != nil {
		return graveyard.ExecutiveAnalysis{}, err
	}
	executive.GeneratedBy = "gemini"

	return executive, nil
}

type geminiRequest struct {
	Contents         []geminiContent        `json:"contents"`
	GenerationConfig geminiGenerationConfig `json:"generationConfig"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiGenerationConfig struct {
	Temperature      float64 `json:"temperature"`
	ResponseMimeType string  `json:"responseMimeType"`
}

type geminiResponse struct {
	Candidates []struct {
		Content geminiContent `json:"content"`
	} `json:"candidates"`
}

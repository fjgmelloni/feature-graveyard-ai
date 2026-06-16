package feature

import (
	"math"
	"sort"
	"time"
)

type Status string

const (
	StatusActive      Status = "ACTIVE"
	StatusAtRisk      Status = "AT_RISK"
	StatusDeadFeature Status = "DEAD_FEATURE"
)

type Risk string

const (
	RiskLow      Risk = "LOW"
	RiskMedium   Risk = "MEDIUM"
	RiskHigh     Risk = "HIGH"
	RiskCritical Risk = "CRITICAL"
)

type AggregatedFeatureUsage struct {
	Feature           string
	LastAccess        time.Time
	UniqueUsers       int
	TotalAccess       int
	FrequencyPerMonth float64
	DaysSinceAccess   int
}

type FeatureAnalysis struct {
	Feature            string   `json:"feature"`
	Status             Status   `json:"status"`
	Risk               Risk     `json:"risk"`
	Summary            string   `json:"summary"`
	BusinessImpact     string   `json:"businessImpact"`
	SuggestedAction    string   `json:"suggestedAction"`
	LastAccess         string   `json:"lastAccess"`
	UniqueUsers        int      `json:"uniqueUsers"`
	TotalAccess        int      `json:"totalAccess"`
	FrequencyPerMonth  float64  `json:"frequencyPerMonth"`
	DaysSinceAccess    int      `json:"daysSinceAccess"`
	ForgottenModules   []string `json:"forgottenModules"`
	Confidence         float64  `json:"confidence"`
	GeneratedBy        string   `json:"generatedBy"`
	ExecutiveRationale string   `json:"executiveRationale"`
}

type GraveyardReport struct {
	GeneratedAt       string            `json:"generatedAt"`
	WindowDays        int               `json:"windowDays"`
	TotalFeatures     int               `json:"totalFeatures"`
	DeadFeatures      int               `json:"deadFeatures"`
	AtRiskFeatures    int               `json:"atRiskFeatures"`
	ActiveFeatures    int               `json:"activeFeatures"`
	RemovalCandidates int               `json:"removalCandidates"`
	Analyses          []FeatureAnalysis `json:"analyses"`
}

func AggregateUsage(logs []UsageLog, now time.Time) []AggregatedFeatureUsage {
	byFeature := map[string]*AggregatedFeatureUsage{}
	usersByFeature := map[string]map[string]struct{}{}

	for _, log := range logs {
		current, exists := byFeature[log.Feature]
		if !exists {
			current = &AggregatedFeatureUsage{Feature: log.Feature, LastAccess: log.LastAccess}
			byFeature[log.Feature] = current
			usersByFeature[log.Feature] = map[string]struct{}{}
		}

		if log.LastAccess.After(current.LastAccess) {
			current.LastAccess = log.LastAccess
		}
		current.TotalAccess += log.TotalAccess
		usersByFeature[log.Feature][log.UserID] = struct{}{}
	}

	result := make([]AggregatedFeatureUsage, 0, len(byFeature))
	for featureName, usage := range byFeature {
		usage.UniqueUsers = len(usersByFeature[featureName])
		usage.DaysSinceAccess = int(now.Sub(usage.LastAccess).Hours() / 24)
		if usage.DaysSinceAccess < 0 {
			usage.DaysSinceAccess = 0
		}

		months := math.Max(float64(usage.DaysSinceAccess)/30.0, 1)
		usage.FrequencyPerMonth = round2(float64(usage.TotalAccess) / months)
		result = append(result, *usage)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Feature < result[j].Feature
	})

	return result
}

func Classify(usage AggregatedFeatureUsage) (Status, Risk, float64) {
	switch {
	case usage.DaysSinceAccess >= 180 && usage.UniqueUsers <= 2 && usage.TotalAccess <= 5:
		return StatusDeadFeature, RiskLow, 0.94
	case usage.DaysSinceAccess >= 180 && usage.FrequencyPerMonth < 2:
		return StatusDeadFeature, RiskMedium, 0.88
	case usage.DaysSinceAccess >= 90 || usage.FrequencyPerMonth < 3 || usage.UniqueUsers <= 3:
		return StatusAtRisk, RiskHigh, 0.79
	default:
		return StatusActive, RiskCritical, 0.86
	}
}

func ForgottenModules(analyses []FeatureAnalysis) []string {
	modules := make([]string, 0)
	for _, analysis := range analyses {
		if analysis.Status == StatusDeadFeature || analysis.Status == StatusAtRisk {
			modules = append(modules, analysis.Feature)
		}
	}
	sort.Strings(modules)
	return modules
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}

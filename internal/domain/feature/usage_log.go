package feature

import (
	"errors"
	"strings"
	"time"
)

type UsageLog struct {
	Feature     string    `json:"feature"`
	UserID      string    `json:"userId"`
	LastAccess  time.Time `json:"lastAccess"`
	TotalAccess int       `json:"totalAccess"`
}

type UsageLogInput struct {
	Feature     string `json:"feature"`
	UserID      string `json:"userId"`
	LastAccess  string `json:"lastAccess"`
	TotalAccess int    `json:"totalAccess"`
}

func NewUsageLog(input UsageLogInput) (UsageLog, error) {
	featureName := strings.TrimSpace(input.Feature)
	userID := strings.TrimSpace(input.UserID)

	if featureName == "" {
		return UsageLog{}, errors.New("feature is required")
	}
	if userID == "" {
		return UsageLog{}, errors.New("userId is required")
	}
	if input.TotalAccess < 0 {
		return UsageLog{}, errors.New("totalAccess cannot be negative")
	}

	lastAccess, err := time.Parse(time.DateOnly, input.LastAccess)
	if err != nil {
		return UsageLog{}, errors.New("lastAccess must use YYYY-MM-DD format")
	}

	return UsageLog{
		Feature:     featureName,
		UserID:      userID,
		LastAccess:  lastAccess,
		TotalAccess: input.TotalAccess,
	}, nil
}

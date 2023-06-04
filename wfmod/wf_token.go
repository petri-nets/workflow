package wfmod

import "time"

type TokenStatusType string

// TokenStatus
const (
	TokenStatusFree      TokenStatusType = "FREE"
	TokenStatusLock      TokenStatusType = "LOCK"
	TokenStatusConsume   TokenStatusType = "CONS"
	TokenStatusCanceller TokenStatusType = "CANC"
)

// WfToken Model
type WfToken struct {
	BaseModel
	AppID         int             `json:"app_id"`
	CaseID        int             `json:"case_id"`
	PlaceID       int             `json:"place_id"`
	WorkflowID    int             `json:"workflow_id"`
	Context       string          `json:"context"`
	Status        TokenStatusType `json:"status"`
	CancelledDate time.Time       `json:"cancelled_date"`
	ConsumedDate  time.Time       `json:"consumed_date"`
}

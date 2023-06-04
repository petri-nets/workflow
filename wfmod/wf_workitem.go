package wfmod

import "time"

type WorkitemStatusType string

// WorkitemStatus
const (
	WorkitemStatusEnabled    WorkitemStatusType = "EN"
	WorkitemStatusInProgress WorkitemStatusType = "IP"
	WorkitemStatusCancelled  WorkitemStatusType = "CA"
	WorkitemStatusFinished   WorkitemStatusType = "FI"
)

// WfWorkitem model
type WfWorkitem struct {
	BaseModel
	AppID             int                   `json:"app_id"`
	CaseID            int                   `json:"case_id"`
	WorkflowID        int                   `json:"workflow_id"`
	TransitionID      int                   `json:"transition_id"`
	TransitionTrigger TransitionTriggerType `json:"transition_trigger"`
	JobID             int                   `json:"job_id"`
	Context           string                `json:"context"`
	Status            WorkitemStatusType    `json:"status"`
	EnabledDate       time.Time             `json:"enabled_date"`
	CancelledDate     time.Time             `json:"cancelled_date"`
	FinishedDate      time.Time             `json:"finished_date"`
	Deadline          time.Time             `json:"deadline"`
	RoleID            int                   `json:"role_id"`
	User              string                `json:"user"`
}

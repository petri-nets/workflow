package wfmod

import "time"

// CaseStatusType case status type
type CaseStatusType string

const (
	CaseStatusOP CaseStatusType = "OP" // open
	CaseStatusCL CaseStatusType = "CL" // closed
	CaseStatusSU CaseStatusType = "SU" // suspended 挂起
	CaseStatusCA CaseStatusType = "CA" // cancelled
)

// WfCase Model
type WfCase struct {
	BaseModel
	AppID      int            `json:"app_id"`
	WorkflowID int            `json:"workflow_id"`
	Context    string         `json:"context"`
	Status     CaseStatusType `json:"status"`
	StartAt    time.Time      `json:"start_at"`
	EndAt      time.Time      `json:"end_at"`
	CreatedBy  string         `json:"created_by"`
	UpdatedBy  string         `json:"updated_by"`
}

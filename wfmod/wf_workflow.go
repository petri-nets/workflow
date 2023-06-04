package wfmod

import "time"

type WfContextType map[string]interface{}
type WorkflowValidType string

const (
	WorkflowValidY WorkflowValidType = "Y"
	WorkflowValidN WorkflowValidType = "N"
)

// WfWorkflow 完整的flow结构
type WfWorkflow struct {
	BaseModel
	AppID      int               `json:"app_id"`
	Name       string            `json:"name"`
	Desc       string            `json:"desc"`
	StartJobID int               `json:"start_job_id"`
	IsValid    WorkflowValidType `json:"is_valid"`
	Errors     string            `json:"errors"`
	StartAt    time.Time         `json:"start_at"`
	EndAt      time.Time         `json:"end_at"`
	CreatedBy  string            `json:"created_by"`
	UpdatedBy  string            `json:"updated_by"`
}

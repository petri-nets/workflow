package wfmod

type TransitionTriggerType string

// TrrasitionTrigger
const (
	TransitionTriggerUser TransitionTriggerType = "USER"
	TransitionTriggerAuto TransitionTriggerType = "AUTO"
	TransitionTriggerMsg  TransitionTriggerType = "MSG"
	TransitionTriggerTime TransitionTriggerType = "TIME"
)

// TransitionTriggerMap with all transition trigger
var TransitionTriggerMap = map[TransitionTriggerType]TransitionTriggerType{
	TransitionTriggerAuto: TransitionTriggerAuto,
	TransitionTriggerMsg:  TransitionTriggerMsg,
	TransitionTriggerTime: TransitionTriggerTime,
	TransitionTriggerUser: TransitionTriggerUser,
}

// WfTransition Model
type WfTransition struct {
	BaseModel
	AppID      int                   `json:"app_id"`
	WorkflowID int                   `json:"workflow_id"`
	Name       string                `json:"name"`
	Desc       string                `json:"desc"`
	Trigger    TransitionTriggerType `json:"trigger"`
	TimeLimit  int                   `json:"time_limit"` // 单位：分钟
	JobID      int                   `json:"job_id"`
	RoleID     int                   `json:"role_id"`
}

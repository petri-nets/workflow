package flow

import (
	"github.com/petri-nets/workflow/wfmod"
)

var flowDao wfmod.FlowDao
var executorHandle func(cs *Case, workitem *Workitem) (wait bool, err error)

// Workflow 完整的flow结构
type Workflow struct {
	WfWorkflow wfmod.WfWorkflow
	Case       Case
}

// Init init
func Init(fd wfmod.FlowDao, executor func(*Case, *Workitem) (wait bool, err error)) {
	flowDao = fd
	executorHandle = executor
}

// GetWorkflowsByStartJob get workflows by start job
func GetWorkflowsByStartJob(appID, startJobID int) []Workflow {
	wfWorkflows := flowDao.GetWorkflowsByStartJob(appID, startJobID)
	workflows := []Workflow{}
	for _, wfWorkflow := range wfWorkflows {
		workflows = append(workflows, Workflow{WfWorkflow: wfWorkflow})
	}
	return workflows
}

// Start 启动工作流新的case
func (w *Workflow) Start(wfContext wfmod.WfContextType, operator string) {
	w.Case = NewCase(wfContext, w, operator)

	place := GetStartPlace(&w.Case)

	token := NewToken(&w.Case, w.Case.WfCase.Context, place.WfPlace.ID)

	place.ReceiveToken(&w.Case, &token)
}

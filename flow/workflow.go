package flow

import (
	"github.com/petri-nets/workflow/wfmod"
)

var flowDao wfmod.FlowDao
var executorHandle func(cs *Case, workitem *Workitem) (wait bool, err error)
var customizeValidator func(tran *wfmod.WfTransition) error

// Workflow 完整的flow结构
type Workflow struct {
	WfWorkflow wfmod.WfWorkflow
	Case       Case
}

// RegFlowDao register flow dao
func RegFlowDao(fd wfmod.FlowDao) {
	flowDao = fd
}

// RegExecutorHandle register executor handle
func RegExecutorHandle(_executorHandle func(*Case, *Workitem) (wait bool, err error)) {
	executorHandle = _executorHandle
}

// RegCustomerValidator register customer validator
func RegCustomerValidator(_customizeValidator func(*wfmod.WfTransition) error) {
	customizeValidator = _customizeValidator
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

package flow

import (
	"github.com/petri-nets/workflow/wfmod"
)

var flowDao wfmod.FlowDao
var executorHandle func(cs *Case, workitem *Workitem) (wait bool, err error)
var customizeValidator func(tran *wfmod.WfTransition) error

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

// StartWorkflow 启动工作流新的case
func StartWorkflow(appID, startJobID int, wfContext wfmod.WfContextType, operator string) {

	workflows := GetWorkflowsByStartJob(appID, startJobID)

	for _, workflow := range workflows {
		// todo 获取用用户信息
		go workflow.Start(wfContext, operator)
	}
}

// ContinueCase continue case
func ContinueCase(appID, jobID int, jobOutput wfmod.WfContextType) {
	transitions := GetTransitionsByWorkitemJob(appID, jobID)

	for _, transition := range transitions {
		// 设置workitem为完成状态
		go transition.WorkitemFinishedHandle(jobOutput)
	}
}

// Workflow 完整的flow结构
type Workflow struct {
	WfWorkflow wfmod.WfWorkflow
	Case       Case
}

// Start 启动工作流
func (w *Workflow) Start(wfContext wfmod.WfContextType, operator string) {

	w.Case = NewCase(wfContext, w, operator)

	place := GetStartPlace(&w.Case)

	token := NewToken(&w.Case, w.Case.WfCase.Context, place.WfPlace.ID)

	place.ReceiveToken(&w.Case, &token)
}

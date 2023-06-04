package flow

import (
	"encoding/json"

	"github.com/petri-nets/workflow/wfmod"
)

// Case Model
type Case struct {
	WfCase    wfmod.WfCase
	WfContext wfmod.WfContextType
}

// NewCase 创建新case，case的起点
func NewCase(wfContext wfmod.WfContextType, workflow *Workflow, operator string) Case {

	cs := Case{
		WfCase:    wfmod.WfCase{},
		WfContext: wfContext,
	}

	wfContextStr, err := json.Marshal(wfContext)
	if err != nil {
		// log.WithFields(log.Fields{
		// 	"workflow_id": workflow.WfWorkflow.ID,
		// 	"data":        wfContext,
		// }).Error("[Workflow][NewCase] encode context error")
	}

	cs.WfCase = wfmod.WfCase{
		AppID:      workflow.WfWorkflow.AppID,
		WorkflowID: workflow.WfWorkflow.ID,
		Context:    string(wfContextStr),
		Status:     wfmod.CaseStatusOP,
		StartAt:    workflow.WfWorkflow.StartAt,
		EndAt:      workflow.WfWorkflow.EndAt,
		CreatedBy:  operator,
		UpdatedBy:  operator,
	}

	cs.createWfCase()

	return cs
}

// GetCasesByIDList get cases by id list
func GetCasesByIDList(appID int, caseIDList []int) []Case {
	wfCases := flowDao.GetCasesByIDList(appID, caseIDList)

	cases := []Case{}
	for _, wfCase := range wfCases {
		wfContext := wfmod.WfContextType{}
		json.Unmarshal([]byte(wfCase.Context), &wfContext)
		cases = append(cases, Case{WfCase: wfCase, WfContext: wfContext})
	}
	return cases
}

// Create 创建新case
func (c *Case) createWfCase() error {
	flowDao.SaveCase(&c.WfCase)
	return nil
}

// Suspended 挂起case
func (c *Case) Suspended() error {
	c.WfCase.Status = wfmod.CaseStatusSU
	flowDao.SaveCase(&c.WfCase)
	return nil
}

// Close 关闭case
func (c *Case) Close() error {
	c.WfCase.Status = wfmod.CaseStatusCL
	flowDao.SaveCase(&c.WfCase)
	return nil
}

// Cancel 取消case
func (c *Case) Cancel() error {
	c.WfCase.Status = wfmod.CaseStatusCA
	flowDao.SaveCase(&c.WfCase)
	return nil
}

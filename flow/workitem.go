package flow

import (
	"encoding/json"

	"github.com/petri-nets/workflow/wfmod"
)

// Workitem workitem
type Workitem struct {
	WfWorkitem wfmod.WfWorkitem
}

// CreateWorkitem create workitem
func CreateWorkitem(t *Transition, token *Token) Workitem {
	params, _ := json.Marshal(token.WfContext)

	workitem := Workitem{}

	workitem.WfWorkitem = wfmod.WfWorkitem{
		AppID:             t.WfTransition.AppID,
		WorkflowID:        t.WfTransition.WorkflowID,
		CaseID:            t.Case.WfCase.ID,
		TransitionID:      t.WfTransition.ID,
		TransitionTrigger: t.WfTransition.Trigger,
		JobID:             t.WfTransition.JobID,
		RoleID:            t.WfTransition.RoleID,
		Context:           string(params),
		Status:            wfmod.WorkitemStatusEnabled,
	}

	flowDao.SaveWorkitem(&workitem.WfWorkitem)

	return workitem
}

// GetOpeningWorkitemsByJobID get opening workitems by job id
func GetOpeningWorkitemsByJobID(appID, jobID int) []Workitem {
	wfWorkitems := flowDao.GetOpeningWorkitemsByJobID(appID, jobID)

	workitems := []Workitem{}
	for _, wfWorkitem := range wfWorkitems {
		workitems = append(workitems, Workitem{WfWorkitem: wfWorkitem})
	}
	return workitems
}

// do work item 开始执行
func (w *Workitem) Do(cs *Case) (isContinue bool) {
	// 修改状态，防止重复执行
	w.setInProgress()

	if w.WfWorkitem.JobID != 0 {
		wait, _ := executorHandle(cs, w)
		return !wait
	} else {
		return true
	}

}

// MergeCaseContext merge case context
func (w *Workitem) MergeCaseContext(cs *Case) wfmod.WfContextType {
	wfContext := wfmod.WfContextType{}
	err := json.Unmarshal([]byte(w.WfWorkitem.Context), &wfContext)
	if err != nil {
		return cs.WfContext
	}

	for k, v := range wfContext {
		cs.WfContext[k] = v
	}

	return cs.WfContext
}

// MergeOutputContext merge output context
func (w *Workitem) MergeOutputContext(output wfmod.WfContextType) wfmod.WfContextType {
	wfContext := wfmod.WfContextType{}
	err := json.Unmarshal([]byte(w.WfWorkitem.Context), &wfContext)
	if err != nil {
		// log.Error("[workflow][workitem][MergeOutputContext] decode workitem context error")
	}

	for k, v := range output {
		wfContext[k] = v
	}

	return wfContext
}

// GetSystemContext get system context
func (w *Workitem) GetSystemContext(cs *Case) wfmod.WfContextType {
	return wfmod.WfContextType{
		"CASE_ID":              cs.WfCase.ID,
		"APP_ID":               cs.WfCase.AppID,
		"WORKFLOW_ID":          cs.WfCase.WorkflowID,
		"TRANSITION_ID":        w.WfWorkitem.TransitionID,
		"TRIGGER_TYPE":         w.WfWorkitem.TransitionTrigger,
		"WORKITEM_ID":          w.WfWorkitem.ID,
		"WORKITEM_STATUS":      w.WfWorkitem.Status,
		"WORKITEM_ENABLED_AT":  w.WfWorkitem.EnabledDate,
		"WORKITEM_DEADLINE":    w.WfWorkitem.Deadline,
		"WORKITEM_FINISHED_AT": w.WfWorkitem.FinishedDate,
	}
}

// HasFinished 判断workitem 是否已经完成
func (w *Workitem) HasFinished() bool {
	wfWorkitems := flowDao.GetWorkitemsByIDList(w.WfWorkitem.AppID, []int{w.WfWorkitem.ID})
	if len(wfWorkitems) > 0 {
		w.WfWorkitem = wfWorkitems[0]
	}
	return w.WfWorkitem.Status == wfmod.WorkitemStatusFinished ||
		w.WfWorkitem.Status == wfmod.WorkitemStatusCancelled
}

// setInProgress processing Workitem
func (w *Workitem) setInProgress() error {
	w.WfWorkitem.Status = wfmod.WorkitemStatusInProgress
	flowDao.SaveWorkitem(&w.WfWorkitem)
	return nil
}

// finish Workitem
func (w *Workitem) finish() error {
	w.WfWorkitem.Status = wfmod.WorkitemStatusFinished
	flowDao.SaveWorkitem(&w.WfWorkitem)
	return nil
}

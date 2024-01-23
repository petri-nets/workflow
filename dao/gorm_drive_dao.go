package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/petri-nets/workflow/wfmod"
)

// GormDrive flow数据库获取接口实现
type GormDrive struct {
	DB *gorm.DB
	tx *gorm.DB
}

// GetWorkflowsByStartJob get workflow by start job id
func (f *GormDrive) GetWorkflowsByStartJob(appID, startJobID int) []wfmod.WfWorkflow {
	list := []wfmod.WfWorkflow{}
	f.DB.Where("app_id = ? AND start_job_id = ? AND is_valid = ?",
		appID, startJobID, wfmod.WorkflowValidY).Find(&list)
	return list
}

// SaveWorkflow create workflow or update workflow
func (f *GormDrive) SaveWorkflow(w *wfmod.WfWorkflow) error {
	return f.DB.Save(w).Error
}

// GetPlaceArcs get place arcs
func (f *GormDrive) GetPlaceArcs(p *wfmod.WfPlace, direct wfmod.ArcDirectionType) []wfmod.WfArc {
	list := []wfmod.WfArc{}
	f.DB.Where("app_id = ? AND workflow_id = ? AND place_id = ? AND direction = ?",
		p.AppID, p.WorkflowID, p.ID, direct).Find(&list)
	return list
}

// GetTransitionArcs get transition arcs
func (f *GormDrive) GetTransitionArcs(t *wfmod.WfTransition, direct wfmod.ArcDirectionType) []wfmod.WfArc {
	list := []wfmod.WfArc{}
	f.DB.Where("app_id = ? AND workflow_id = ? AND transition_id = ? AND direction = ?",
		t.AppID, t.WorkflowID, t.ID, direct).Find(&list)
	return list
}

// SaveArc save arc or update when give id
func (f *GormDrive) SaveArc(a *wfmod.WfArc) error {
	return f.DB.Save(a).Error
}

// GetCasesByIDList get cases by id list
func (f *GormDrive) GetCasesByIDList(appID int, caseIDList []int) []wfmod.WfCase {
	list := []wfmod.WfCase{}
	f.DB.Where("app_id = ? AND id in (?)", appID, caseIDList).Find(&list)
	return list
}

// SaveCase save case or update when give id
func (f *GormDrive) SaveCase(c *wfmod.WfCase) error {
	return f.DB.Save(c).Error
}

// GetPlacesByType get places by place type
func (f *GormDrive) GetPlacesByType(appID int, workflowID int, placeType wfmod.PlaceTypeType) []wfmod.WfPlace {
	list := []wfmod.WfPlace{}
	f.DB.Where("app_id = ? AND workflow_id = ? AND type = ?", appID, workflowID, placeType).Find(&list)
	return list
}

// GetPlacesByIDList get places by id list
func (f *GormDrive) GetPlacesByIDList(appID int, idList []int) []wfmod.WfPlace {
	list := []wfmod.WfPlace{}
	f.DB.Where("app_id = ? AND id IN (?)", appID, idList).Find(&list)
	return list
}

// SavePlace save place or update when give id
func (f *GormDrive) SavePlace(p *wfmod.WfPlace) error {
	return f.DB.Save(p).Error
}

// GetTransitionsByIDList get transitions by id list
func (f *GormDrive) GetTransitionsByIDList(appID int, idList []int) []wfmod.WfTransition {
	list := []wfmod.WfTransition{}
	f.DB.Where("app_id = ? AND id IN (?)", appID, idList).Find(&list)
	return list
}

// SaveTransition save transition or update transition when give id
func (f *GormDrive) SaveTransition(p *wfmod.WfTransition) error {
	return f.DB.Save(p).Error
}

// GetTokenByIDList get tokens by id list
func (f *GormDrive) GetTokenByIDList(appID int, idList []int) []wfmod.WfToken {
	list := []wfmod.WfToken{}
	f.DB.Where("app_id = ? AND id IN (?)", appID, idList).Find(&list)
	return list
}

// GetTokensByPlaces get token by place id list
func (f *GormDrive) GetTokensByPlaces(cs *wfmod.WfCase, placeIDList []int,
	status wfmod.TokenStatusType) []wfmod.WfToken {
	list := []wfmod.WfToken{}
	f.DB.Where("app_id = ? AND workflow_id = ? AND case_id = ? AND place_id IN (?) AND status = ?",
		cs.AppID, cs.WorkflowID, cs.ID, placeIDList, status).Find(&list)
	return list
}

// SaveToken save token or update token when give id
func (f *GormDrive) SaveToken(t *wfmod.WfToken) error {
	return f.DB.Save(t).Error
}

// GetOpeningWorkitemsByJobID get opening workitems by job id
func (f *GormDrive) GetOpeningWorkitemsByJobID(appID, jobID int) []wfmod.WfWorkitem {
	list := []wfmod.WfWorkitem{}
	f.DB.Where("app_id = ? AND job_id = ? AND status IN (?)",
		appID, jobID, []wfmod.WorkitemStatusType{wfmod.WorkitemStatusEnabled, wfmod.WorkitemStatusInProgress},
	).Find(&list)
	return list
}

// GetWorkitemsByIDList get workitems by id list
func (f *GormDrive) GetWorkitemsByIDList(appID int, idList []int) []wfmod.WfWorkitem {
	list := []wfmod.WfWorkitem{}
	f.DB.Where("app_id = ? AND id IN (?)", appID, idList).Find(&list)
	return list
}

// SaveWorkitem save workitem or update workitem when give id
func (f *GormDrive) SaveWorkitem(w *wfmod.WfWorkitem) error {
	return f.DB.Save(w).Error
}

// BeginTransaction begin a transaction
func (f *GormDrive) BeginTransaction() {
	f.tx = f.DB.Begin()
}

// CommitTransaction commit a transaction
func (f *GormDrive) CommitTransaction() {
	f.tx.Commit()
}

// RollbackTransaction commit a transaction
func (f *GormDrive) RollbackTransaction() {
	f.tx.Rollback()
}

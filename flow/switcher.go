package flow

import (
	"fmt"
	"time"

	"github.com/petri-nets/workflow/wfmod"
)

// place
type place struct {
	Name      string              `json:"name"`
	Type      wfmod.PlaceTypeType `json:"type"`
	condition string
}

// workflow
type workflow struct {
	ID         int       `json:"id"`
	AppID      int       `json:"app_id"`
	Name       string    `json:"name"`
	Desc       string    `json:"desc"`
	StartJobID int       `json:"start_job_id"`
	StartAt    time.Time `json:"start_date"`
	EndAt      time.Time `json:"end_date"`
	Operator   string    `json:"operator"`
}

// arc
type arc struct {
	Place     string            `json:"place_name"`
	Type      wfmod.ArcTypeType `json:"type"`
	Condition string            `json:"condition"`
}

// transition
type transition struct {
	Name      string                      `json:"name"`
	Desc      string                      `json:"desc"`
	Trigger   wfmod.TransitionTriggerType `json:"trigger"`
	JobID     int                         `json:"job_id"`
	TimeLimit int                         `json:"time_limit"` // 单位：分钟
	RoleID    int                         `json:"role_id"`
	In        []arc                       `json:"in"`
	Out       []arc                       `json:"out"`
}

// Switcher 编排文件数据
type Switcher struct {
	Workflow    workflow
	Places      []place
	Transitions []transition
}

// Save 保存工作流
func (s *Switcher) Save() error {
	if err := s.Validate(); err != nil {
		return err
	}

	flowDao.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			flowDao.RollbackTransaction() // 出现panic时，回滚事务
		}
	}()

	workflow, err := s.SaveWfWorkflow()
	if err != nil {
		flowDao.RollbackTransaction()
		return err
	}

	placeMap, err := s.SavePlaces(workflow)
	if err != nil {
		flowDao.RollbackTransaction()
		return err
	}

	if err := s.SaveTransitions(workflow, placeMap); err != nil {
		flowDao.RollbackTransaction()
		return err
	}

	flowDao.CommitTransaction()

	return nil
}

// Validate validate workflow
func (s *Switcher) Validate() error {
	return nil
}

// SaveWfWorkflow save workflow
func (s *Switcher) SaveWfWorkflow() (*wfmod.WfWorkflow, error) {
	wfWorkflow := wfmod.WfWorkflow{
		AppID:      s.Workflow.AppID,
		Name:       s.Workflow.Name,
		Desc:       s.Workflow.Desc,
		IsValid:    wfmod.WorkflowValidY,
		StartJobID: s.Workflow.StartJobID,
		StartAt:    s.Workflow.StartAt,
		EndAt:      s.Workflow.EndAt,
		// UpdatedBy:   s.Workflow.Operator,
	}
	if s.Workflow.ID == 0 {
		wfWorkflow.CreatedBy = wfWorkflow.UpdatedBy
	}
	err := flowDao.SaveWorkflow(&wfWorkflow)
	return &wfWorkflow, err
}

// SavePlaces save places
func (s *Switcher) SavePlaces(wfWorkflow *wfmod.WfWorkflow) (map[string]*wfmod.WfPlace, error) {
	wfPlaces := make(map[string]*wfmod.WfPlace)
	for _, place := range s.Places {
		wfPlace := wfmod.WfPlace{
			AppID:      wfWorkflow.AppID,
			WorkflowID: wfWorkflow.ID,
			Type:       place.Type,
			Name:       place.Name,
		}
		err := flowDao.SavePlace(&wfPlace)
		if err != nil {
			return wfPlaces, err
		}
		wfPlaces[place.Name] = &wfPlace
	}
	return wfPlaces, nil
}

// SaveTransitions save transitions
func (s *Switcher) SaveTransitions(wfWorkflow *wfmod.WfWorkflow, wfPlaceMap map[string]*wfmod.WfPlace) error {
	for _, transition := range s.Transitions {

		if transition.Trigger == "" {
			transition.Trigger = wfmod.TransitionTriggerAuto
		}

		wfTransition := wfmod.WfTransition{
			AppID:      wfWorkflow.AppID,
			WorkflowID: wfWorkflow.ID,
			Name:       transition.Name,
			Desc:       transition.Desc,
			Trigger:    transition.Trigger,
			JobID:      transition.JobID,
			TimeLimit:  transition.TimeLimit,
			RoleID:     transition.RoleID,
		}

		if customizeValidator != nil {
			if err := customizeValidator(&wfTransition); err != nil {
				return fmt.Errorf("transition check error:%v", err)
			}
		}

		if err := flowDao.SaveTransition(&wfTransition); err != nil {
			return fmt.Errorf("save transition error:%v", err)
		}

		if wfTransition.ID > 0 {
			if err := s.SaveArcs(&wfTransition, transition.In, wfPlaceMap, wfmod.ArcDirectionIn); err != nil {
				return err
			}

			if err := s.SaveArcs(&wfTransition, transition.Out, wfPlaceMap, wfmod.ArcDirectionOut); err != nil {
				return err
			}
		}
	}

	return nil
}

// SaveArcs save arcs
func (s *Switcher) SaveArcs(wfTransition *wfmod.WfTransition, inArcs []arc,
	wfPlaceMap map[string]*wfmod.WfPlace, arcDir wfmod.ArcDirectionType) error {

	for _, arc := range inArcs {
		wfPlace, ok := wfPlaceMap[arc.Place]
		if !ok {
			continue
		}

		wfArc := wfmod.WfArc{
			AppID:        wfTransition.AppID,
			WorkflowID:   wfTransition.WorkflowID,
			TransitionID: wfTransition.ID,
			PlaceID:      wfPlace.ID,
			Direction:    arcDir,
			Type:         arc.Type,
			Condition:    arc.Condition,
		}

		if err := flowDao.SaveArc(&wfArc); err != nil {
			return fmt.Errorf("arc[%#v] save error:%v", wfArc, err)
		}
	}

	return nil
}

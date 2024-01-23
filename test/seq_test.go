package main

// package test

// import (
// 	"encoding/json"
// 	"testing"

// 	"github.com/go-playground/assert/v2"
// 	"github.com/golang/mock/gomock"
// 	"github.com/petri-nets/workflow/flow"
// 	"github.com/petri-nets/workflow/test/wfmock"
// 	"github.com/petri-nets/workflow/wfmod"
// )

// type seqMocker struct {
// 	mockCtrl *gomock.Controller
// 	t        *testing.T
// }

// func newSeqMocker(t *testing.T) *seqMocker {
// 	return &seqMocker{
// 		t:        t,
// 		mockCtrl: gomock.NewController(t),
// 	}
// }

// func (s *seqMocker) do(wfCtx wfmod.WfContextType) *wfmock.MockFlowDao {
// 	mocker := wfmock.NewMockFlowDao(s.mockCtrl)
// 	mocker.EXPECT().SaveCase(gomock.Any()).AnyTimes().Return(nil)

// 	gomock.InOrder(
// 		mocker.EXPECT().SaveToken(gomock.Any()).Do(func(token *wfmod.WfToken) {
// 			wfCtxStr, _ := json.Marshal(&wfCtx)
// 			assert.Equal(s.t, token.AppID, 1)
// 			assert.Equal(s.t, token.WorkflowID, 1)
// 			assert.Equal(s.t, token.PlaceID, 1)
// 			assert.Equal(s.t, token.Context, string(wfCtxStr))
// 			assert.Equal(s.t, token.Status, wfmod.TokenStatusFree)
// 		}),
// 	)

// 	mocker.EXPECT().GetPlacesByType(1, 1, wfmod.StartPlaceType).AnyTimes().Return([]wfmod.WfPlace{{
// 		BaseModel: wfmod.BaseModel{ID: 1}, AppID: 1, WorkflowID: 1,
// 		Type: wfmod.StartPlaceType,
// 		Desc: "start place",
// 		Name: "place-1",
// 	}})

// 	gomock.InOrder(
// 		mocker.EXPECT().GetPlaceArcs(gomock.Any(), wfmod.ArcDirectionIn).Return([]wfmod.WfArc{{
// 			BaseModel: wfmod.BaseModel{ID: 1},
// 			AppID:     1, WorkflowID: 1, PlaceID: 1, TransitionID: 1, Direction: wfmod.ArcDirectionIn,
// 			Type: wfmod.ArcSEQ,
// 		}}),
// 		mocker.EXPECT().GetPlaceArcs(gomock.Any(), wfmod.ArcDirectionIn).Return([]wfmod.WfArc{{
// 			BaseModel: wfmod.BaseModel{ID: 2},
// 			AppID:     1, WorkflowID: 1, PlaceID: 2, TransitionID: 2, Direction: wfmod.ArcDirectionIn,
// 			Type: wfmod.ArcSEQ,
// 		}}),
// 	)

// 	mocker.EXPECT().GetTransitionsByIDList(1, []int{1}).Return([]wfmod.WfTransition{{
// 		BaseModel: wfmod.BaseModel{ID: 1}, AppID: 1, WorkflowID: 1,
// 		Name: "trans-1", Trigger: wfmod.TransitionTriggerAuto,
// 	}})

// 	mocker.EXPECT().GetTransitionArcs(gomock.Any(), wfmod.ArcDirectionOut).Return([]wfmod.WfArc{{
// 		BaseModel: wfmod.BaseModel{ID: 2},
// 		AppID:     1, WorkflowID: 1, PlaceID: 2, TransitionID: 1, Direction: wfmod.ArcDirectionOut,
// 		Type: wfmod.ArcSEQ,
// 	}})

// 	mocker.EXPECT().GetPlacesByIDList(1, []int{2}).Return([]wfmod.WfPlace{{
// 		BaseModel: wfmod.BaseModel{ID: 2}, AppID: 1, WorkflowID: 1,
// 		Type: wfmod.MiddlePlaceType,
// 		Desc: "start place",
// 		Name: "middle-1",
// 	}})

// 	return mocker
// }

// func (s *seqMocker) close() {
// 	s.mockCtrl.Finish()
// }

// func TestStartSeq(t *testing.T) {
// 	workflow := flow.Workflow{
// 		WfWorkflow: wfmod.WfWorkflow{BaseModel: wfmod.BaseModel{ID: 1}, AppID: 1},
// 	}
// 	wfContext := wfmod.WfContextType{}

// 	seqMocker := newSeqMocker(t)
// 	flowDao := seqMocker.do(wfContext)
// 	defer seqMocker.close()

// 	flow.RegFlowDao(flowDao)

// 	workflow.Start(wfContext, "edenzou")
// }

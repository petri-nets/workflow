package flow

import (
	"github.com/petri-nets/workflow/wfmod"
)

// Transition Model
type Transition struct {
	WfTransition wfmod.WfTransition
	Case         *Case
	workitem     Workitem
}

// GetTransitionsByWorkitemJob get transitions by workitem job
func GetTransitionsByWorkitemJob(appID, workitemJobID int) []Transition {
	workitems := GetOpeningWorkitemsByJobID(appID, workitemJobID)

	var transitionIDList []int
	var caseIDList []int
	for _, workitem := range workitems {
		transitionIDList = append(transitionIDList, workitem.WfWorkitem.TransitionID)
		caseIDList = append(caseIDList, workitem.WfWorkitem.CaseID)
	}

	cases := GetCasesByIDList(appID, caseIDList)

	wfTransitions := GetWfTransitionsByIDList(appID, transitionIDList)

	transitions := []Transition{}
	for _, wfTransition := range wfTransitions {
		workitem := Workitem{}
		for _, item := range workitems {
			if item.WfWorkitem.TransitionID == wfTransition.ID {
				workitem = item
			}
		}

		cs := Case{}
		for _, item := range cases {
			if workitem.WfWorkitem.CaseID == item.WfCase.ID {
				cs = item
			}
		}

		transitions = append(transitions, Transition{
			WfTransition: wfTransition, workitem: workitem, Case: &cs})
	}

	return transitions
}

// GetWfTransitionsByIDList get wfTransitions by id list
func GetWfTransitionsByIDList(appID int, transitionIDList []int) []wfmod.WfTransition {
	return flowDao.GetTransitionsByIDList(appID, transitionIDList)
}

// Activate 激活transition
// 1. AND join: A transition with two or more input places and one output place.
// This will only be enabled once there is a token in all of the input places, which would be after each parallel thread of execution has finished.
// 2. Implicit OR split: An example of conditional routing where the decision is made as late as possible.
// Implicit or-splits are modeled as two or more arcs going from the same place but to different transitions.
// That way, the transition that happens to fire first will get the token. Once the token is gone, the remaining transitions are cancelled and thus cannot be fired.
// One of the transitions may have a timer as its trigger so that it will be fired if none of the other transitions is activated before the time limit expires.
// Expired transitions can either be triggered automatically via a background process which is running on a timer (e.g. cron), or manually via an online screen.
func (t *Transition) Activate() bool {
	tokens, complete := t.inPlaceFreeTokens()
	if !complete {
		return false
	}

	// 锁定token
	if err := t.lockTokens(tokens); err != nil {
		return false
	}

	// 合并token，并且合并上下文
	mergeToken := MergeTokens(tokens)

	// 创建新workitem
	t.workitem = CreateWorkitem(t, &mergeToken)

	// 创建workitem后，如果是自动触发类型，则直接触发workitem异步执行
	if t.WfTransition.Trigger == wfmod.TransitionTriggerAuto {
		if t.workitem.Do(t.Case) {
			t.WorkitemFinishedHandle(nil)
		}
	}

	return true
}

// WorkitemFinishedHandle 当workitem 执行完成后，触发执行
func (t *Transition) WorkitemFinishedHandle(output wfmod.WfContextType) {
	// 获取lock tokens
	places, tokens := t.inPlacesAndTokens(wfmod.TokenStatusLock)
	if len(places) != len(tokens) {
		// 异常
		return
	}

	// 扭转workitem状态至完成
	t.workitem.finish()

	// consume token
	if err := t.consumeTokens(tokens); err != nil {
		return
	}

	// Find the out arc
	// 1. AND split: An example of parallel routing where several tasks are performed in parallel or in no particular order.
	// It is modeled by a transition with one input place and two or more output places. When fired the transition will create tokens in all output places.

	// 2. Explicit OR split : An example of conditional routing where the decision is made as early as possible.
	// It is modeled by attaching conditions or guards to the arcs going out of a transition to different places.
	// Guard - An expression attached to an arc, shown in brackets, that evaluates to either TRUE or FALSE.
	// Tokens can only travel over arcs when their guard evaluates to TRUE. The expression will typically involve the case attributes.
	// More than two arcs of this type can come out of the same transition. They must all have a condition except the last arc as this will be
	// used as the default path if none of the other conditions evaluates to TRUE.

	mergeTokenContext := t.workitem.MergeOutputContext(output)

	outArcs := t.outArcs()

	// 获取畅通的路径列表
	for _, outArc := range outArcs {
		if outArc.CheckPass(mergeTokenContext) {
			place := GetPlaceByID(outArc.WfArc.AppID, outArc.WfArc.PlaceID)
			// 如果是end place则不生产新token
			if place.IsEndPlace() {
				t.Case.Close()
			} else {
				// 生产新token，并递交给place
				token := NewToken(t.Case, t.workitem.WfWorkitem.Context, outArc.WfArc.PlaceID)
				// 传递token
				go place.ReceiveToken(t.Case, &token)
			}

			if outArc.IsOutExplicitORSplit() {
				break
			}
		}
	}
}

func (t *Transition) outArcs() []Arc {
	return GetTransitionArcs(t, wfmod.ArcDirectionOut)
}

func (t *Transition) inPlacesAndTokens(tokenStatus wfmod.TokenStatusType) ([]Place, []Token) {
	// 通过arc找到inPlaces
	// var places []Arc
	inArcs := t.inArcs()

	var placeIDList []int
	for _, inArc := range inArcs {
		placeIDList = append(placeIDList, inArc.WfArc.PlaceID)
	}

	inPlaces := GetPlacesByIDList(t.WfTransition.AppID, placeIDList)
	tokens := GetPlacesTokens(t.Case, placeIDList, tokenStatus)
	return inPlaces, tokens
}

func (t *Transition) inArcs() []Arc {
	return GetTransitionArcs(t, wfmod.ArcDirectionIn)
}

// transition in
func (t *Transition) inPlaceFreeTokens() ([]Token, bool) {
	inPlaces, freeTokens := t.inPlacesAndTokens(wfmod.TokenStatusFree)
	complete := false
	if len(freeTokens) > 0 && len(inPlaces) == len(freeTokens) {
		complete = true
	}
	return freeTokens, complete

}

func (t *Transition) lockTokens(tokens []Token) error {
	for _, token := range tokens {
		if err := token.Lock(); err != nil {
			return err
		}
	}
	return nil
}

func (t *Transition) consumeTokens(tokens []Token) error {
	for _, token := range tokens {
		if err := token.Consume(); err != nil {
			return err
		}
	}
	return nil
}

func (t *Transition) mergeTokensContext(tokens []Token) error {
	for _, token := range tokens {
		if err := token.Consume(); err != nil {
			return err
		}
	}
	return nil
}

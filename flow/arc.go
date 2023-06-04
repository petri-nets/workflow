package flow

import (
	"sort"

	"github.com/Knetic/govaluate"
	"github.com/petri-nets/workflow/wfmod"
)

// Arc is model of wf_arc
type Arc struct {
	WfArc wfmod.WfArc
}

// Start arc启动
func (a *Arc) CheckPass(wfContext wfmod.WfContextType) bool {
	switch a.WfArc.Type {
	case wfmod.ArcSEQ, wfmod.ArcOutANDSplit, wfmod.ArcOutORJoin:
		return true
	case wfmod.ArcOutExplicitORsplit:
		if a.WfArc.Condition == "" {
			return true
		}
		return a.passCondition(wfContext)
	default:
		return false
	}
}

// IsOutExplicitORSplit check is out explicit or split type of arc
func (a *Arc) IsOutExplicitORSplit() bool {
	return a.WfArc.Type == wfmod.ArcOutExplicitORsplit
}

// passCondition can pass to transition
func (a *Arc) passCondition(wfContext wfmod.WfContextType) bool {
	expr, err := govaluate.NewEvaluableExpression(a.WfArc.Condition)
	if err != nil {
		// log.Fatal("syntax error:", err)
		return false
	}

	result, err := expr.Evaluate(wfContext)
	if err != nil {
		// log.Fatal("evaluate error:", err)
		return false
	}

	bol, ok := result.(bool)
	if !ok {
		return false
	}

	return bol
}

// GetTransitionArcs get transition arcs
func GetTransitionArcs(t *Transition, direct wfmod.ArcDirectionType) []Arc {
	wfArcs := flowDao.GetTransitionArcs(&t.WfTransition, direct)
	// 对arcs按id从小到大排序
	sort.SliceStable(wfArcs, func(i, j int) bool {
		if wfArcs[i].ID < wfArcs[j].ID {
			return true
		}
		return false
	})

	var arcs []Arc
	for _, wfArc := range wfArcs {
		arcs = append(arcs, Arc{WfArc: wfArc})
	}
	return arcs
}

// GetPlaceArcs get place arcs
func GetPlaceArcs(p *Place, direct wfmod.ArcDirectionType) []Arc {
	wfArcs := flowDao.GetPlaceArcs(&p.WfPlace, direct)

	var arcs []Arc
	for _, wfArc := range wfArcs {
		arcs = append(arcs, Arc{WfArc: wfArc})
	}
	return arcs
}

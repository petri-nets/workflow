package flow

import (
	"fmt"
	"time"

	"github.com/petri-nets/workflow/wfmod"
)

// RouterType router type
type RouterType string

const (
	RouterEIF RouterType = "eif" // 选择路由
	RouterVIE RouterType = "vie" // 争夺路由
	RouterSEQ RouterType = "seq" // 顺序路由
	RouterPAR RouterType = "par" // 并行路由
)

const (
	StartPlace string = "start"
	EndPlace   string = "end"
)

type task struct {
	Name      string                      `yaml:"name" json:"name"`
	Desc      string                      `yaml:"desc" json:"desc"`
	Trigger   wfmod.TransitionTriggerType `yaml:"trigger" json:"trigger"`
	Job       int                         `yaml:"job_id" json:"job_id"`
	Role      int                         `yaml:"role_id" json:"role_id"`
	Router    RouterType                  `yaml:"router" json:"router"`
	Condition string                      `yaml:"condition" json:"condition"`
	Tasks     []*task                     `yaml:"tasks" json:"tasks"`
	Then      *task                       `yaml:"then" json:"then"`
	Goto      string                      `yaml:"goto" json:"goto"` // 仅在if路由下有效
}

// Pipeline pipeline
type Pipeline struct {
	ID         int        `yaml:"id" json:"id"`
	Name       string     `yaml:"name" json:"name"`
	Desc       string     `yaml:"desc" json:"desc"`
	AppID      int        `yaml:"app_id" json:"app_id"`
	StartJobID int        `yaml:"start_job_id" json:"start_job_id"`
	StartDate  time.Time  `yaml:"start_date" json:"start_date"`
	EndDate    time.Time  `yaml:"end_date" json:"end_date"`
	Operator   string     `yaml:"operator" json:"operator"`
	Router     RouterType `yaml:"router" json:"router"`
	Tasks      []*task    `yaml:"tasks" json:"tasks"`
	Then       *task      `yaml:"then" json:"then"`
}

// Parser parser
type Parser struct {
	Pipeline    *Pipeline
	transitions []transition
	places      []place
	index       int
}

// NewParser new parser
func NewParser(pipe *Pipeline) *Parser {
	parser := Parser{
		Pipeline: pipe,
	}
	return &parser
}

// Start start
func (p *Parser) Start() error {
	if err := p.Validate(); err != nil {
		return err
	}
	// 编排
	p.Arrangement()
	// 保存模型
	return p.SaveModel()
}

// SaveModel save model
func (p *Parser) SaveModel() error {
	switcher := Switcher{
		Workflow: workflow{
			ID:         p.Pipeline.ID,
			AppID:      p.Pipeline.AppID,
			Name:       p.Pipeline.Name,
			Desc:       p.Pipeline.Desc,
			StartJobID: p.Pipeline.StartJobID,
			StartAt:    p.Pipeline.StartDate,
			EndAt:      p.Pipeline.EndDate,
		},
		Places:      p.places,
		Transitions: p.transitions,
	}

	return switcher.Save()
}

// Arrangement arrangement
func (p *Parser) Arrangement() error {
	startPlace := p.newPlace(wfmod.StartPlaceType, "")
	endPlace := p.newPlace(wfmod.EndPlaceType, "")

	startTask := task{
		Name:    p.Pipeline.Name,
		Desc:    p.Pipeline.Desc,
		Trigger: wfmod.TransitionTriggerAuto,
		Router:  p.Pipeline.Router,
		Then:    p.Pipeline.Then,
	}

	// 中间任务编排
	inPlaces, outPlaces := p.ArrangeTasks(&startTask, p.Pipeline.Tasks)
	if len(inPlaces) == 0 {
		return fmt.Errorf("arrange task error: inPlaces empty")
	}

	if len(outPlaces) == 0 {
		return fmt.Errorf("arrange task error: outPlaces empty")
	}

	// startPlace && startTransition
	p.newTransition(&startTask, wfmod.ArcSEQ, p.arcOutTypeByRouter(startTask.Router), []place{startPlace}, inPlaces)

	// endPlaces && endTransition
	p.replaceAndDeletePlaces(outPlaces, endPlace)

	return nil
}

func (p *Parser) arcOutTypeByRouter(router RouterType) wfmod.ArcTypeType {
	switch router {
	case RouterEIF:
		return wfmod.ArcOutExplicitORsplit
	case RouterPAR:
		return wfmod.ArcOutANDSplit
	case RouterSEQ, RouterVIE:
		return wfmod.ArcSEQ
	}
	return wfmod.ArcSEQ
}

func (p *Parser) thenInArcTypeByRouter(router RouterType) wfmod.ArcTypeType {
	switch router {
	case RouterSEQ, RouterEIF:
		return wfmod.ArcSEQ
	case RouterPAR:
		return wfmod.ArcInANDJoin
	case RouterVIE:
		return wfmod.ArcInImplicitORsplit
	}
	return wfmod.ArcSEQ
}

func (p *Parser) Validate() error {
	return nil
}

// ArrangeTasks save tasks
func (p *Parser) ArrangeTasks(t *task, tasks []*task) (inPlaces, outPlaces []place) {
	switch t.Router {
	case RouterEIF:
		inPlaces, outPlaces = p.EIfRouter(tasks)
	case RouterSEQ:
		inPlaces, outPlaces = p.SeqRouter(tasks)
	case RouterPAR:
		inPlaces, outPlaces = p.ParRouter(tasks)
		// 并行路由必须存在一个回收transition
		if t.Then == nil {
			t.Then = &task{
				Name:    "t-auto",
				Trigger: wfmod.TransitionTriggerAuto,
			}
		}
	case RouterVIE:
		inPlaces, outPlaces = p.VieRouter(tasks)
	}

	if t.Then != nil {
		inArcType := p.thenInArcTypeByRouter(t.Router)
		outPlace := p.newPlace(wfmod.MiddlePlaceType, "")
		p.newTransition(t.Then, inArcType, wfmod.ArcSEQ, outPlaces, []place{outPlace})
		outPlaces = []place{outPlace}
	}

	return inPlaces, outPlaces
}

func (p *Parser) newPlace(placeType wfmod.PlaceTypeType, cond string) place {
	name := fmt.Sprintf("p-%d", p.index)
	if placeType == wfmod.StartPlaceType {
		name = StartPlace
	} else if placeType == wfmod.EndPlaceType {
		name = EndPlace
	}
	place := place{
		Name:      name,
		Type:      placeType,
		condition: cond,
	}
	p.places = append(p.places, place)
	p.index++
	return place
}

func (p *Parser) newTransition(
	t *task, InRouter, outRouter wfmod.ArcTypeType, inPlaces, outPlaces []place) *transition {

	inArcs := p.buildTransitionInArcs(InRouter, t, inPlaces)
	outArcs := p.buildTransitionOutArcs(outRouter, t, outPlaces)
	tran := transition{
		Name:    t.Name,
		Desc:    t.Desc,
		Trigger: t.Trigger,
		JobID:   t.Job,
		RoleID:  t.Role,
		In:      inArcs,
		Out:     outArcs,
	}
	p.transitions = append(p.transitions, tran)
	return &tran
}

func (p *Parser) buildTransitionInArcs(router wfmod.ArcTypeType, t *task, inPlaces []place) []arc {
	// transition in arc
	// 当只有一个inPlace时，有两种类型的arc，分别是seq，和Implicit split
	var inArcs []arc
	for _, inPlace := range inPlaces {
		inArcs = append(inArcs, arc{Place: inPlace.Name, Type: router})
	}
	return inArcs
}

func (p *Parser) buildTransitionOutArcs(router wfmod.ArcTypeType, t *task, outPlaces []place) []arc {
	var inArcs []arc
	for _, outPlace := range outPlaces {
		arc := arc{Place: outPlace.Name}
		if router == wfmod.ArcOutExplicitORsplit {
			arc.Condition = outPlace.condition
			arc.Type = wfmod.ArcOutExplicitORsplit
		} else {
			arc.Type = router
		}
		inArcs = append(inArcs, arc)
	}
	return inArcs
}

// EIfRouter else if router
func (p *Parser) EIfRouter(tasks []*task) (inPlaces, outPlaces []place) {
	// ifRouter 前一个transition 通过ifRouter 拆分成多条子路，每条子路一个place
	outPlace := p.newPlace(wfmod.MiddlePlaceType, "")
	outPlaces = append(outPlaces, outPlace)
	// 获取无条件的task作为默认task，当且仅当其他没有任何满足条件的task情况下，才生效默认task，如果未设置默认路由，则自动生成一条默认路由
	tasks = p.setEIFRouterDefaultTask(tasks)

	for _, task := range tasks {

		inPlace := p.newPlace(wfmod.MiddlePlaceType, task.Condition)
		inPlaces = append(inPlaces, inPlace)

		if len(task.Tasks) > 0 {
			nextInPlaces, nextOutPlaces := p.ArrangeTasks(task, task.Tasks)
			p.newTransition(task, wfmod.ArcSEQ, wfmod.ArcOutORJoin, []place{inPlace}, nextInPlaces)
			// 需要将nextOutPlaces 删除，并将所有引用的地方修改成outPlace
			p.replaceAndDeletePlaces(nextOutPlaces, outPlace)
		} else {
			p.newTransition(task, wfmod.ArcSEQ, wfmod.ArcOutORJoin, []place{inPlace}, outPlaces)
		}
	}

	return inPlaces, outPlaces
}

func (p *Parser) setEIFRouterDefaultTask(tasks []*task) []*task {
	newTasks := []*task{}
	var defaultTask *task
	for _, task := range tasks {
		if task.Condition != "" {
			newTasks = append(newTasks, task)
		} else {
			defaultTask = task
		}
	}
	// 当未设置默认task时，自动生成一条无条件的task
	if defaultTask == nil {
		defaultTask = &task{
			Name: "t-default",
		}
	}
	newTasks = append(newTasks, defaultTask)
	return newTasks
}

// VieRouter vie router
func (p *Parser) VieRouter(tasks []*task) (inPlaces, outPlaces []place) {
	outPlace := p.newPlace(wfmod.MiddlePlaceType, "")
	outPlaces = append(outPlaces, outPlace)
	inPlace := p.newPlace(wfmod.MiddlePlaceType, "")
	inPlaces = append(inPlaces, inPlace)

	for _, task := range tasks {
		if len(task.Tasks) > 0 {
			nextInPlaces, nextOutPlaces := p.ArrangeTasks(task, task.Tasks)
			p.newTransition(task, wfmod.ArcInImplicitORsplit, wfmod.ArcOutORJoin, inPlaces, nextInPlaces)
			// 需要将nextOutPlaces 删除，并将所有引用的地方修改成outPlace
			p.replaceAndDeletePlaces(nextOutPlaces, outPlace)
		} else {
			p.newTransition(task, wfmod.ArcInImplicitORsplit, wfmod.ArcOutORJoin, inPlaces, outPlaces)
		}
	}

	return inPlaces, outPlaces
}

// SeqRouter sequel router
func (p *Parser) SeqRouter(tasks []*task) ([]place, []place) {
	inPlace := p.newPlace(wfmod.MiddlePlaceType, "")

	inPlaces := []place{inPlace}
	inArcType := wfmod.ArcSEQ

	for _, task := range tasks {
		if len(task.Tasks) > 0 {
			nextInPlaces, nextOutPlaces := p.ArrangeTasks(task, task.Tasks)
			p.newTransition(task, inArcType, p.arcOutTypeByRouter(task.Router), inPlaces, nextInPlaces)

			inPlaces = nextOutPlaces
		} else {
			outPlace := p.newPlace(wfmod.MiddlePlaceType, "")
			p.newTransition(task, inArcType, p.arcOutTypeByRouter(task.Router), inPlaces, []place{outPlace})
			inPlaces = []place{outPlace}
		}
		inArcType = p.thenInArcTypeByRouter(task.Router)
	}

	return []place{inPlace}, inPlaces
}

// ParRouter parallel router
func (p *Parser) ParRouter(tasks []*task) (inPlaces, outPlaces []place) {
	for _, task := range tasks {
		inPlace := p.newPlace(wfmod.MiddlePlaceType, "")
		inPlaces = append(inPlaces, inPlace)

		var nextInPlaces []place
		var nextOutPlaces []place
		if len(task.Tasks) > 0 {
			// outPlaces
			nextInPlaces, nextOutPlaces = p.ArrangeTasks(task, task.Tasks)
			p.newTransition(task, wfmod.ArcSEQ, p.arcOutTypeByRouter(task.Router), []place{inPlace}, nextInPlaces)
			outPlaces = append(outPlaces, nextOutPlaces...)
		} else {
			outPlace := p.newPlace(wfmod.MiddlePlaceType, "")
			p.newTransition(task, wfmod.ArcSEQ, wfmod.ArcSEQ, []place{inPlace}, []place{outPlace})
			outPlaces = append(outPlaces, outPlace)
		}
	}
	return inPlaces, outPlaces
}

func (p *Parser) replaceAndDeletePlaces(fromPlaces []place, toPlace place) {
	newPlaces := []place{}
	for _, place := range p.places {
		bol := false
		for _, fromPlace := range fromPlaces {
			if place.Name == fromPlace.Name {
				bol = true
				break
			}
		}

		if bol {
			// 替换所有
			// 替换arcs
			for i, transition := range p.transitions {
				inArcs := []arc{}
				for _, a := range transition.In {
					if a.Place == place.Name {
						a.Place = toPlace.Name
					}
					inArcs = append(inArcs, a)
				}
				transition.In = inArcs

				outArcs := []arc{}
				for _, a := range transition.Out {
					if a.Place == place.Name {
						a.Place = toPlace.Name
					}
					outArcs = append(outArcs, a)
				}
				transition.Out = outArcs

				p.transitions[i] = transition
			}

		} else {
			newPlaces = append(newPlaces, place)
		}
	}

	p.places = newPlaces

}

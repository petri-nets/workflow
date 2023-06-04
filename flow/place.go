package flow

import (
	"github.com/petri-nets/workflow/wfmod"
)

// Place Model
type Place struct {
	WfPlace wfmod.WfPlace
}

// GetStartPlace get workflow start place
func GetStartPlace(cs *Case) Place {
	wfPlaces := getWfPlaces(cs, wfmod.StartPlaceType)
	return Place{WfPlace: wfPlaces[0]}
}

func getWfPlaces(cs *Case, placeType wfmod.PlaceTypeType) []wfmod.WfPlace {
	return flowDao.GetPlacesByType(cs.WfCase.AppID, cs.WfCase.WorkflowID, placeType)
}

// GetPlaceByID get place by id
func GetPlaceByID(appID, placeID int) Place {
	places := GetPlacesByIDList(appID, []int{placeID})
	if len(places) > 0 {
		return places[0]
	}
	return Place{}
}

// ReceiveToken receive token start
func (p *Place) ReceiveToken(cs *Case, token *Token) {
	if !p.IsEndPlace() {
		p.activateTransitions(cs, token)
	}
}

func (p *Place) activateTransitions(cs *Case, token *Token) {
	// 激活transition
	transitions := p.outTransitions(cs)
	for _, transition := range transitions {
		transition.Activate()
	}
}

// IsEndPlace check is end place
func (p *Place) IsEndPlace() bool {
	return p.WfPlace.Type == wfmod.EndPlaceType
}

func (p *Place) outTransitions(cs *Case) []Transition {
	outArcs := GetPlaceArcs(p, wfmod.ArcDirectionIn)
	transitionIDList := []int{}
	for _, outArc := range outArcs {
		transitionIDList = append(transitionIDList, outArc.WfArc.TransitionID)
	}

	transitions := []Transition{}
	wfTransitions := GetWfTransitionsByIDList(p.WfPlace.AppID, transitionIDList)
	for _, wfTransition := range wfTransitions {
		transitions = append(transitions, Transition{
			WfTransition: wfTransition,
			Case:         cs,
		})
	}
	return transitions
}

// GetPlacesByIDList get places by id list
func GetPlacesByIDList(appID int, idList []int) []Place {
	wfPlaces := flowDao.GetPlacesByIDList(appID, idList)

	var places []Place
	for _, wfPlace := range wfPlaces {
		places = append(places, Place{WfPlace: wfPlace})
	}

	return places
}

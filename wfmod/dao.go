package wfmod

// FlowDao flow data interface
// mockgen -source=./wfmod/dao.go -destination=./wfmod/wfmock/mock.go -package=wfmock
type FlowDao interface {
	// GetWorkflowsByStartJob get workflow by start job id
	GetWorkflowsByStartJob(appID, startJobID int) []WfWorkflow
	// SaveWorkflow create workflow or update workflow when give id
	SaveWorkflow(w *WfWorkflow) error

	// GetPlaceArcs get place arcs
	GetPlaceArcs(p *WfPlace, direct ArcDirectionType) []WfArc
	// GetTransitionArcs get transition arcs
	GetTransitionArcs(t *WfTransition, direct ArcDirectionType) []WfArc
	// SaveArc save arc or update when give id
	SaveArc(a *WfArc) error

	// GetCasesByIDList get cases by id list
	GetCasesByIDList(appID int, idList []int) []WfCase
	// SaveCase save case or update when give id
	SaveCase(c *WfCase) error

	// GetPlacesByType get places by place type
	GetPlacesByType(appID int, workflowID int, placeType PlaceTypeType) []WfPlace
	// GetPlacesByIDList get places by id list
	GetPlacesByIDList(appID int, idList []int) []WfPlace
	// SavePlace save place or update when give id
	SavePlace(p *WfPlace) error

	// GetTransitionsByIDList get transitions by id list
	GetTransitionsByIDList(appID int, idList []int) []WfTransition
	// SaveTransition save transition or update transition when give id
	SaveTransition(p *WfTransition) error

	// GetTokenByIDList get tokens by id list
	GetTokenByIDList(appID int, idList []int) []WfToken
	// GetTokensByPlaces get token by place id list
	GetTokensByPlaces(cs *WfCase, placeIDList []int, status TokenStatusType) []WfToken
	// SaveToken save token or update token when give id
	SaveToken(t *WfToken) error

	// GetOpeningWorkitemsByJobID get opening workitems by job id
	GetOpeningWorkitemsByJobID(appID, jobID int) []WfWorkitem
	// GetWorkitemsByIDList get workitems by id list
	GetWorkitemsByIDList(appID int, idList []int) []WfWorkitem
	// SaveWorkitem save workitem or update workitem when give id
	SaveWorkitem(w *WfWorkitem) error

	// BeginTransaction begin a transaction
	BeginTransaction()
	// CommitTransaction commit a transaction
	CommitTransaction()
	// RollbackTransaction commit a transaction
	RollbackTransaction()
}

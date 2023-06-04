package wfmod

// ArcDirectionType
type ArcDirectionType string

const (
	ArcDirectionIn  ArcDirectionType = "IN"  //入弧
	ArcDirectionOut ArcDirectionType = "OUT" //出弧
)

// ArcTypeType
type ArcTypeType string

const (
	ArcSEQ                ArcTypeType = "SEQ"               //连接弧
	ArcOutExplicitORsplit ArcTypeType = "Explicit OR split" //显式或拆分
	ArcInImplicitORsplit  ArcTypeType = "Implicit OR split" //隐式或拆分
	ArcOutORJoin          ArcTypeType = "OR join"           //或合并，跟或拆分联合使用
	ArcOutANDSplit        ArcTypeType = "AND split"         //与拆分
	ArcInANDJoin          ArcTypeType = "AND join"          //与合并
)

// OutArcTypes a map contain all valid OutArcTypes
var OutArcTypes = map[ArcTypeType]ArcTypeType{
	ArcSEQ:                ArcSEQ,
	ArcOutORJoin:          ArcOutORJoin,
	ArcOutExplicitORsplit: ArcOutExplicitORsplit,
	ArcOutANDSplit:        ArcOutANDSplit}

// InArcTypes a map contain all valid InArcTypes
var InArcTypes = map[ArcTypeType]ArcTypeType{
	ArcSEQ:               ArcSEQ,
	ArcInImplicitORsplit: ArcInImplicitORsplit,
	ArcInANDJoin:         ArcInANDJoin}

// WfArc is model of wf_arc
type WfArc struct {
	BaseModel
	AppID        int              `json:"app_id"`
	WorkflowID   int              `json:"workflow_id"`
	TransitionID int              `json:"transition_id"`
	PlaceID      int              `json:"place_id"`
	Direction    ArcDirectionType `json:"direction"`
	Type         ArcTypeType      `json:"type"`
	Condition    string           `json:"condition"`
}

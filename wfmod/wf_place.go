package wfmod

type PlaceTypeType int

// PlaceTypeType
const (
	StartPlaceType  PlaceTypeType = 1 // 开始类型
	MiddlePlaceType PlaceTypeType = 5 // 中间类型
	EndPlaceType    PlaceTypeType = 9 // 结束类型
)

// WfPlace Model
type WfPlace struct {
	BaseModel
	AppID      int           `json:"app_id"`
	WorkflowID int           `json:"workflow_id"`
	Type       PlaceTypeType `json:"type"`
	Desc       string        `json:"desc"`
	Name       string        `json:"name"`
}

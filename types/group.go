package types

// GroupListResponse 表示群列表接口响应参数
type GroupListResponse struct {
	Items []*GroupItem `json:"items"` // 群列表
}

// GroupItem 表示单个群的信息
type GroupItem struct {
	GroupID   int32  `json:"group_id"`   // 群 ID
	GroupName string `json:"group_name"` // 群名称
	Avatar    string `json:"avatar"`     // 群头像
	Profile   string `json:"profile"`    // 群简介
	Leader    int32  `json:"leader"`     // 群主 ID
	CreatorID int32  `json:"creator_id"` // 群创建者 ID
}

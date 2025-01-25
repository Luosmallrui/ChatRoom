package types

// ContactListResponse 表示联系人列表响应
type ContactListResponse struct {
	Items []*ContactItem `json:"items"` // 联系人列表
}

// ContactItem 表示联系人信息
type ContactItem struct {
	UserID   int32  `json:"user_id"`  // 用户 ID
	Nickname string `json:"nickname"` // 昵称
	Gender   int32  `json:"gender"`   // 性别 [0:未知;1:男;2:女]
	Motto    string `json:"motto"`    // 座右铭
	Avatar   string `json:"avatar"`   // 头像
	Remark   string `json:"remark"`   // 备注
	GroupID  int32  `json:"group_id"` // 联系人分组 ID
}

// ContactGroupListResponse 表示联系人分组列表响应参数
type ContactGroupListResponse struct {
	Items []*ContactGroupItem `json:"items"` // 分组列表
}

// ContactGroupItem 表示单个联系人分组信息
type ContactGroupItem struct {
	ID    int32  `json:"id"`    // 分组 ID
	Name  string `json:"name"`  // 分组名称
	Count int32  `json:"count"` // 联系人数
	Sort  int32  `json:"sort"`  // 分组排序
}

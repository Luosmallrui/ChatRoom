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

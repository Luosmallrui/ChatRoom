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

type GetInviteFriendsRequest struct {
	GroupID int32 `form:"group_id"`
}

// GetInviteFriendsResponse 获取邀请好友响应
type GetInviteFriendsResponse struct {
	Items []*GetInviteFriendsResponseItem `json:"items" form:"items"`
}

// GetInviteFriendsResponseItem 邀请好友项
type GetInviteFriendsResponseItem struct {
	UserID   int32  `json:"user_id" form:"user_id"`   // 用户ID
	Nickname string `json:"nickname" form:"nickname"` // 用户昵称
	Avatar   string `json:"avatar" form:"avatar"`     // 头像地址
	Gender   int32  `json:"gender" form:"gender"`     // 性别(1:男 2:女)
	Remark   string `json:"remark" form:"remark"`     // 备注
}

// GroupCreateRequest 创建群聊请求参数
type GroupCreateRequest struct {
	Name    string  `json:"name" form:"name" binding:"required"`         // 群名称
	UserIds []int32 `json:"user_ids" form:"user_ids" binding:"required"` // 用户ID列表
}

// GroupCreateResponse 创建群聊响应参数
type GroupCreateResponse struct {
	GroupID int32 `json:"group_id"` // 群组ID
}

// GroupMemberListRequest 群成员列表请求参数
type GroupMemberListRequest struct {
	GroupID int32 `json:"group_id" form:"group_id" binding:"required"` // 群组ID
}

// GroupMemberListResponse 群成员列表响应参数
type GroupMemberListResponse struct {
	Items []*GroupMemberListResponseItem `json:"items"` // 群成员列表
}

// GroupMemberListResponseItem 群成员信息项
type GroupMemberListResponseItem struct {
	UserID   int32  `json:"user_id"`  // 用户ID
	Nickname string `json:"nickname"` // 用户昵称
	Avatar   string `json:"avatar"`   // 头像地址
	Gender   int32  `json:"gender"`   // 性别
	Leader   int32  `json:"leader"`   // 是否群主/管理员
	IsMute   int32  `json:"is_mute"`  // 是否被禁言
	Remark   string `json:"remark"`   // 备注
	Motto    string `json:"motto"`    // 个性签名
}

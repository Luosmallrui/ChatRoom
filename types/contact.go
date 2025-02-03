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

type ContactSearchRequest struct {
	Mobile string `form:"mobile" binding:"required"`
}

type ContactSearchResponse struct {
	UserID   int32  `json:"user_id"`  // 用户ID
	Mobile   string `json:"mobile"`   // 手机号
	Nickname string `json:"nickname"` // 昵称
	Avatar   string `json:"avatar"`   // 头像URL
	Gender   int32  `json:"gender"`   // 性别
	Motto    string `json:"motto"`    // 个性签名
}

type ContactDetailRequest struct {
	UserID int32 `form:"user_id" binding:"required"` // 用户ID，必填
}

type ContactDetailResponse struct {
	UserID     int32       `json:"user_id"`     // 用户ID
	Mobile     string      `json:"mobile"`      // 手机号
	Nickname   string      `json:"nickname"`    // 昵称
	Avatar     string      `json:"avatar"`      // 头像URL
	Gender     int32       `json:"gender"`      // 性别
	Motto      string      `json:"motto"`       // 个性签名
	Email      string      `json:"email"`       // 邮箱
	FriendInfo *FriendInfo `json:"friend_info"` // 好友信息
}

type FriendInfo struct {
	IsFriend string `json:"is_friend"` // 是否是好友
	GroupID  int32  `json:"group_id"`  // 分组ID
	Remark   string `json:"remark"`    // 好友备注
}

type ContactApplyCreateRequest struct {
	UserID int32  `json:"user_id" binding:"required"` // 用户ID，必填
	Remark string `json:"remark" binding:"required"`  // 添加好友备注，必填
}

type ContactApplyCreateResponse struct{}

type ContactApplyAcceptRequest struct {
	ApplyID int32  `json:"apply_id" binding:"required"` // 申请ID，必填
	Remark  string `json:"remark"`
}

type ContactApplyAcceptResponse struct{}

type ContactApplyListResponse struct {
	Items []*ContactApplyListItem `json:"items"` // 联系人申请列表
}

type ContactApplyListItem struct {
	ID        int32  `json:"id"`         // 申请记录ID
	UserID    int32  `json:"user_id"`    // 当前用户ID
	FriendID  int32  `json:"friend_id"`  // 好友ID
	Remark    string `json:"remark"`     // 申请备注
	Nickname  string `json:"nickname"`   // 好友昵称
	Avatar    string `json:"avatar"`     // 好友头像URL
	CreatedAt string `json:"created_at"` // 申请创建时间
}

type ContactOnlineStatusRequest struct {
	UserID int32 `json:"user_id"` // 用户ID，必填
}

type ContactOnlineStatusResponse struct {
	OnlineStatus int32 `json:"online_status"` // 在线状态 [1:离线; 2:在线]
}

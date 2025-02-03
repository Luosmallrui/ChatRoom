package types

// 聊天模式
const (
	ChatPrivateMode = 1 // 私信模式
	ChatGroupMode   = 2 // 群聊模式
)

const (
	PushEventImMessage         = "im.message"          // 对话消息推送
	PushEventImMessageKeyboard = "im.message.keyboard" // 键盘输入事件推送
	PushEventImMessageRevoke   = "im.message.revoke"   // 聊天消息撤销推送
	PushEventContactApply      = "im.contact.apply"    // 好友申请消息推送
	PushEventContactStatus     = "im.contact.status"   // 用户在线状态推送
	PushEventGroupApply        = "im.group.apply"      // 用户在线状态推送
)

// IM消息类型
// 1-999    自定义消息类型
// 1000-1999 系统消息类型
const (
	ChatMsgTypeText        = 1  // 文本消息
	ChatMsgTypeCode        = 2  // 代码消息
	ChatMsgTypeImage       = 3  // 图片文件
	ChatMsgTypeAudio       = 4  // 语音文件
	ChatMsgTypeVideo       = 5  // 视频文件
	ChatMsgTypeFile        = 6  // 其它文件
	ChatMsgTypeLocation    = 7  // 位置消息
	ChatMsgTypeCard        = 8  // 名片消息
	ChatMsgTypeForward     = 9  // 转发消息
	ChatMsgTypeLogin       = 10 // 登录消息
	ChatMsgTypeVote        = 11 // 投票消息
	ChatMsgTypeMixed       = 12 // 图文消息
	ChatMsgTypeGroupNotice = 13 // 群公告消息

	ChatMsgSysText                   = 1000 // 系统文本消息
	ChatMsgSysGroupCreate            = 1101 // 创建群聊消息
	ChatMsgSysGroupMemberJoin        = 1102 // 加入群聊消息
	ChatMsgSysGroupMemberQuit        = 1103 // 群成员退出群消息
	ChatMsgSysGroupMemberKicked      = 1104 // 踢出群成员消息
	ChatMsgSysGroupMessageRevoke     = 1105 // 管理员撤回成员消息
	ChatMsgSysGroupDismissed         = 1106 // 群解散
	ChatMsgSysGroupMuted             = 1107 // 群禁言
	ChatMsgSysGroupCancelMuted       = 1108 // 群解除禁言
	ChatMsgSysGroupMemberMuted       = 1109 // 群成员禁言
	ChatMsgSysGroupMemberCancelMuted = 1110 // 群成员解除禁言
	ChatMsgSysGroupNotice            = 1111 // 编辑群公告
	ChatMsgSysGroupTransfer          = 1113 // 变更群主
)

var ChatMsgTypeMapping = map[int]string{
	ChatMsgTypeImage:                 "[图片消息]",
	ChatMsgTypeAudio:                 "[语音消息]",
	ChatMsgTypeVideo:                 "[视频消息]",
	ChatMsgTypeFile:                  "[文件消息]",
	ChatMsgTypeLocation:              "[位置消息]",
	ChatMsgTypeCard:                  "[名片消息]",
	ChatMsgTypeForward:               "[转发消息]",
	ChatMsgTypeLogin:                 "[登录消息]",
	ChatMsgTypeVote:                  "[投票消息]",
	ChatMsgTypeCode:                  "[代码消息]",
	ChatMsgTypeMixed:                 "[图文消息]",
	ChatMsgSysText:                   "[系统消息]",
	ChatMsgSysGroupCreate:            "[创建群消息]",
	ChatMsgSysGroupMemberJoin:        "[加入群消息]",
	ChatMsgSysGroupMemberQuit:        "[退出群消息]",
	ChatMsgSysGroupMemberKicked:      "[踢出群消息]",
	ChatMsgSysGroupMessageRevoke:     "[撤回消息]",
	ChatMsgSysGroupDismissed:         "[群解散消息]",
	ChatMsgSysGroupMuted:             "[群禁言消息]",
	ChatMsgSysGroupCancelMuted:       "[群解除禁言消息]",
	ChatMsgSysGroupMemberMuted:       "[群成员禁言消息]",
	ChatMsgSysGroupMemberCancelMuted: "[群成员解除禁言消息]",
}

type TalkLastMessage struct {
	MsgId      string // 消息ID
	Sequence   int    // 消息时序ID（消息排序）
	MsgType    int    // 消息类型
	UserId     int    // 发送者ID
	ReceiverId int    // 接受者ID
	Content    string // 消息内容
	Mention    []int  // 提及列表
	CreatedAt  string // 消息发送时间
}

type TalkRecord struct {
	MsgId      string `json:"msg_id"`
	Sequence   int    `json:"sequence"`
	MsgType    int    `json:"msg_type"`
	FromUserId int    `json:"from_user_id"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	Extra      any    `json:"extra"`
	CreatedAt  string `json:"created_at"`
}

type TalkSessionItem struct {
	// ID represents the unique identifier of the talk session
	ID int32 `json:"id"`

	// TalkMode indicates the chat mode: 1 for private chat, 2 for group chat
	TalkMode int32 `json:"talk_mode"`

	// ToFromID represents the ID of the other participant
	ToFromID int32 `json:"to_from_id"`

	// IsTop indicates if the chat is pinned to the top (1 for true, 0 for false)
	IsTop int32 `json:"is_top"`

	// IsDisturb indicates if "do not disturb" is enabled for the chat (1 for true, 0 for false)
	IsDisturb int32 `json:"is_disturb"`

	// IsOnline indicates if the other participant is online (1 for true, 0 for false)
	IsOnline int32 `json:"is_online"`

	// IsRobot indicates if the other participant is a robot (1 for true, 0 for false)
	IsRobot int32 `json:"is_robot"`

	// Name is the display name of the other participant or group
	Name string `json:"name"`

	// Avatar is the URL of the avatar for the chat
	Avatar string `json:"avatar"`

	// Remark is an optional note or alias for the chat
	Remark string `json:"remark"`

	// UnreadNum represents the number of unread messages in the chat
	UnreadNum int32 `json:"unread_num"`

	// MsgText contains the most recent message text
	MsgText string `json:"msg_text"`

	// UpdatedAt is the timestamp of the last update to the chat session
	UpdatedAt string `json:"updated_at"`
}

// TalkSessionCreateResponse represents the response parameters for the session creation API.
type TalkSessionCreateResponse struct {
	// ID represents the unique identifier of the talk session
	ID int32 `json:"id"`

	// TalkMode indicates the chat mode: 1 for private chat, 2 for group chat
	TalkMode int32 `json:"talk_mode"`

	// ToFromID represents the ID of the other participant
	ToFromID int32 `json:"to_from_id"`

	// IsTop indicates if the chat is pinned to the top (1 for true, 0 for false)
	IsTop int32 `json:"is_top"`

	// IsDisturb indicates if "do not disturb" is enabled for the chat (1 for true, 0 for false)
	IsDisturb int32 `json:"is_disturb"`

	// IsOnline indicates if the other participant is online (1 for true, 0 for false)
	IsOnline int32 `json:"is_online"`

	// IsRobot indicates if the other participant is a robot (1 for true, 0 for false)
	IsRobot int32 `json:"is_robot"`

	// Name is the display name of the other participant or group
	Name string `json:"name"`

	// Avatar is the URL of the avatar for the chat
	Avatar string `json:"avatar"`

	// Remark is an optional note or alias for the chat
	Remark string `json:"remark"`

	// UnreadNum represents the number of unread messages in the chat
	UnreadNum int32 `json:"unread_num"`

	// MsgText contains the most recent message text
	MsgText string `json:"msg_text"`

	// UpdatedAt is the timestamp of the last update to the chat session
	UpdatedAt string `json:"updated_at"`
}

// TalkSessionDeleteRequest represents the request parameters for deleting a talk session.
type TalkSessionDeleteRequest struct {
	// TalkMode indicates the chat mode (e.g., 1 for private chat, 2 for group chat)
	TalkMode int32 `json:"talk_mode" binding:"required"`

	// ToFromID represents the ID of the other participant in the chat
	ToFromID int32 `json:"to_from_id" binding:"required"`
}

type TalkSessionDeleteResponse struct{}

// TalkSessionTopRequest represents the request parameters for toggling the "top" (pinned) status of a talk session.
type TalkSessionTopRequest struct {
	// TalkMode indicates the chat mode (e.g., 1 for private chat, 2 for group chat)
	TalkMode int32 `json:"talk_mode" binding:"required"`

	// ToFromID represents the ID of the other participant in the chat
	ToFromID int32 `json:"to_from_id" binding:"required"`

	// Action specifies the action to perform: 1 to pin the session to the top, 2 to unpin
	Action int32 `json:"action" binding:"required,oneof=1 2"`
}

type TalkSessionTopResponse struct {
}

// TalkSessionDisturbRequest represents the request parameters for toggling the "Do Not Disturb" status of a talk session.
type TalkSessionDisturbRequest struct {
	// TalkMode indicates the chat mode (e.g., 1 for private chat, 2 for group chat)
	TalkMode int32 `json:"talk_mode" binding:"required"`

	// ToFromID represents the ID of the other participant in the chat
	ToFromID int32 `json:"to_from_id" binding:"required"`

	// Action specifies the action to perform: 1 to enable "Do Not Disturb", 2 to disable it
	Action int32 `json:"action" binding:"oneof=1 2"`
}

type TalkSessionDisturbResponse struct {
}

// TalkSessionListResponse represents the response parameters for the session list API.
type TalkSessionListResponse struct {
	// Items is a list of talk session items
	Items []*TalkSessionItem `json:"items"`
}

// TalkSessionClearUnreadNumRequest represents the request parameters for clearing the unread message count of a talk session.
type TalkSessionClearUnreadNumRequest struct {
	// TalkMode indicates the chat mode (e.g., 1 for private chat, 2 for group chat)
	TalkMode int32 `json:"talk_mode" binding:"required,oneof=1 2"`

	// ToFromID represents the ID of the other participant in the chat
	ToFromID int32 `json:"to_from_id" binding:"required"`
}

type TalkSessionClearUnreadNumResponse struct {
}

type GetTalkRecordsRequest struct {
	TalkMode int `form:"talk_mode" json:"talk_mode" binding:"required,oneof=1 2"`       // 对话类型
	ToFromId int `form:"to_from_id" json:"to_from_id" binding:"required,numeric,min=1"` // 接收者ID
	MsgType  int `form:"msg_type" json:"msg_type" binding:"numeric"`                    // 消息类型
	Cursor   int `form:"cursor" json:"cursor" binding:"min=0,numeric"`                  // 上次查询的游标
	Limit    int `form:"limit" json:"limit" binding:"required,numeric,max=100"`         // 数据行数
}

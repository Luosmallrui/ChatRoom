package types

type UserSettingResponse struct {
	UserInfo *UserInfo   `json:"user_info"` // 用户信息
	Setting  *ConfigInfo `json:"setting"`   // 配置信息
}

// UserInfo 表示用户信息
type UserInfo struct {
	Uid      int32  `json:"uid"`      // 用户 ID
	Nickname string `json:"nickname"` // 昵称
	Avatar   string `json:"avatar"`   // 头像
	Motto    string `json:"motto"`    // 个性签名
	Gender   int32  `json:"gender"`   // 性别
	IsQiye   bool   `json:"is_qiye"`  // 是否是企业用户
	Mobile   string `json:"mobile"`   // 手机号码
	Email    string `json:"email"`    // 邮箱
}

// ConfigInfo 表示用户配置信息
type ConfigInfo struct {
	ThemeMode           string `json:"theme_mode"`            // 主题模式
	ThemeBagImg         string `json:"theme_bag_img"`         // 背景图片
	ThemeColor          string `json:"theme_color"`           // 主题颜色
	NotifyCueTone       string `json:"notify_cue_tone"`       // 通知提示音
	KeyboardEventNotify string `json:"keyboard_event_notify"` // 键盘事件通知
}

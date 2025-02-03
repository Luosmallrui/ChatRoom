package types

type EmoticonListResponse struct {
	Items []*EmoticonItem `json:"items"` // 表情包列表
}

type EmoticonItem struct {
	EmoticonID int32  `json:"emoticon_id"` // 表情包ID
	URL        string `json:"url"`         // 表情包的URL
}

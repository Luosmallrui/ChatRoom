package types

// UploadAvatarRequest 头像上传接口请求参数
type UploadAvatarRequest struct{}

// UploadAvatarResponse 头像上传接口响应参数
type UploadAvatarResponse struct {
	Avatar string `json:"avatar"` // 头像地址
}

// UploadImageRequest 图片上传接口请求参数
type UploadImageRequest struct{}

// UploadImageResponse 图片上传接口响应参数
type UploadImageResponse struct {
	Src string `json:"src"` // 图片地址
}

// UploadInitiateMultipartRequest 批量上传文件初始化接口请求参数
type UploadInitiateMultipartRequest struct {
	FileName string `json:"file_name" binding:"required"` // 文件名
	FileSize int64  `json:"file_size" binding:"required"` // 文件大小
}

// UploadInitiateMultipartResponse 批量上传文件初始化接口响应参数
type UploadInitiateMultipartResponse struct {
	UploadID  string `json:"upload_id"`  // 上传ID
	ShardSize int32  `json:"shard_size"` // 分片大小
	ShardNum  int32  `json:"shard_num"`  // 分片数量
}

// UploadMultipartRequest 批量上传文件接口请求参数
type UploadMultipartRequest struct {
	UploadID   string `form:"upload_id" binding:"required"`       // 上传ID
	SplitIndex int32  `form:"split_index" binding:"min=1"`        // 分片索引
	SplitNum   int32  `form:"split_num" binding:"required,min=1"` // 分片总数
}

// UploadMultipartResponse 批量上传文件接口响应参数
type UploadMultipartResponse struct {
	UploadID string `json:"upload_id"` // 上传ID
	IsMerge  bool   `json:"is_merge"`  // 是否已合并
}

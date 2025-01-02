package types

// AuthLoginRequest 管理员登录接口请求参数
type AuthLoginRequest struct {
	// 登录账号
	Username string `json:"username" binding:"required"`
	// 登录密码
	Password string `json:"password" binding:"required"`
	// 图形验证码
	Captcha string `json:"captcha" binding:"required"`
	// 图形验证码凭据
	CaptchaVoucher string `json:"captcha_voucher" binding:"required"`
}

// AuthLoginResponse 管理员登录接口响应参数
type AuthLoginResponse struct {
	Auth AccessToken `json:"auth"`
}

// AccessToken 包含授权信息
type AccessToken struct {
	// Token 类型
	Type string `json:"type"`
	// Token
	AccessToken string `json:"access_token"`
	// 过期时间（秒）
	ExpiresIn int32 `json:"expires_in"`
}

// AuthCaptchaResponse 图形验证接口响应参数
type AuthCaptchaResponse struct {
	// 验证码唯一凭证
	Voucher string `json:"voucher"`
	// 验证码图像 base64
	Captcha string `json:"captcha"`
}

package types

// AuthLoginRequest 管理员登录接口请求参数
type AuthLoginRequest struct {
	// 登录账号
	Username string `json:"username"`
	// 登录密码
	Password string `json:"password" binding:"required"`
	// 图形验证码
	Captcha string `json:"captcha"`
	// 图形验证码凭据
	CaptchaVoucher string `json:"captcha_voucher"`

	Mobile string `json:"mobile"`
}

// AuthLoginResponse 管理员登录接口响应参数
type AuthLoginResponse struct {
	AccessToken
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

// AuthRegisterRequest 注册请求参数
type AuthRegisterRequest struct {
	// 用户昵称
	Nickname string `json:"nickname"`
	// 登录手机号
	Mobile string `json:"mobile"`
	// 登录密码
	Password string `json:"password"`
	// 登录平台
	Platform string `json:"platform"`
	// 短信验证码
	SmsCode string `json:"sms_code"`
	Email   string `json:"email"`
}

type AuthRegisterResponse struct {
}

// AuthRefreshResponse 刷新令牌响应参数
type AuthRefreshResponse struct {
	// Token 类型
	Type string `json:"type"`
	// 访问 Token
	AccessToken string `json:"access_token"`
	// Token 过期时间（秒）
	ExpiresIn int32 `json:"expires_in"`
}

// AuthForgetRequest 忘记密码请求参数
type AuthForgetRequest struct {
	// 手机号
	Mobile string `json:"mobile" binding:"required,len=11,phone"`
	// 登录密码
	Password string `json:"password" binding:"required,min=6,max=16"`
	// 短信验证码
	SmsCode string `json:"sms_code" binding:"required"`
}

type AuthForgetResponse struct{}

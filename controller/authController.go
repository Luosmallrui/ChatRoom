package controller

import (
	"chatroom/context"
	"chatroom/dao"
	jwt "chatroom/middleware"
	"chatroom/model"
	"chatroom/pkg/encrypt"
	"chatroom/types"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"

	"chatroom/config"
	"github.com/mojocn/base64Captcha"
	"gorm.io/gorm"
)

type AuthController struct {
	Config          *config.Config
	AdminRepo       *dao.Admin
	JwtTokenStorage *dao.JwtTokenStorage
	ICaptcha        *base64Captcha.Captcha
}

func (u *AuthController) RegisterRouter(r gin.IRouter) {
	auth := r.Group("/auth")
	auth.POST("/login", context.HandlerFunc(u.Login)) // 登录
	//auth.POST("/register", context.HandlerFunc(u.R))          // 注册
	auth.POST("/refresh", context.HandlerFunc(u.Refresh)) // 刷新 Token
	auth.POST("/logout", context.HandlerFunc(u.Logout))   // 退出登录
	//auth.POST("/forget", context.HandlerFunc(handler.V1.Auth.Forget)) // 找回密码
}

// Login 登录接口
func (c *AuthController) Login(ctx *context.Context) error {

	var in types.AuthLoginRequest
	if err := ctx.Context.ShouldBindJSON(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	if !c.ICaptcha.Verify(in.CaptchaVoucher, in.Captcha, true) {
		return ctx.InvalidParams("验证码填写不正确")
	}

	adminInfo, err := c.AdminRepo.FindByWhere(ctx.Ctx(), "username = ? or email = ?", in.Username, in.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.InvalidParams("账号不存在或密码填写错误!")
		}

		return ctx.Error(err.Error())
	}

	password, err := encrypt.RsaDecrypt(in.Password, c.Config.App.PrivateKey)
	if err != nil {
		return ctx.Error(err.Error())
	}

	if !encrypt.VerifyPassword(adminInfo.Password, string(password)) {
		return ctx.InvalidParams("账号不存在或密码填写错误!")
	}

	if adminInfo.Status != model.AdminStatusNormal {
		return ctx.ErrorBusiness("账号已被管理员禁用，如有问题请联系管理员！")
	}

	expiresAt := time.Now().Add(12 * time.Hour)

	// 生成登录凭证
	token := jwt.GenerateToken("admin", c.Config.Jwt.Secret, &jwt.Options{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		ID:        strconv.Itoa(adminInfo.Id),
		Issuer:    "im.admin",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	return ctx.Success(&types.AuthLoginResponse{
		Auth: types.AccessToken{
			Type:        "Bearer",
			AccessToken: token,
			ExpiresIn:   int32(expiresAt.Unix() - time.Now().Unix()),
		},
	})
}

// Captcha 图形验证码
func (c *AuthController) Captcha(ctx *context.Context) error {
	voucher, captcha, _, err := c.ICaptcha.Generate()
	if err != nil {
		return ctx.ErrorBusiness(err)
	}

	return ctx.Success(&types.AuthCaptchaResponse{
		Voucher: voucher,
		Captcha: captcha,
	})
}

// Logout 退出登录接口
func (c *AuthController) Logout(ctx *context.Context) error {

	session := ctx.JwtSession()
	if session != nil {
		if ex := session.ExpiresAt - time.Now().Unix(); ex > 0 {
			_ = c.JwtTokenStorage.SetBlackList(ctx.Ctx(), session.Token, time.Duration(ex)*time.Second)
		}
	}

	return ctx.Success(nil)
}

// Refresh Token 刷新接口
func (c *AuthController) Refresh(ctx *context.Context) error {

	// TODO 业务逻辑 ...

	return ctx.Success(nil)
}

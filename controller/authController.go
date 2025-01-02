package controller

import (
	"chatroom/context"
	"chatroom/dao"
	jwt "chatroom/middleware"
	"chatroom/pkg/encrypt"
	"chatroom/service"
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
	UserRepo        *dao.Users
	JwtTokenStorage *dao.JwtTokenStorage
	ICaptcha        *base64Captcha.Captcha
	UserService     service.IUserService
}

func (u *AuthController) RegisterRouter(r gin.IRouter) {
	auth := r.Group("/auth")
	auth.POST("/login", context.HandlerFunc(u.Login))       // 登录
	auth.POST("/register", context.HandlerFunc(u.Register)) // 注册
	auth.POST("/refresh", context.HandlerFunc(u.Refresh))   // 刷新 Token
	auth.POST("/logout", context.HandlerFunc(u.Logout))     // 退出登录
	auth.POST("/forget", context.HandlerFunc(u.Forget))     // 找回密码
}

// Login 登录接口
func (u *AuthController) Login(ctx *context.Context) error {

	var in types.AuthLoginRequest
	if err := ctx.Context.ShouldBindJSON(&in); err != nil {
		return ctx.InvalidParams(err)
	}

	//if !u.ICaptcha.Verify(in.CaptchaVoucher, in.Captcha, true) {
	//	return ctx.InvalidParams("验证码填写不正确")
	//}

	adminInfo, err := u.UserRepo.FindByWhere(ctx.Ctx(), "nickname = ? ", in.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.InvalidParams("账号不存在或密码填写错误!")
		}

		return ctx.Error(err.Error())
	}
	if !encrypt.VerifyPassword(adminInfo.Password, in.Password) {
		return ctx.InvalidParams("账号不存在或密码填写错误!")
	}

	//if adminInfo.Status != model.AdminStatusNormal {
	//	return ctx.ErrorBusiness("账号已被管理员禁用，如有问题请联系管理员！")
	//}

	expiresAt := time.Now().Add(12 * time.Hour)

	// 生成登录凭证
	token := jwt.GenerateToken("admin", u.Config.Jwt.Secret, &jwt.Options{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		ID:        strconv.Itoa(adminInfo.Id),
		Issuer:    "im.admin",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	return ctx.Success(&types.AuthLoginResponse{
		AccessToken: types.AccessToken{
			Type:        "Bearer",
			AccessToken: token,
			ExpiresIn:   int32(expiresAt.Unix() - time.Now().Unix()),
		},
	})
}

// Captcha 图形验证码
func (u *AuthController) Captcha(ctx *context.Context) error {
	voucher, captcha, _, err := u.ICaptcha.Generate()
	if err != nil {
		return ctx.ErrorBusiness(err)
	}

	return ctx.Success(&types.AuthCaptchaResponse{
		Voucher: voucher,
		Captcha: captcha,
	})
}

// Logout 退出登录接口
func (u *AuthController) Logout(ctx *context.Context) error {

	u.toBlackList(ctx)

	return ctx.Success(nil)
}

// Register 注册接口
func (u *AuthController) Register(ctx *context.Context) error {
	in := &types.AuthRegisterRequest{}
	if err := ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	//// 验证短信验证码是否正确
	//if !u.SmsService.Verify(ctx.Ctx(), entity.SmsRegisterChannel, in.Mobile, in.SmsCode) {
	//	return ctx.InvalidParams("短信验证码填写错误！")
	//}

	if _, err := u.UserService.Register(ctx.Ctx(), &service.UserRegisterOpt{
		Nickname: in.Nickname,
		Mobile:   in.Mobile,
		Password: in.Password,
		Platform: in.Platform,
		Email:    in.Email,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	//u.SmsService.Delete(ctx.Ctx(), entity.SmsRegisterChannel, in.Mobile)

	return ctx.Success(&types.AuthRegisterResponse{})
}

// Refresh Token 刷新接口
func (u *AuthController) Refresh(ctx *context.Context) error {

	u.toBlackList(ctx)

	return ctx.Success(&types.AuthRefreshResponse{
		Type:        "Bearer",
		AccessToken: u.token(ctx.UserId()),
		ExpiresIn:   int32(u.Config.Jwt.ExpiresTime),
	})
}

// 设置黑名单
func (u *AuthController) toBlackList(ctx *context.Context) {
	session := ctx.JwtSession()
	if session != nil {
		if ex := session.ExpiresAt - time.Now().Unix(); ex > 0 {
			_ = u.JwtTokenStorage.SetBlackList(ctx.Ctx(), session.Token, time.Duration(ex)*time.Second)
		}
	}
}

func (u *AuthController) token(uid int) string {

	expiresAt := time.Now().Add(time.Second * time.Duration(u.Config.Jwt.ExpiresTime))

	// 生成登录凭证
	token := jwt.GenerateToken("api", u.Config.Jwt.Secret, &jwt.Options{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		ID:        strconv.Itoa(uid),
		Issuer:    "im.web",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	return token
}

// Forget 账号找回接口
func (u *AuthController) Forget(ctx *context.Context) error {
	in := &types.AuthForgetRequest{}
	if err := ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	//// 验证短信验证码是否正确
	//if !c.SmsService.Verify(ctx.Ctx(), entity.SmsForgetAccountChannel, in.Mobile, in.SmsCode) {
	//	return ctx.InvalidParams("短信验证码填写错误！")
	//}

	if _, err := u.UserService.Forget(&service.UserForgetOpt{
		Mobile:   in.Mobile,
		Password: in.Password,
		SmsCode:  in.SmsCode,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	//c.SmsService.Delete(ctx.Ctx(), entity.SmsForgetAccountChannel, in.Mobile)

	return ctx.Success(&types.AuthForgetResponse{})
}

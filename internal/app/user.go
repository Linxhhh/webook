package app

import (
	"github.com/Linxhhh/webook/internal/domain"
	"github.com/Linxhhh/webook/internal/service"
	"github.com/Linxhhh/webook/pkg/jwts"
	"github.com/Linxhhh/webook/pkg/res"
	"github.com/gin-gonic/gin"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
)

const (
	emailRegexPattern    = `^\w+([-+.]\\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?.&])[A-Za-z\d@$!%*#?.&]{8,}$`
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (hdl *UserHandler) RegistryRouter(router *gin.Engine) {
	ug := router.Group("user")
	ug.POST("signup", hdl.SignUp)    // 用户注册
	ug.POST("login", hdl.LoginByJWT) // 用户登录
}

/*
SignUp 用户注册API：
先绑定前端的注册请求，再进行邮箱校验，密码校验，最后调用下层服务来创建用户
*/
func (hdl *UserHandler) SignUp(ctx *gin.Context) {

	// 注册请求
	type SignUpReq struct {
		Email           string `json:"email" binding:"required"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirmPassword" binding:"required"`
	}
	var req SignUpReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.String(200, "填写信息不完整")
		return
	}

	// 校验邮箱
	if ok, err := IsValidEmail(req.Email); err != nil {
		ctx.String(200, "系统错误")
		return
	} else if !ok {
		ctx.String(200, "非法邮箱格式")
		return
	}

	// 校验密码
	if req.Password != req.ConfirmPassword {
		ctx.String(200, "两次密码不一致")
		return
	}
	if ok, err := IsValidPassword(req.Password); err != nil {
		ctx.String(200, "系统错误")
		return
	} else if !ok {
		ctx.String(200, "密码必须包含字母、数字、特殊字符，并且不少于八位")
		return
	}

	// 调用服务
	err := hdl.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		ctx.String(200, "注册成功")
	case service.ErrDuplicateEmail:
		ctx.String(200, "邮箱冲突")
	default:
		ctx.String(200, "系统错误")
	}
}

// IsValidEmail 通过正则表达式校验邮箱格式
func IsValidEmail(email string) (bool, error) {
	var emailRegex = regexp.MustCompile(emailRegexPattern, regexp.None)
	return emailRegex.MatchString(email)
}

// IsValidPassword 通过正则表达式校验密码格式
func IsValidPassword(pwd string) (bool, error) {
	var passwordRegex = regexp.MustCompile(passwordRegexPattern, regexp.None)
	return passwordRegex.MatchString(pwd)
}

/*
LoginBySession 用户登录API：
先绑定前端的登录请求，再调用下层服务进行校验，最后设置 Session
*/
func (hdl *UserHandler) LoginBySession(ctx *gin.Context) {

	// 登录请求
	type LoginReq struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.String(200, "系统错误")
		return
	}

	// 调用服务
	user, err := hdl.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidEmailOrPassword {
		ctx.String(200, "邮箱或密码错误")
		return
	}
	if err != nil {
		ctx.String(200, "系统错误")
		return
	}

	// 设置 Session
	session := sessions.Default(ctx)
	session.Set("userId", user.Id)
	session.Options(sessions.Options{MaxAge: 3600})
	if err = session.Save(); err != nil {
		ctx.String(200, "系统错误")
		return
	}
	ctx.String(200, "登录成功")
}

/*
LoginByJWT 用户登录API：
先绑定前端的登录请求，再调用下层服务进行校验，最后生成并返回用户 Token
*/
func (hdl *UserHandler) LoginByJWT(ctx *gin.Context) {

	// 登录请求
	type LoginReq struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}

	// 调用服务
	user, err := hdl.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidEmailOrPassword {
		res.FailWithMsg("邮箱或密码错误", ctx)
		return
	}
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}

	// 生成 Token
	token, err := jwts.GenToken(jwts.JwtPayload{
		UserId: user.Id,
		UserAgent: ctx.GetHeader("User-Agent"),
	})
	if err != nil {
		res.FailWithMsg("生成用户令牌错误！", ctx)
		return
	}

	// 返回用户token
	res.OKWithData(token, ctx)
}

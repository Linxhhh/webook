package app

import (
	"time"
	"unicode/utf8"

	"github.com/Linxhhh/webook/internal/domain"
	"github.com/Linxhhh/webook/internal/service"
	"github.com/Linxhhh/webook/pkg/jwts"
	"github.com/Linxhhh/webook/pkg/res"
	"github.com/gin-gonic/gin"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
)

const (
	biz                  = "login"
	emailRegexPattern    = `^\w+([-+.]\\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?.&])[A-Za-z\d@$!%*#?.&]{8,}$`
	phoneRegexPattern    = `^(\+?0?86\-?)?1[345789]\d{9}$`
)

type UserHandler struct {
	svc     *service.UserService
	codeSvc *service.CodeService
}

func NewUserHandler(svc *service.UserService, codeSvc *service.CodeService) *UserHandler {
	return &UserHandler{
		svc: svc,
		codeSvc: codeSvc,
	}
}

func (hdl *UserHandler) RegistryRouter(router *gin.Engine) {
	ug := router.Group("user")
	ug.POST("signup", hdl.SignUp)    // 用户注册
	ug.POST("login", hdl.LoginByJWT) // 用户登录

	ug.PUT("sms/send", hdl.SendSmsCode)      // 短信验证码登录：发送验证码
	ug.POST("sms/verify", hdl.VerifySmsCode) // 短信验证码登录：校验验证码

	ug.POST("edit", hdl.Edit)      // 信息编辑
	ug.GET("profile", hdl.Profile) // 信息获取
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
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 校验邮箱
	if ok, err := isValidEmail(req.Email); err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	} else if !ok {
		res.FailWithMsg("非法邮箱格式", ctx)
		return
	}

	// 校验密码
	if req.Password != req.ConfirmPassword {
		res.FailWithMsg("两次密码不一致", ctx)
		return
	}
	if ok, err := isValidPassword(req.Password); err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	} else if !ok {
		res.FailWithMsg("密码必须包含字母、数字、特殊字符，并且不少于八位", ctx)
		return
	}

	// 调用服务
	err := hdl.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		res.OKWithMsg("注册成功", ctx)
	case service.ErrDuplicateEmailorPhone:
		res.FailWithMsg("邮箱冲突", ctx)
	default:
		res.FailWithMsg("系统错误", ctx)
	}
}

// isValidEmail 通过正则表达式校验邮箱格式
func isValidEmail(email string) (bool, error) {
	var emailRegex = regexp.MustCompile(emailRegexPattern, regexp.None)
	return emailRegex.MatchString(email)
}

// isValidPassword 通过正则表达式校验密码格式
func isValidPassword(pwd string) (bool, error) {
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
		res.FailWithMsg("参数错误", ctx)
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

	// 设置 Session
	session := sessions.Default(ctx)
	session.Set("userId", user.Id)
	session.Options(sessions.Options{MaxAge: 3600})
	if err = session.Save(); err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}
	res.OKWithMsg("登录成功", ctx)
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
		res.FailWithMsg("参数错误", ctx)
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
		UserId:    user.Id,
		UserAgent: ctx.GetHeader("User-Agent"),
	})
	if err != nil {
		res.FailWithMsg("生成用户令牌错误！", ctx)
		return
	}

	// 返回用户token
	ctx.Header("jwt-token", token)
	res.OKWithMsg("登陆成功", ctx)
}

/*
SendSmsCode 发送短信验证码API：
绑定前端手机号，调用底层服务发送短信
*/
func (hdl *UserHandler) SendSmsCode(ctx *gin.Context) {

	// 发送验证码请求
	type SendSmsCodeReq struct {
		Phone string `json:"phone" binding:"required"`
	}
	var req SendSmsCodeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 校验手机号码
	if ok, err := isValidPhone(req.Phone); err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	} else if !ok {
		res.FailWithMsg("非法手机号码", ctx)
		return
	}

	// 调用底层服务
	err := hdl.codeSvc.Send(ctx, biz, req.Phone)
	switch err {
	case nil:
		res.OKWithMsg("短信发送成功", ctx)
	case service.ErrSendCodeTooMany:
		res.FailWithMsg("短信发送频繁", ctx)
	default:
		res.FailWithMsg("系统错误", ctx)
	}
}

// isValidPhone 通过正则表达式校验手机格式
func isValidPhone(phone string) (bool, error) {
	var emailRegex = regexp.MustCompile(phoneRegexPattern, regexp.None)
	return emailRegex.MatchString(phone)
}

/*
VerifySmsCode 短信验证API：
调用下次服务，对验证码进行校验，若校验通过，则查询/创建用户，返回用户 Token
*/
func (hdl *UserHandler) VerifySmsCode(ctx *gin.Context) {

	// 验证码校验请求
	type VerifySmsCodeReq struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	var req VerifySmsCodeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 验证码校验
	err := hdl.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	switch err {
	case nil:
	case service.ErrVerifyCodeFailed:
		res.FailWithMsg("校验失败", ctx)
		return
	case service.ErrVerifyCodeTooMany:
		res.FailWithMsg("校验频繁", ctx)
		return
	default:
		res.FailWithMsg("系统错误", ctx)
		return
	}

	// 查找或创建 User
	uid, err := hdl.svc.FindOrCreate(ctx, req.Phone)
	switch err {
	case nil:
	case service.ErrDuplicateEmailorPhone:  // 这种情况是，一个未注册用户，通过验证码同时登录两台设备
	default:
		res.FailWithMsg("系统错误", ctx)
		return
	}

	// 生成 Token
	token, err := jwts.GenToken(jwts.JwtPayload{
		UserId:    uid,
		UserAgent: ctx.GetHeader("User-Agent"),
	})
	if err != nil {
		res.FailWithMsg("生成用户令牌错误！", ctx)
		return
	}

	// 返回用户token
	ctx.Header("jwt-token", token)
	res.OKWithMsg("登陆成功", ctx)
}

/*
edit 信息编辑API：
绑定前端信息，调用下次服务进行存储
*/
func (hdl *UserHandler) Edit(ctx *gin.Context) {

	// 信息修改请求
	type EditReq struct {
		NickName     string `json:"nickName"`
		Birthday     string `json:"birthday"`
		Introduction string `json:"introduction"`
	}
	var req EditReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	user := domain.User{Id: claims.UserId}

	// 校验昵称长度
	if req.NickName != "" {
		if !isValidateNickName(req.NickName) {
			res.FailWithMsg("用户昵称超长", ctx)
			return
		}
		user.NickName = req.NickName
	}

	// 校验日期格式
	var birthday time.Time
	if req.Birthday != "" {
		date, ok := validateDateFormat(req.Birthday)
		if !ok {
			res.FailWithMsg("生日日期格式错误", ctx)
			return
		}
		birthday = date
		user.Birthday = birthday
	}

	// 校验简介长度
	if req.Introduction != "" {
		if !isValidateIntroduction(req.Introduction) {
			res.FailWithMsg("个人简介超长", ctx)
			return
		}
		user.Introduction = req.Introduction
	}

	// 调用下层服务
	err := hdl.svc.Edit(ctx, user)
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}
	res.OKWithMsg("编辑成功", ctx)
}

// isValidateNickName 判断昵称长度
func isValidateNickName(name string) bool {
	const MaxLength, MaxBytes = 10, 20
	return utf8.RuneCountInString(name) < MaxLength || len(name) < MaxBytes
}

// validateDateFormat 校验日期格式是否为 "2006-01-02"
func validateDateFormat(dateStr string) (time.Time, bool) {
	date, err := time.Parse("2006-01-02", dateStr)
	return date, err == nil
}

// isValidateIntroduction 校验简介长度
func isValidateIntroduction(s string) bool {
	const MaxLength, MaxBytes = 50, 100
	return utf8.RuneCountInString(s) < MaxLength || len(s) < MaxBytes
}

/*
Profile 获取用户信息API：
调用下次服务，响应用户信息
*/
func (hdl *UserHandler) Profile(ctx *gin.Context) {

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	// 调用下层服务
	user, err := hdl.svc.Profile(ctx, claims.UserId)
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}

	// 用户信息响应
	type ProfileResp struct {
		Email        string `json:"email"`
		Phone        string `json:"phone"`
		NickName     string `json:"nickName"`
		Birthday     string `json:"birthday"`
		Introduction string `json:"introduction"`
	}
	resp := ProfileResp{
		Email:        user.Email,
		Phone:        user.Phone,
		NickName:     user.NickName,
		Birthday:     user.Birthday.Format("2006-01-02"),
		Introduction: user.Introduction,
	}
	res.OKWithData(resp, ctx)
}

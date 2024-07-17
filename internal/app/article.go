package app

import (
	"errors"
	"log"
	"strconv"

	"github.com/Linxhhh/webook/internal/domain"
	"github.com/Linxhhh/webook/internal/service"
	"github.com/Linxhhh/webook/pkg/jwts"
	"github.com/Linxhhh/webook/pkg/res"
	"github.com/gin-gonic/gin"
)

var ErrIncorrectArticleorAuthor = service.ErrIncorrectArticleorAuthor

type ArticleHandler struct {
	svc      *service.ArticleService
	interSvc *service.InteractionService
	biz      string
}

func NewArticleHandler(svc *service.ArticleService, interSvc *service.InteractionService) *ArticleHandler {
	return &ArticleHandler{
		svc:      svc,
		interSvc: interSvc,
		biz:      "article",
	}
}

func (hdl *ArticleHandler) RegistryRouter(router *gin.Engine) {
	// 作者接口
	ag := router.Group("article")
	ag.POST("edit", hdl.Edit)
	ag.POST("publish", hdl.Publish)
	ag.DELETE("withdraw", hdl.Withdraw)
	ag.GET("list", hdl.List)
	ag.GET("detail", hdl.Detail)

	// 读者接口
	pg := router.Group("pub")
	pg.GET("detail", hdl.PubDetail, hdl.Read)
	pg.POST("like", hdl.Like)
	pg.POST("collect", hdl.Collect)
	pg.GET("interaction", hdl.Interaction)
}

// Edit 新建帖子，或编辑旧帖子
func (hdl *ArticleHandler) Edit(ctx *gin.Context) {

	// 绑定参数
	var req ArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	// 调用下层服务
	aid, err := hdl.svc.Save(ctx, domain.Article{
		Id:       req.Id,
		Title:    req.Title,
		Content:  req.Content,
		AuthorId: claims.UserId,
	})
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}
	res.OKWithData(gin.H{"article_id": aid}, ctx)
}

// Pubish 帖子发表
func (hdl *ArticleHandler) Publish(ctx *gin.Context) {

	// 绑定参数
	var req ArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	// 调用下层服务
	aid, err := hdl.svc.Publish(ctx, domain.Article{
		Id:       req.Id,
		Title:    req.Title,
		Content:  req.Content,
		AuthorId: claims.UserId,
	})
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}
	res.OKWithData(gin.H{"article_id": aid}, ctx)
}

// Withdraw 撤销发表
func (hdl *ArticleHandler) Withdraw(ctx *gin.Context) {

	// 绑定参数
	type Req struct{ Id int64 }
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	// 调用下层服务
	err := hdl.svc.Withdraw(ctx, claims.UserId, req.Id)
	if err != nil {
		if errors.Is(err, ErrIncorrectArticleorAuthor) {
			res.FailWithMsg("非法撤销", ctx)
			return
		}
		res.FailWithMsg("系统错误", ctx)
		return
	}
	res.OKWithMsg("撤销成功", ctx)
}

type ArticleRequest struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// List 获取用户制作库的帖子列表
func (hdl *ArticleHandler) List(ctx *gin.Context) {

	// 绑定参数
	type ListReq struct {
		Page     int `json:"page"`
		PageSize int `json:"pageSize"`
	}
	var req ListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	// 调用下层服务
	articleList, err := hdl.svc.List(ctx, claims.UserId, req.Page, req.PageSize)
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}

	// 返回响应
	res.OKWithData(articleList, ctx)
}

// Detail 获取制作库的帖子详情
func (hdl *ArticleHandler) Detail(ctx *gin.Context) {

	// 绑定参数
	aid, err := strconv.ParseInt(ctx.Query("id"), 10, 64)
	if aid == 0 || err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	// 调用下层服务
	art, err := hdl.svc.Detail(ctx, claims.UserId, aid)
	if err != nil {
		if errors.Is(err, ErrIncorrectArticleorAuthor) {
			res.FailWithMsg("非法查询", ctx)
			return
		}
		res.FailWithMsg("系统错误", ctx)
		return
	}

	// 返回响应
	res.OKWithData(art, ctx)
}

// PubDetail 获取线上库的帖子详情
func (hdl *ArticleHandler) PubDetail(ctx *gin.Context) {

	// 绑定参数
	aid, err := strconv.ParseInt(ctx.Query("id"), 10, 64)
	if aid == 0 || err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 调用下层服务
	art, err := hdl.svc.PubDetail(ctx, aid)
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}

	// 设置上下文，然后返回响应
	ctx.Set("article_id", aid)
	res.OKWithData(art, ctx)
}

/*
阅读、点赞、收藏功能待优化
*/

func (hdl *ArticleHandler) Read(ctx *gin.Context) {

	aid, exists := ctx.Get("article_id")
	if !exists {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 调用下层服务
	err := hdl.interSvc.IncrReadCnt(ctx, hdl.biz, aid.(int64))
	if err != nil {
		log.Panicln("IncrReadCnt 报错：err : ", err.Error())
	}
}

func (hdl *ArticleHandler) Like(ctx *gin.Context) {

	// 绑定参数
	type Req struct {
		Id   int64 `json:"id"`
		Like bool  `json:"like"` // true 表示点赞，false 表示取消
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	var err error
	if req.Like {
		err = hdl.interSvc.Like(ctx, hdl.biz, req.Id, claims.UserId)
	} else {
		err = hdl.interSvc.CancelLike(ctx, hdl.biz, req.Id, claims.UserId)
	}
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}
	res.OKWithMsg("操作成功", ctx)
}

func (hdl *ArticleHandler) Collect(ctx *gin.Context) {
	// 绑定参数
	type Req struct {
		Id      int64 `json:"id"`
		Collect bool  `json:"collect"` // true 表示点赞，false 表示取消
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	var err error
	if req.Collect {
		err = hdl.interSvc.Collect(ctx, hdl.biz, req.Id, claims.UserId)
	} else {
		err = hdl.interSvc.CancelCollect(ctx, hdl.biz, req.Id, claims.UserId)
	}
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}
	res.OKWithMsg("操作成功", ctx)
}

func (hdl *ArticleHandler) CollectionList(ctx *gin.Context) {
	
	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	// 调用下层服务
	arts, err := hdl.interSvc.CollectionList(ctx, hdl.biz, claims.UserId)
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}
	res.OKWithData(arts, ctx)
}

func (hdl *ArticleHandler) Interaction(ctx *gin.Context) {
	
	// 绑定参数
	aid, err := strconv.ParseInt(ctx.Query("id"), 10, 64)
	if aid == 0 || err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	// 调用下层服务
	i, err := hdl.interSvc.Get(ctx, hdl.biz, aid, claims.UserId)
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}

	// 返回响应
	res.OKWithData(i, ctx)
}


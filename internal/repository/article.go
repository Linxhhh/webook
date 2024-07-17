package repository

import (
	"context"
	"time"

	"github.com/Linxhhh/webook/internal/domain"
	"github.com/Linxhhh/webook/internal/repository/cache"
	"github.com/Linxhhh/webook/internal/repository/dao"
)

var ErrIncorrectArticleorAuthor = dao.ErrIncorrectArticleorAuthor

type ArticleRepository interface {
	Insert(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	Sync(ctx context.Context, article domain.Article) (int64, error)
	SyncStatus(ctx context.Context, uid int64, aid int64, status domain.ArticleStatus) error
	GetListByAuthor(ctx context.Context, uid int64, offset, limit int) ([]domain.ArticleListElem, error)
	GetById(ctx context.Context, aid int64) (domain.Article, error)
	GetPubById(ctx context.Context, aid int64) (domain.Article, error)
}

type CacheArticleRepository struct {
	dao   dao.ArticleDAO
	cache cache.ArticleCache
}

func NewArticleRepository(dao dao.ArticleDAO, cache cache.ArticleCache) ArticleRepository {
	return &CacheArticleRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *CacheArticleRepository) Insert(ctx context.Context, article domain.Article) (int64, error) {

	aid, err := repo.dao.Insert(ctx, dao.Article{
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.AuthorId,
		Status:   uint8(article.Status),
	})
	if err == nil {
		// 清除首页缓存
		repo.cache.DelFirstPage(ctx, article.AuthorId)
	}
	return aid, err
}

func (repo *CacheArticleRepository) Update(ctx context.Context, article domain.Article) error {

	// 清除首页缓存
	repo.cache.DelFirstPage(ctx, article.AuthorId)

	err := repo.dao.Update(ctx, dao.Article{
		Id:      article.Id,
		Title:   article.Title,
		Content: article.Content,
		Status:  uint8(article.Status),
	})
	if err == nil {
		// 清除首页缓存
		repo.cache.DelFirstPage(ctx, article.AuthorId)
	}
	return err
}

func (repo *CacheArticleRepository) Sync(ctx context.Context, article domain.Article) (int64, error) {

	aid, err := repo.dao.Sync(ctx, dao.Article{
		Id:       article.Id,
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.AuthorId,
		Status:   uint8(article.Status),
	})
	if err == nil {
		go func() {
			// 清除首页缓存
			repo.cache.DelFirstPage(ctx, article.AuthorId)
			// 设置帖子缓存
			repo.cache.SetPub(ctx, article)
		}()
	}
	return aid, err
}

func (repo *CacheArticleRepository) SyncStatus(ctx context.Context, uid int64, aid int64, status domain.ArticleStatus) error {

	err := repo.dao.SyncStatus(ctx, uid, aid, uint8(status))
	if err == nil {
		// 清除首页缓存
		repo.cache.DelFirstPage(ctx, uid)
	}
	return err
}

func (repo *CacheArticleRepository) GetListByAuthor(ctx context.Context, uid int64, offset, limit int) ([]domain.ArticleListElem, error) {

	// 如果是首页，则先查询缓存
	if offset == 0 {
		res, err := repo.cache.GetFirstPage(ctx, uid)
		if err == nil {
			return res, err
		}
	}

	// 查询数据库
	arts, err := repo.dao.GetListByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}

	// 类型转换
	var articleList []domain.ArticleListElem
	for _, art := range arts {
		article := domain.ArticleListElem{
			Id:       art.Id,
			Title:    art.Title,
			Abstract: domain.Abstract(art.Content),
			Ctime:    time.UnixMilli(art.Ctime),
			Utime:    time.UnixMilli(art.Utime),
			Status:   domain.ArticleStatus(art.Status),
		}
		articleList = append(articleList, article)
	}

	// 回写缓存
	if offset == 0 && len(articleList) > 0 {
		go func() {
			// 缓存首页
			repo.cache.SetFirstPage(ctx, uid, articleList)

			// 预加载第一个帖子
			const size = 1024 * 1024
			if len(arts[0].Content) < size {
				article := domain.Article{
					Id:       arts[0].Id,
					Title:    arts[0].Title,
					Content:  arts[0].Content,
					AuthorId: arts[0].AuthorId,
					Ctime:    time.UnixMilli(arts[0].Ctime),
					Utime:    time.UnixMilli(arts[0].Utime),
					Status:   domain.ArticleStatus(arts[0].Status),
				}
				repo.cache.Set(ctx, article)
			}
		}()
	}

	return articleList, err
}

func (repo *CacheArticleRepository) GetById(ctx context.Context, aid int64) (domain.Article, error) {

	// 查询缓存
	article, err := repo.cache.Get(ctx, aid)
	if err == nil {
		return article, err
	}

	// 查询数据库
	art, err := repo.dao.GetById(ctx, aid)
	if err != nil {
		return domain.Article{}, err
	}
	article = domain.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.AuthorId,
		Ctime:    time.UnixMilli(art.Ctime),
		Utime:    time.UnixMilli(art.Utime),
		Status:   domain.ArticleStatus(art.Status),
	}

	// 回写缓存
	go func() {
		repo.cache.Set(ctx, article)
	}()

	return article, err
}

func (repo *CacheArticleRepository) GetPubById(ctx context.Context, aid int64) (domain.Article, error) {

	// 查询缓存
	article, err := repo.cache.GetPub(ctx, aid)
	if err == nil {
		return article, err
	}

	// 查询数据库
	art, err := repo.dao.GetPubById(ctx, aid)
	if err != nil {
		return domain.Article{}, err
	}
	article = domain.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.AuthorId,
		Ctime:    time.UnixMilli(art.Ctime),
		Utime:    time.UnixMilli(art.Utime),
		Status:   domain.ArticleStatus(art.Status),
	}

	// 回写缓存
	go func() {
		repo.cache.SetPub(ctx, article)
	}()

	return article, err
}

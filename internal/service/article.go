package service

import (
	"context"
	"errors"
	"time"

	"github.com/Linxhhh/webook/internal/domain"
	"github.com/Linxhhh/webook/internal/repository"
)

var ErrIncorrectArticleorAuthor = repository.ErrIncorrectArticleorAuthor

type ArticleService struct {
	repo     repository.ArticleRepository
	userRepo repository.UserRepository
}

func NewArticleService(repo repository.ArticleRepository, userRepo repository.UserRepository) *ArticleService {
	return &ArticleService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (as *ArticleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnpublished
	if art.Id > 0 {
		return art.Id, as.repo.Update(ctx, art)
	}
	return as.repo.Insert(ctx, art)
}

func (as *ArticleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	return as.repo.Sync(ctx, art)
}

func (as *ArticleService) Withdraw(ctx context.Context, uid int64, aid int64) error {
	return as.repo.SyncStatus(ctx, uid, aid, domain.ArticleStatusPrivate)
}

func (as *ArticleService) List(ctx context.Context, uid int64, page, pageSize int) ([]domain.ArticleListElem, error) {
	limit := pageSize
	offset := (page - 1) * pageSize
	return as.repo.GetListByAuthor(ctx, uid, offset, limit)
}

func (as *ArticleService) Detail(ctx context.Context, uid, aid int64) (domain.Article, error) {
	art, err := as.repo.GetById(ctx, aid)
	if err == nil && art.AuthorId != uid {
		return domain.Article{}, ErrIncorrectArticleorAuthor
	}
	return art, err
}

func (as *ArticleService) PubDetail(ctx context.Context, aid int64) (domain.Article, error) {
	art, err := as.repo.GetPubById(ctx, aid)
	if err != nil {
		return domain.Article{}, err
	}

	// 获取 AuthorName
	user, err := as.userRepo.SearchById(ctx, art.AuthorId)
	if err != nil {
		return domain.Article{}, errors.New("查找用户失败")
	}
	art.AuthorName = user.NickName
	return art, nil
}

func (as *ArticleService) PubList(ctx context.Context, uid int64, limit, offset int) ([]domain.Article, error) {
	// 获取一周内的帖子
	startTime := time.Now().Add(-7 * 24 * time.Hour)
	return as.repo.GetPubList(ctx, startTime, offset, limit)
}
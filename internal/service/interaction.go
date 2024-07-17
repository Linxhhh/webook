package service

import (
	"context"
	"sync"

	"github.com/Linxhhh/webook/internal/domain"
	"github.com/Linxhhh/webook/internal/repository"
)

type InteractionService struct {
	repo    repository.InteractionRepository
	artRepo repository.ArticleRepository
}

func NewInteractionService(repo repository.InteractionRepository, artRepo repository.ArticleRepository) *InteractionService {
	return &InteractionService{
		repo:    repo,
		artRepo: artRepo,
	}
}

func (svc *InteractionService) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return svc.repo.IncrReadCnt(ctx, biz, bizId)
}

func (svc *InteractionService) Like(ctx context.Context, biz string, bizId int64, uid int64) error {
	return svc.repo.Like(ctx, biz, bizId, uid)
}

func (svc *InteractionService) CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	return svc.repo.CancelLike(ctx, biz, bizId, uid)
}

func (svc *InteractionService) Collect(ctx context.Context, biz string, bizId int64, uid int64) error {
	return svc.repo.Collect(ctx, biz, bizId, uid)
}

func (svc *InteractionService) CancelCollect(ctx context.Context, biz string, bizId int64, uid int64) error {
	return svc.repo.CancelCollect(ctx, biz, bizId, uid)
}

/*
后续优化：分页查询
*/
func (svc *InteractionService) CollectionList(ctx context.Context, biz string, uid int64) ([]domain.Article, error) {
	
	// 获取收藏帖子的 aid
	aids, err := svc.repo.GetCollectionList(ctx, biz, uid)
	if err != nil {
		return nil, err
	}

	// 获取收藏的帖子
	var arts []domain.Article
	for _, aid := range aids {
		art, err := svc.artRepo.GetById(ctx, aid)
		if err != nil {
			return nil, err
		}
		arts = append(arts, art)
	}
	return arts, nil
}

func (svc *InteractionService) Get(ctx context.Context, biz string, bizId int64, uid int64) (domain.Interaction, error) {

	// 获取（阅读、点赞、收藏）数据
	i, err := svc.repo.Get(ctx, biz, bizId)
	if err != nil {
		return domain.Interaction{}, err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// 创建一个 error channel
	errCh := make(chan error, 2)

	// 是否已经点赞
	go func() {
		defer wg.Done()
		i.IsLiked, err = svc.repo.GetLike(ctx, biz, bizId, uid)
		if err != nil {
			errCh <- err // 发送错误到 channel
		}
	}()

	// 是否已经收藏
	go func() {
		defer wg.Done()
		i.IsCollected, err = svc.repo.GetCollection(ctx, biz, bizId, uid)
		if err != nil {
			errCh <- err // 发送错误到 channel
		}
	}()

	wg.Wait()

	// 检查协程中是否有错误（只需要记录第一个错误）
	for err := range errCh {
		if err != nil {
			close(errCh)
			return domain.Interaction{}, err
		}
	}
	close(errCh)

	return i, nil
}

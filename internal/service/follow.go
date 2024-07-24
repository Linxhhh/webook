package service

import (
	"context"

	"github.com/Linxhhh/webook/internal/domain"
	"github.com/Linxhhh/webook/internal/repository"
)

type FollowService struct {
	repo repository.FollowRepository
}

func NewFollowService(repo repository.FollowRepository) *FollowService {
	return &FollowService{
		repo: repo,
	}
}

func (svc *FollowService) Follow(ctx context.Context, follower_id, followee_id int64) error {
	return svc.repo.Follow(ctx, follower_id, followee_id)
}

func (svc *FollowService) CancelFollow(ctx context.Context, follower_id, followee_id int64) error {
	return svc.repo.CancelFollow(ctx, follower_id, followee_id)
}

func (svc *FollowService) GetFollowData(ctx context.Context, uid int64) (domain.FollowData, error) {
	return svc.repo.GetFollowData(ctx, uid)
}

func (svc *FollowService) GetFollowed(ctx context.Context, follower_id, followee_id int64) (bool, error) {
	return svc.repo.GetFollowed(ctx, follower_id, followee_id)
}

/*
func (svc *FollowService) GetFolloweeList(ctx context.Context, follower_id, limit, offset int64) ([]domain.User, error) {
	list, err := svc.repo.GetFolloweeList(ctx, follower_id, limit, offset)
}
*/
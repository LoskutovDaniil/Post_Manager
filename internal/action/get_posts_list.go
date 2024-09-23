package action

import (
	"context"
	"fmt"

	"github.com/LoskutovDaniil/OzonTestTask2024/internal/model"
)

//go:generate mockgen -source=get_posts_list.go -destination=mock/get_posts_list_mock.go -package=mock

type GetPostsListRepository interface {
	GetPostsList(ctx context.Context, filter model.PostsFilter) ([]model.Post, error)
}

type GetPostsList struct {
	repo GetPostsListRepository
}

func NewGetPostsList(repo GetPostsListRepository) GetPostsList {
	return GetPostsList{repo}
}

func (act GetPostsList) Do(ctx context.Context, filter model.PostsFilter) ([]model.Post, error) {
	filter = filter.SetDefault()

	post, err := act.repo.GetPostsList(ctx, filter.SetDefault())
	if err != nil {
		return nil, fmt.Errorf("get posts list: %w", err)
	}

	return post, nil
}

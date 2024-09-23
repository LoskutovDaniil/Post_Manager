package action

import (
	"context"
	"fmt"

	"github.com/LoskutovDaniil/OzonTestTask2024/internal/model"
)

//go:generate mockgen -source=get_comments_list.go -destination=mock/get_comments_list_mock.go -package=mock

type GetCommentsListRepository interface {
	GetCommentsList(ctx context.Context, filter model.CommentsFilter) ([]model.Comment, error)
}

type GetCommentsList struct {
	repo GetCommentsListRepository
}

func NewGetCommentsList(repo GetCommentsListRepository) GetCommentsList {
	return GetCommentsList{repo}
}

func (act GetCommentsList) Do(ctx context.Context, filter model.CommentsFilter) ([]model.Comment, error) {
	filter = filter.SetDefault()

	comments, err := act.repo.GetCommentsList(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get comments list: %w", err)
	}

	return comments, nil
}

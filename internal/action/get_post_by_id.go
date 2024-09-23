package action

import (
	"context"
	"fmt"

	"github.com/LoskutovDaniil/OzonTestTask2024/internal/model"

	"github.com/google/uuid"
)

//go:generate mockgen -source=get_post_by_id.go -destination=mock/get_post_by_id_mock.go -package=mock

type GetPostByIDRepository interface {
	GetPostByID(ctx context.Context, id uuid.UUID) (model.Post, error)
}

type GetPostByID struct {
	repo GetPostByIDRepository
}

func NewGetPostByID(repo GetPostByIDRepository) GetPostByID {
	return GetPostByID{repo}
}

func (act GetPostByID) Do(ctx context.Context, id uuid.UUID) (model.Post, error) {
	post, err := act.repo.GetPostByID(ctx, id)
	if err != nil {
		return model.Post{}, fmt.Errorf("get post by id: %w", err)
	}

	return post, nil
}

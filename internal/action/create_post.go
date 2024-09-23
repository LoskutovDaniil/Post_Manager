package action

import (
	"context"
	"fmt"
	"time"

	"github.com/LoskutovDaniil/OzonTestTask2024/internal/model"

	"github.com/google/uuid"
)

//go:generate mockgen -source=create_post.go -destination=mock/create_post_mock.go -package=mock

type CreatePostRepository interface {
	CreatePost(ctx context.Context, post model.Post) error
}

type CreatePost struct {
	repo CreatePostRepository
}

func NewCreatePost(repo CreatePostRepository) CreatePost {
	return CreatePost{repo}
}

func (act CreatePost) Do(ctx context.Context, newPost model.NewPost) (model.Post, error) {
	post := model.Post{
		ID:              uuid.New(),
		CreatedAt:       time.Now(),
		Text:            newPost.Text,
		DisableComments: newPost.DisableComments,
	}

	err := act.repo.CreatePost(ctx, post)
	if err != nil {
		return model.Post{}, fmt.Errorf("create post: %w", err)
	}

	return post, nil
}

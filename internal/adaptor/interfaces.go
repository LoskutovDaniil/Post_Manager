package adaptor

import (
	"context"

	"github.com/LoskutovDaniil/OzonTestTask2024/internal/model"

	"github.com/google/uuid"
)

type Repository interface {
	CreateComment(ctx context.Context, comment model.Comment) error
	CreatePost(ctx context.Context, post model.Post) error
	GetCommentsList(ctx context.Context, filter model.CommentsFilter) ([]model.Comment, error)
	GetPostByID(ctx context.Context, id uuid.UUID) (model.Post, error)
	GetPostsList(ctx context.Context, filter model.PostsFilter) ([]model.Post, error)
}

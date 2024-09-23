package action

import (
	"context"
	"fmt"
	"time"

	"github.com/LoskutovDaniil/OzonTestTask2024/internal/model"

	"github.com/google/uuid"
)

//go:generate mockgen -source=create_comment.go -destination=mock/create_comment_mock.go -package=mock

type CreateCommentRepository interface {
	GetPostByID(ctx context.Context, id uuid.UUID) (model.Post, error)
	CreateComment(ctx context.Context, comment model.Comment) error
}

type CreateComment struct {
	repo CreateCommentRepository
}

func NewCreateComment(repo CreateCommentRepository) CreateComment {
	return CreateComment{repo}
}

func (act CreateComment) Do(ctx context.Context, newComment model.NewComment) (model.Comment, error) {
	post, err := act.repo.GetPostByID(ctx, newComment.PostID)
	if err != nil {
		return model.Comment{}, fmt.Errorf("get post by ID: %w", err)
	}
	if post.DisableComments {
		return model.Comment{}, model.NewBadRequestErrf(
			"comments disabled for post with ID \"%s\"", newComment.PostID,
		)
	}

	comment := model.Comment{
		ID:              uuid.New(),
		CreatedAt:       time.Now(),
		PostID:          newComment.PostID,
		ParentCommentID: newComment.ParentCommentID,
		Text:            newComment.Text,
	}
	err = act.repo.CreateComment(ctx, comment)
	if err != nil {
		return model.Comment{}, fmt.Errorf("create comment: %w", err)
	}

	return comment, nil
}

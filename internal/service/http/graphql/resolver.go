package graphql

import (
	"context"
	"errors"

	"github.com/LoskutovDaniil/OzonTestTask2024/internal/action"
	"github.com/LoskutovDaniil/OzonTestTask2024/internal/adaptor"
	"github.com/LoskutovDaniil/OzonTestTask2024/internal/model"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type Resolver struct {
	mutation Mutation
	query    Query
}

type Mutation struct {
	repo adaptor.Repository
}

type Query struct {
	repo adaptor.Repository
}

func NewResolver(repo adaptor.Repository) Resolver {
	return Resolver{
		mutation: Mutation{repo},
		query:    Query{repo},
	}
}

//nolint:ireturn
func (r Resolver) Mutation() MutationResolver {
	return r.mutation
}

//nolint:ireturn
func (r Resolver) Query() QueryResolver {
	return r.query
}

func wrapError(ctx context.Context, e error) error {
	var knownErr model.Error
	if errors.As(e, &knownErr) {
		return &gqlerror.Error{
			Path:    graphql.GetPath(ctx),
			Message: e.Error(),
			Extensions: map[string]interface{}{
				"code": knownErr.Code(),
			},
		}
	}

	return &gqlerror.Error{
		Path:    graphql.GetPath(ctx),
		Message: "internal error: " + e.Error(),
		Extensions: map[string]interface{}{
			"code": 0,
		},
	}
}

func (srv Mutation) CreatePost(ctx context.Context, input model.NewPost) (*model.Post, error) {
	act := action.NewCreatePost(srv.repo)
	post, err := act.Do(ctx, input)
	if err != nil {
		return nil, wrapError(ctx, err)
	}

	return &post, nil
}

func (srv Mutation) CreateComment(ctx context.Context, input model.NewComment) (*model.Comment, error) {
	act := action.NewCreateComment(srv.repo)
	comment, err := act.Do(ctx, input)
	if err != nil {
		return nil, wrapError(ctx, err)
	}

	return &comment, nil
}

func (srv Query) Post(ctx context.Context, id uuid.UUID) (*model.Post, error) {
	act := action.NewGetPostByID(srv.repo)
	post, err := act.Do(ctx, id)
	if err != nil {
		return nil, wrapError(ctx, err)
	}

	return &post, nil
}

func (srv Query) Posts(ctx context.Context, input *model.PostsFilter) ([]model.Post, error) {
	var filter model.PostsFilter
	if input != nil {
		filter = *input
	}

	act := action.NewGetPostsList(srv.repo)
	posts, err := act.Do(ctx, filter)
	if err != nil {
		return nil, wrapError(ctx, err)
	}

	return posts, nil
}

func (srv Query) Comments(ctx context.Context, input *model.CommentsFilter) ([]model.Comment, error) {
	var filter model.CommentsFilter
	if input != nil {
		filter = *input
	}

	act := action.NewGetCommentsList(srv.repo)
	comments, err := act.Do(ctx, filter)
	if err != nil {
		return nil, wrapError(ctx, err)
	}

	return comments, nil
}

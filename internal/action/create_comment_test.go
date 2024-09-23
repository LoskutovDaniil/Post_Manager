package action_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/LoskutovDaniil/OzonTestTask2024/internal/action"
	"github.com/LoskutovDaniil/OzonTestTask2024/internal/action/mock"
	"github.com/LoskutovDaniil/OzonTestTask2024/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	tMock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestCreateComment(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(CreateCommentSuite))
}

type CreateCommentSuite struct {
	suite.Suite

	ctrl *gomock.Controller
	repo *mock.MockCreateCommentRepository
}

func (s *CreateCommentSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.repo = mock.NewMockCreateCommentRepository(s.ctrl)
}

func (s *CreateCommentSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *CreateCommentSuite) TestGetPostByIDError() {
	ctx := context.Background()
	newComment := model.NewComment{
		PostID: uuid.New(),
	}

	s.repo.EXPECT().
		GetPostByID(ctx, newComment.PostID).
		Return(model.Post{}, assert.AnError).
		Times(1)

	act := action.NewCreateComment(s.repo)
	comment, err := act.Do(ctx, newComment)

	s.Require().ErrorIs(err, assert.AnError)
	s.Require().Equal(model.Comment{}, comment)
}

func (s *CreateCommentSuite) TestPostCommentsDisabled() {
	ctx := context.Background()
	post := model.Post{
		ID:              uuid.New(),
		DisableComments: true,
	}
	newComment := model.NewComment{
		PostID: post.ID,
	}

	s.repo.EXPECT().
		GetPostByID(ctx, post.ID).
		Return(post, nil).
		Times(1)

	act := action.NewCreateComment(s.repo)
	comment, err := act.Do(ctx, newComment)

	s.Require().ErrorContains(err, fmt.Sprintf(
		"comments disabled for post with ID \"%s\"", post.ID,
	))
	s.Require().Equal(model.Comment{}, comment)
}

func (s *CreateCommentSuite) TestCreatePostError() {
	ctx := context.Background()
	post := model.Post{
		ID: uuid.New(),
	}
	newComment := model.NewComment{
		PostID: post.ID,
	}

	s.repo.EXPECT().
		GetPostByID(ctx, post.ID).
		Return(post, nil).
		Times(1)

	s.repo.EXPECT().
		CreateComment(ctx, gomock.Any()).
		Return(assert.AnError).
		Times(1)

	act := action.NewCreateComment(s.repo)
	comment, err := act.Do(ctx, newComment)

	s.Require().ErrorIs(err, assert.AnError)
	s.Require().Equal(model.Comment{}, comment)
}

func (s *CreateCommentSuite) TestSuccess() {
	ctx := context.Background()
	post := model.Post{
		ID: uuid.New(),
	}
	newComment := model.NewComment{
		PostID:          post.ID,
		ParentCommentID: uuid.New(),
		Text:            "My comment",
	}

	s.repo.EXPECT().
		GetPostByID(ctx, post.ID).
		Return(post, nil).
		Times(1)

	s.repo.EXPECT().
		CreateComment(ctx, tMock.MatchedBy(func(comment model.Comment) bool {
			return comment.PostID == newComment.PostID &&
				comment.ParentCommentID == newComment.ParentCommentID &&
				comment.Text == newComment.Text
		})).
		Return(nil).
		Times(1)

	act := action.NewCreateComment(s.repo)
	comment, err := act.Do(ctx, newComment)

	s.Require().NoError(err)

	s.Require().NotEqual(uuid.Nil, comment.ID)
	s.Require().NotEqual(time.Time{}, comment.CreatedAt)
	s.Require().Empty(comment.Comments)

	s.Require().Equal(newComment.PostID, comment.PostID)
	s.Require().Equal(newComment.ParentCommentID, comment.ParentCommentID)
	s.Require().Equal(newComment.Text, comment.Text)
}

package action_test

import (
	"context"
	"testing"

	"github.com/LoskutovDaniil/OzonTestTask2024/internal/action"
	"github.com/LoskutovDaniil/OzonTestTask2024/internal/action/mock"
	"github.com/LoskutovDaniil/OzonTestTask2024/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestGetCommentsList(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(GetCommentsListSuite))
}

type GetCommentsListSuite struct {
	suite.Suite

	ctrl *gomock.Controller
	repo *mock.MockGetCommentsListRepository
}

func (s *GetCommentsListSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.repo = mock.NewMockGetCommentsListRepository(s.ctrl)
}

func (s *GetCommentsListSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *GetCommentsListSuite) TestError() {
	ctx := context.Background()
	filter := model.CommentsFilter{
		PostID: uuid.New(),
	}.SetDefault()

	s.repo.EXPECT().
		GetCommentsList(ctx, filter).
		Return(nil, assert.AnError).
		Times(1)

	act := action.NewGetCommentsList(s.repo)
	comments, err := act.Do(ctx, filter)

	s.Require().ErrorIs(err, assert.AnError)
	s.Require().Empty(comments)
}

func (s *GetCommentsListSuite) TestSuccess() {
	ctx := context.Background()
	filter := model.CommentsFilter{
		PostID: uuid.New(),
	}.SetDefault()
	comments := []model.Comment{
		{ID: uuid.New(), Text: "My comment #1"},
		{ID: uuid.New(), Text: "My comment #2"},
	}

	s.repo.EXPECT().
		GetCommentsList(ctx, filter).
		Return(comments, nil).
		Times(1)

	act := action.NewGetCommentsList(s.repo)
	out, err := act.Do(ctx, filter)

	s.Require().NoError(err)
	s.Require().Equal(comments, out)
}

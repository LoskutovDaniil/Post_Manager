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

func TestGetPostsList(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(GetPostsListSuite))
}

type GetPostsListSuite struct {
	suite.Suite

	ctrl *gomock.Controller
	repo *mock.MockGetPostsListRepository
}

func (s *GetPostsListSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.repo = mock.NewMockGetPostsListRepository(s.ctrl)
}

func (s *GetPostsListSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *GetPostsListSuite) TestError() {
	ctx := context.Background()
	filter := model.PostsFilter{}.SetDefault()

	s.repo.EXPECT().
		GetPostsList(ctx, filter).
		Return(nil, assert.AnError).
		Times(1)

	act := action.NewGetPostsList(s.repo)
	posts, err := act.Do(ctx, filter)

	s.Require().ErrorIs(err, assert.AnError)
	s.Require().Empty(posts)
}

func (s *GetPostsListSuite) TestSuccess() {
	ctx := context.Background()
	filter := model.PostsFilter{}.SetDefault()
	posts := []model.Post{
		{ID: uuid.New(), Text: "My post #1"},
		{ID: uuid.New(), Text: "My post #2", DisableComments: true},
	}

	s.repo.EXPECT().
		GetPostsList(ctx, filter).
		Return(posts, nil).
		Times(1)

	act := action.NewGetPostsList(s.repo)
	out, err := act.Do(ctx, filter)

	s.Require().NoError(err)
	s.Require().Equal(posts, out)
}

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

func TestGetPostByID(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(GetPostByIDSuite))
}

type GetPostByIDSuite struct {
	suite.Suite

	ctrl *gomock.Controller
	repo *mock.MockGetPostByIDRepository
}

func (s *GetPostByIDSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.repo = mock.NewMockGetPostByIDRepository(s.ctrl)
}

func (s *GetPostByIDSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *GetPostByIDSuite) TestError() {
	ctx := context.Background()
	postID := uuid.New()

	s.repo.EXPECT().
		GetPostByID(ctx, postID).
		Return(model.Post{}, assert.AnError).
		Times(1)

	act := action.NewGetPostByID(s.repo)
	post, err := act.Do(ctx, postID)

	s.Require().ErrorIs(err, assert.AnError)
	s.Require().Equal(model.Post{}, post)
}

func (s *GetPostByIDSuite) TestSuccess() {
	ctx := context.Background()
	post := model.Post{
		ID:              uuid.New(),
		Text:            "My post",
		DisableComments: true,
	}

	s.repo.EXPECT().
		GetPostByID(ctx, post.ID).
		Return(post, nil).
		Times(1)

	act := action.NewGetPostByID(s.repo)
	out, err := act.Do(ctx, post.ID)

	s.Require().NoError(err)
	s.Require().Equal(post, out)
}

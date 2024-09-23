package action_test

import (
	"context"
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

func TestCreatePost(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(CreatePostSuite))
}

type CreatePostSuite struct {
	suite.Suite

	ctrl *gomock.Controller
	repo *mock.MockCreatePostRepository
}

func (s *CreatePostSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.repo = mock.NewMockCreatePostRepository(s.ctrl)
}

func (s *CreatePostSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *CreatePostSuite) TestError() {
	ctx := context.Background()

	s.repo.EXPECT().
		CreatePost(ctx, gomock.Any()).
		Return(assert.AnError).
		Times(1)

	act := action.NewCreatePost(s.repo)
	post, err := act.Do(ctx, model.NewPost{})

	s.Require().ErrorIs(err, assert.AnError)
	s.Require().Equal(model.Post{}, post)
}

func (s *CreatePostSuite) TestSuccess() {
	ctx := context.Background()
	newPost := model.NewPost{
		Text:            "My post text",
		DisableComments: true,
	}

	s.repo.EXPECT().
		CreatePost(ctx, tMock.MatchedBy(func(post model.Post) bool {
			return post.Text == newPost.Text &&
				post.DisableComments == newPost.DisableComments
		})).
		Return(nil).
		Times(1)

	act := action.NewCreatePost(s.repo)
	post, err := act.Do(ctx, newPost)

	s.Require().NoError(err)

	s.Require().NotEqual(uuid.Nil, post.ID)
	s.Require().NotEqual(time.Time{}, post.CreatedAt)

	s.Require().Equal(newPost.Text, post.Text)
	s.Require().Equal(newPost.DisableComments, post.DisableComments)
}

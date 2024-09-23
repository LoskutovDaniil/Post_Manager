package adaptor_test

import (
	"context"
	"testing"
	"time"

	"github.com/LoskutovDaniil/OzonTestTask2024/internal/adaptor"
	"github.com/LoskutovDaniil/OzonTestTask2024/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

func runSharedRepositoryTests(t *testing.T, repo adaptor.Repository) {
	t.Helper()
	suite.Run(t, &RepositorySharedTestsSuite{repo: repo})
}

type RepositorySharedTestsSuite struct {
	suite.Suite

	repo adaptor.Repository
}

func (*RepositorySharedTestsSuite) fixPostTime(expected, actual model.Post) model.Post {
	expected.CreatedAt = actual.CreatedAt

	return expected
}

func (*RepositorySharedTestsSuite) fixPostsTime(expected, actual []model.Post) []model.Post {
	if len(expected) != len(actual) {
		return expected
	}

	out := make([]model.Post, len(expected))
	for i, v := range expected {
		v.CreatedAt = actual[i].CreatedAt
		out[i] = v
	}

	return out
}

func (s *RepositorySharedTestsSuite) fixCommentsTime(expected, actual []model.Comment) []model.Comment {
	if len(expected) != len(actual) {
		return expected
	}

	out := make([]model.Comment, len(expected))
	for i, v := range expected {
		v.CreatedAt = actual[i].CreatedAt
		if len(v.Comments) != 0 {
			v.Comments = s.fixCommentsTime(v.Comments, actual[i].Comments)
		}
		out[i] = v
	}

	return out
}

func (s *RepositorySharedTestsSuite) TestCreatePost() {
	ctx := context.Background()
	post := model.Post{
		ID:              uuid.New(),
		CreatedAt:       time.Now(),
		Text:            "Create post test",
		DisableComments: true,
	}

	err := s.repo.CreatePost(ctx, post)
	s.Require().NoError(err)

	out, err := s.repo.GetPostByID(ctx, post.ID)
	s.Require().NoError(err)
	s.Equal(s.fixPostTime(post, out), out)
}

func (s *RepositorySharedTestsSuite) TestGetPostsList() {
	ctx := context.Background()

	time1 := time.Now()
	time.Sleep(1 * time.Second)
	time2 := time.Now()
	time.Sleep(1 * time.Second)
	time3 := time.Now()
	time.Sleep(1 * time.Second)

	posts := []model.Post{
		{
			ID:              uuid.New(),
			CreatedAt:       time3,
			Text:            "Get posts list #3",
			DisableComments: true,
		},
		{
			ID:        uuid.New(),
			CreatedAt: time2,
			Text:      "Get posts list #2",
		},
		{
			ID:              uuid.New(),
			CreatedAt:       time1,
			Text:            "Get posts list #1",
			DisableComments: true,
		},
	}
	filter := model.PostsFilter{
		Limit: 3,
	}

	for _, post := range posts {
		err := s.repo.CreatePost(ctx, post)
		s.Require().NoError(err)
	}

	out, err := s.repo.GetPostsList(ctx, filter)
	s.Require().NoError(err)
	s.Require().Equal(s.fixPostsTime(posts, out), out)
}

func (s *RepositorySharedTestsSuite) TestCreateComment() {
	ctx := context.Background()
	post := model.Post{
		ID:              uuid.New(),
		CreatedAt:       time.Now(),
		Text:            "Create comment test (post)",
		DisableComments: true,
	}
	comment := model.Comment{
		ID:        uuid.New(),
		PostID:    post.ID,
		CreatedAt: time.Now(),
		Text:      "Create comment test (comment)",
	}
	filter := model.CommentsFilter{
		PostID: post.ID,
		Depth:  1,
	}

	err := s.repo.CreatePost(ctx, post)
	s.Require().NoError(err)

	err = s.repo.CreateComment(ctx, comment)
	s.Require().NoError(err)

	out, err := s.repo.GetCommentsList(ctx, filter)
	s.Require().NoError(err)
	s.Equal(s.fixCommentsTime([]model.Comment{comment}, out), out)
}

func (s *RepositorySharedTestsSuite) TestAddCommentToNotExistsPost() {
	ctx := context.Background()
	comment := model.Comment{
		ID:        uuid.New(),
		PostID:    uuid.New(), // ID несуществующего поста
		CreatedAt: time.Now(),
		Text:      "Add comment to not exists post",
	}

	err := s.repo.CreateComment(ctx, comment)

	var knownErr model.Error
	s.Require().ErrorAs(err, &knownErr)
	s.Require().Equal(model.NotFoundCode, knownErr.Code())
}

//nolint:funlen
func (s *RepositorySharedTestsSuite) TestRecursiveComments() {
	ctx := context.Background()
	post := model.Post{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		Text:      "Recursive comments",
	}

	id1 := uuid.New()
	id2 := uuid.New()
	id3 := uuid.New()
	comments := []model.Comment{
		{
			ID:        uuid.New(),
			PostID:    post.ID,
			CreatedAt: time.Now(),
			Text:      "Recursive comments (comment #1)",
		},
		{
			ID:        id1,
			PostID:    post.ID,
			CreatedAt: time.Now(),
			Text:      "Recursive comments (comment #2)",
			Comments: []model.Comment{
				{
					ID:              id2,
					PostID:          post.ID,
					ParentCommentID: id1,
					CreatedAt:       time.Now(),
					Text:            "Recursive comments (comment #3)",
					Comments: []model.Comment{
						{
							ID:              id3,
							PostID:          post.ID,
							ParentCommentID: id2,
							CreatedAt:       time.Now(),
							Text:            "Recursive comments (comment #4)",
						},
					},
				},
			},
		},
	}

	err := s.repo.CreatePost(ctx, post)
	s.Require().NoError(err)

	var addCommentRecursive func([]model.Comment)
	addCommentRecursive = func(comments []model.Comment) {
		for _, comment := range comments {
			err = s.repo.CreateComment(ctx, comment)
			s.Require().NoError(err)

			addCommentRecursive(comment.Comments)
		}
	}
	addCommentRecursive(comments)

	// Добавляем комментарий который не будет отображаться в выборке
	// при Depth=3 и ParentCommentID=Nil, так как находится на 4 уровне.
	err = s.repo.CreateComment(ctx, model.Comment{
		ID:              uuid.New(),
		PostID:          post.ID,
		ParentCommentID: id3,
		CreatedAt:       time.Now(),
		Text:            "Recursive comments (comment #5)",
	})
	s.Require().NoError(err)

	out, err := s.repo.GetCommentsList(ctx, model.CommentsFilter{PostID: post.ID, Depth: 3})
	s.Require().NoError(err)
	s.Require().Equal(s.fixCommentsTime(comments, out), out)

	out, err = s.repo.GetCommentsList(ctx, model.CommentsFilter{PostID: post.ID, ParentCommentID: id1, Depth: 2})
	s.Require().NoError(err)
	s.Require().Equal(s.fixCommentsTime(comments[1].Comments, out), out)
}

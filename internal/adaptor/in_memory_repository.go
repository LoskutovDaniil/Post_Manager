package adaptor

import (
	"context"
	"sort"
	"sync"

	"github.com/LoskutovDaniil/OzonTestTask2024/internal/model"

	"github.com/google/uuid"
)

type inMemoryPost struct {
	model.Post

	Comments []model.Comment
}

type InMemoryRepository struct {
	posts map[uuid.UUID]inMemoryPost
	mu    sync.Mutex
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		posts: make(map[uuid.UUID]inMemoryPost),
	}
}

func (rep *InMemoryRepository) CreatePost(_ context.Context, post model.Post) error {
	rep.mu.Lock()
	defer rep.mu.Unlock()

	rep.posts[post.ID] = inMemoryPost{Post: post}

	return nil
}

func (rep *InMemoryRepository) GetPostsList(_ context.Context, filter model.PostsFilter) ([]model.Post, error) {
	rep.mu.Lock()
	defer rep.mu.Unlock()

	//nolint:prealloc
	var posts []model.Post
	for _, post := range rep.posts {
		posts = append(posts, post.Post)
	}
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})

	start := int(filter.Offset)
	end := int(filter.Offset + filter.Limit)

	if start > len(posts) {
		return []model.Post{}, nil
	}

	if end > len(posts) {
		end = len(posts)
	}

	return posts[start:end], nil
}

func (rep *InMemoryRepository) GetPostByID(_ context.Context, id uuid.UUID) (model.Post, error) {
	rep.mu.Lock()
	defer rep.mu.Unlock()

	post, exists := rep.posts[id]
	if !exists {
		return model.Post{}, model.NewNotFoundErrf("post with ID \"%s\" does not exist", id)
	}

	return post.Post, nil
}

func (rep *InMemoryRepository) addCommentRecursive(
	comments []model.Comment,
	targetID uuid.UUID,
	newComment model.Comment,
) ([]model.Comment, bool) {
	for i, comment := range comments {
		if comment.ID == targetID {
			comment.Comments = append(comment.Comments, newComment)
			comments[i] = comment

			return comments, true
		}

		var found bool
		comment.Comments, found = rep.addCommentRecursive(comment.Comments, targetID, newComment)
		if found {
			return comments, true
		}
	}

	return comments, false
}

func (rep *InMemoryRepository) CreateComment(_ context.Context, comment model.Comment) error {
	rep.mu.Lock()
	defer rep.mu.Unlock()

	post, exists := rep.posts[comment.PostID]
	if !exists {
		return model.NewNotFoundErrf("post with ID \"%s\" does not exist", comment.PostID)
	}

	comment.Comments = nil
	if comment.ParentCommentID == uuid.Nil {
		post.Comments = append(post.Comments, comment)
	} else {
		updatedComments, exists := rep.addCommentRecursive(post.Comments, comment.ParentCommentID, comment)
		if !exists {
			return model.NewNotFoundErrf("comment with ID \"%s\" does not exist", comment.ParentCommentID)
		}
		post.Comments = updatedComments
	}

	rep.posts[comment.PostID] = post

	return nil
}

func (rep *InMemoryRepository) getCommentByID(comments []model.Comment, id uuid.UUID) (model.Comment, bool) {
	for _, comment := range comments {
		if comment.ID == id {
			return comment, true
		}
		if len(comment.Comments) != 0 {
			out, ok := rep.getCommentByID(comment.Comments, id)
			if ok {
				return out, true
			}
		}
	}

	return model.Comment{}, false
}

func (rep *InMemoryRepository) getCommentsRecursive(comments []model.Comment, remainingDepth uint32) []model.Comment {
	if len(comments) == 0 || remainingDepth == 0 {
		return nil
	}

	out := make([]model.Comment, len(comments))
	for i, comment := range comments {
		comment.Comments = rep.getCommentsRecursive(comment.Comments, remainingDepth-1)
		out[i] = comment
	}

	return out
}

func (rep *InMemoryRepository) GetCommentsList(
	_ context.Context,
	filter model.CommentsFilter,
) ([]model.Comment, error) {
	rep.mu.Lock()
	defer rep.mu.Unlock()

	post, exists := rep.posts[filter.PostID]
	if !exists {
		return nil, model.NewNotFoundErrf("post with ID \"%s\" does not exist", filter.PostID)
	}

	comments := post.Comments
	if filter.ParentCommentID != uuid.Nil {
		comment, ok := rep.getCommentByID(comments, filter.ParentCommentID)
		if !ok {
			return nil, model.NewNotFoundErrf("comment with ID \"%s\" does not exist", filter.ParentCommentID)
		}
		comments = comment.Comments
	}
	comments = rep.getCommentsRecursive(comments, filter.Depth)

	return comments, nil
}

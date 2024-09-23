package adaptor

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LoskutovDaniil/OzonTestTask2024/internal/model"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db}
}

func (rep *PostgresRepository) CreatePost(ctx context.Context, post model.Post) error {
	_, err := rep.db.ExecContext(
		ctx,
		"INSERT INTO posts(id, created_at, post_text, disable_comments) VALUES ($1, $2, $3, $4);",
		post.ID,
		post.CreatedAt,
		post.Text,
		post.DisableComments,
	)
	if err != nil {
		return fmt.Errorf("create post: %w", err)
	}

	return nil
}

func (rep *PostgresRepository) GetPostsList(ctx context.Context, filter model.PostsFilter) ([]model.Post, error) {
	rows, err := rep.db.QueryContext(
		ctx,
		"SELECT id, created_at, post_text, disable_comments FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2;",
		filter.Limit,
		filter.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("get post: %w", err)
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		err = rows.Scan(&post.ID, &post.CreatedAt, &post.Text, &post.DisableComments)
		if err != nil {
			return nil, fmt.Errorf("get post: %w", err)
		}
		posts = append(posts, post)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("get post comments list: %w", err)
	}

	return posts, nil
}

func (rep *PostgresRepository) GetPostByID(ctx context.Context, id uuid.UUID) (model.Post, error) {
	row := rep.db.QueryRowContext(
		ctx,
		"SELECT id, created_at, post_text, disable_comments FROM posts WHERE id = $1;",
		id,
	)

	var post model.Post
	err := row.Scan(&post.ID, &post.CreatedAt, &post.Text, &post.DisableComments)
	if err != nil {
		return model.Post{}, fmt.Errorf("get post by ID: %w", err)
	}

	return post, nil
}

func (rep *PostgresRepository) CreateComment(ctx context.Context, comment model.Comment) error {
	var parentCommentID uuid.NullUUID
	if comment.ParentCommentID != uuid.Nil {
		parentCommentID.Valid = true
		parentCommentID.UUID = comment.ParentCommentID
	}

	_, err := rep.db.ExecContext(
		ctx,
		"INSERT INTO comments (id, post_id, parent_comment_id, created_at, comment_text) VALUES ($1, $2, $3, $4, $5);",
		comment.ID,
		comment.PostID,
		parentCommentID,
		comment.CreatedAt,
		comment.Text,
	)
	if err != nil {
		//nolint:errorlint
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == "23503" && pqErr.Constraint == "comments_post_id_fkey" {
			return model.NewNotFoundErrf("post with ID \"%s\" does not exist", comment.PostID)
		}

		return fmt.Errorf("add comment to post: %w", err)
	}

	return nil
}

func (rep *PostgresRepository) GetCommentsList(
	ctx context.Context,
	filter model.CommentsFilter,
) ([]model.Comment, error) {
	if filter.Depth == 0 {
		return nil, nil
	}

	rows, err := rep.db.QueryContext(
		ctx,
		"SELECT id, parent_comment_id, created_at, comment_text FROM comments WHERE post_id = $1 ORDER BY created_at;",
		filter.PostID,
	)
	if err != nil {
		return nil, fmt.Errorf("get post comments list: %w", err)
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		comment := model.Comment{PostID: filter.PostID}
		err = rows.Scan(&comment.ID, &comment.ParentCommentID, &comment.CreatedAt, &comment.Text)
		if err != nil {
			return nil, fmt.Errorf("get post comments list: %w", err)
		}
		comments = append(comments, comment)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("get post comments list: %w", err)
	}

	comments = rep.buildCommentsTree(comments, uuid.Nil)
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

func (rep *PostgresRepository) buildCommentsTree(comments []model.Comment, parentCommentID uuid.UUID) []model.Comment {
	if len(comments) == 0 {
		return nil
	}

	var out []model.Comment
	for _, comment := range comments {
		if comment.ParentCommentID == parentCommentID {
			out = append(out, comment)
		}
	}

	for i, comment := range out {
		out[i].Comments = rep.buildCommentsTree(comments, comment.ID)
	}

	return out
}

func (rep *PostgresRepository) getCommentByID(comments []model.Comment, id uuid.UUID) (model.Comment, bool) {
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

func (rep *PostgresRepository) getCommentsRecursive(comments []model.Comment, remainingDepth uint32) []model.Comment {
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

package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	DefaultFilterLimit   = 100
	DefaultFilterDepth   = 2
	CommentMaximumLength = 2000
)

type NewPost struct {
	Text            string `json:"text"`
	DisableComments bool   `json:"disableComments"`
}

type Post struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"createdAt"`       // Дата создания.
	Text            string    `json:"text"`            // Текст поста.
	DisableComments bool      `json:"disableComments"` // Возможность выключить комментарии.
}

type PostsFilter struct {
	Limit  uint32 `json:"limit,omitempty"`
	Offset uint32 `json:"offset,omitempty"`
}

func (p PostsFilter) SetDefault() PostsFilter {
	if p.Limit == 0 {
		p.Limit = DefaultFilterLimit
	}

	return p
}

type NewComment struct {
	PostID          uuid.UUID `json:"postId"`
	ParentCommentID uuid.UUID `json:"parentCommentId,omitempty"`
	Text            string    `json:"text"`
}

type Comment struct {
	ID              uuid.UUID `json:"id"`
	PostID          uuid.UUID `json:"postId"`
	ParentCommentID uuid.UUID `json:"parentCommentId,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	Text            string    `json:"text"`               // Текст комментария, максимум 2000 символов.
	Comments        []Comment `json:"comments,omitempty"` // Вложенные комментарии.
}

type CommentsFilter struct {
	PostID          uuid.UUID `json:"postId"`
	ParentCommentID uuid.UUID `json:"parentCommentId,omitempty"`
	Depth           uint32    `json:"depth,omitempty"` // Глубина разворачивания дерева комментариев.
}

func (p CommentsFilter) SetDefault() CommentsFilter {
	if p.Depth == 0 {
		p.Depth = DefaultFilterDepth
	}

	return p
}

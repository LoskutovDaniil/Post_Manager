-- +goose Up
CREATE TABLE posts (
    id               UUID      PRIMARY KEY,
    created_at       TIMESTAMP NOT NULL,
    post_text        TEXT      NOT NULL,
    disable_comments BOOLEAN   NOT NULL
);

CREATE TABLE comments (
    id                 UUID      PRIMARY KEY,
    post_id            UUID      NOT NULL    REFERENCES posts(id)    ON DELETE CASCADE,
    parent_comment_id  UUID                  REFERENCES comments(id) ON DELETE SET NULL,
    created_at         TIMESTAMP NOT NULL,
    comment_text       TEXT      NOT NULL
);

-- +goose Down
DROP TABLE posts;
DROP TABLE comments;

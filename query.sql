-- PostgreSQL queries with $1, $2 placeholders
-- name: CreateUser :exec
INSERT INTO users (name, email) VALUES ($1, $2);

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: GetUserPosts :many
SELECT posts.* FROM posts
JOIN users ON posts.user_id = users.id
WHERE users.id = $1;

-- name: UpdateUserEmail :exec
UPDATE users SET email = $1 WHERE id = $2;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: CreateComment :exec
INSERT INTO comments (post_id, user_id, content) 
VALUES ($1, $2, $3);

-- name: GetPostComments :many
SELECT * FROM comments 
WHERE post_id = $1 
ORDER BY created_at DESC;

-- name: GetUserComments :many
SELECT * FROM comments 
WHERE user_id = $1 
ORDER BY created_at DESC;

-- name: DeleteComment :exec
DELETE FROM comments WHERE id = $1;

-- name: GetUserWithPosts :many
SELECT u.*, p.* FROM users u
LEFT JOIN posts p ON u.id = p.user_id
WHERE u.id = $1;

-- name: GetRecentUsers :many
SELECT * FROM users 
WHERE created_at > NOW() - INTERVAL '7 days'
ORDER BY created_at DESC;

-- name: UpdateUserStatus :exec
UPDATE users SET status = $1 WHERE id = $2;

-- name: GetUserWithPosts :many
SELECT u.*, p.* FROM users u
LEFT JOIN posts p ON u.id = p.user_id
WHERE u.id = $1;

-- name: GetRecentUsers :many
SELECT * FROM users 
WHERE created_at > NOW() - INTERVAL '7 days'
ORDER BY created_at DESC;

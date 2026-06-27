-- +goose Up
-- +goose StatementBegin
INSERT INTO users (name, email, status) VALUES 
    ('Admin', 'admin@example.com', 'active'),
    ('Alice', 'alice@example.com', 'active'),
    ('Bob', 'bob@example.com', 'inactive');

INSERT INTO posts (user_id, title, content, published) VALUES 
    (1, 'Welcome Post', 'Welcome to our platform!', true),
    (2, 'Hello World', 'My first post.', true),
    (3, 'Hidden Post', 'This should be hidden.', false);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM posts WHERE user_id IN (1,2,3);
DELETE FROM users WHERE email IN ('admin@example.com', 'alice@example.com', 'bob@example.com');
-- +goose StatementEnd

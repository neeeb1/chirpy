-- +goose Up
ALTER TABLE users
ADD CONSTRAINT unique_user_id UNIQUE (id);

ALTER TABLE chirps
ADD CONSTRAINT unique_chirp_id UNIQUE (id);

ALTER TABLE refresh_tokens
ADD CONSTRAINT unique_token UNIQUE (token);


-- +goose Down
ALTER TABLE users
DROP CONSTRAINT unique_id;

ALTER TABLE chirps
DROP CONSTRAINT unique_id;

ALTER TABLE refresh_tokens
DROP CONSTRAINT unique_token;
-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(),
    time.Now(),
    time.Now(),
    $1
)
RETURNING *;
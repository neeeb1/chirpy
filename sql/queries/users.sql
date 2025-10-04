-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, hashed_password, email)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserbyID :one
SELECT * FROM users
where id = $1;

-- name: UpdateUserEmailAndPassword :one
UPDATE users
SET email = $1, hashed_password = $2, updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: UpgradeChirpyRed :one
UPDATE users
SET updated_at = NOW(), is_chirpy_red = TRUE
where id = $1
RETURNING *;

-- name: DowngradeChirpyRed :one
UPDATE users
SET updated_at = NOW(), is_chirpy_red = FALSE
where id = $1
RETURNING *;
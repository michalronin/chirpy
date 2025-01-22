-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, is_chirpy_red)
VALUES(
		gen_random_uuid(),
		NOW(),
		NOW(),
		$1,
		$2,
		FALSE
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByID :one
SELECT id, created_at, updated_at, email, is_chirpy_red FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET email = $2, hashed_password = $3, updated_at = NOW()
WHERE id = $1
RETURNING users.id, users.email, users.created_at, users.updated_at, users.is_chirpy_red;

-- name: UpgradeUserToChirpyRed :exec
UPDATE users
SET is_chirpy_red = TRUE, updated_at = NOW()
WHERE id = $1;

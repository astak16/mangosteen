-- name: CreateTag :one
INSERT INTO tags (user_id, name, sign, kind) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateTag :one
UPDATE tags SET user_id = $2, name = $3, sign = $4, kind = $5 WHERE id = $1 RETURNING *;
-- name: CreateItem :one
INSERT INTO items (user_id, amount, kind, happened_at, tag_ids ) VALUES ($1, $2, $3, $4, $5) RETURNING *;
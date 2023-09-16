-- name: CreateItem :one
INSERT INTO items (user_id, amount, kind, happened_at, tag_ids ) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: ListItem :many
SELECT * FROM items ORDER BY happened_at DESC OFFSET $1 LIMIT $2;

-- name: ListItemsHappenedBetween :many
SELECT * FROM items WHERE happened_at >= sqlc.arg(happened_after) AND happened_at < sqlc.arg(happened_before) ORDER BY happened_at DESC;

-- name: ListItemsByHappenedAtAndKind :many
SELECT * FROM items 
WHERE happened_at >= @happened_after 
AND happened_at < @happened_before 
AND kind = @kind 
AND user_id = @user_id
ORDER BY happened_at DESC;

-- name: CountItems :one
SELECT COUNT(*) FROM items;

-- name: DeleteAllItems :exec
DELETE FROM items;
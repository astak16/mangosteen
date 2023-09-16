-- name: CreateTag :one
INSERT INTO tags (user_id, name, sign, kind) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateTag :one
UPDATE tags 
SET 
  user_id = @user_id, 
  name = CASE WHEN @name::VARCHAR = '' THEN name ELSE @name END, 
  sign = CASE WHEN @sign::VARCHAR = '' THEN sign ELSE @sign END, 
  kind = CASE WHEN @kind::VARCHAR = '' THEN kind ELSE @kind END 
WHERE id = @id 
RETURNING *;

-- name: DeleteTag :exec
UPDATE tags SET deleted_at = NOW() WHERE id = $1;

-- name: FindTag :one
SELECT * FROM tags WHERE id = $1 AND deleted_at IS NULL;

-- name: ListTags :many
SELECT * FROM tags WHERE user_id = @user_id AND kind = @kind AND deleted_at IS NULL ORDER BY created_at DESC OFFSET $1 LIMIT $2;

-- name: GetTag :one
SELECT * FROM tags WHERE user_id = @user_id AND id = @id; 

-- name: CreateValidationCode :one
INSERT INTO validation_codes (email, code) VALUES ($1, $2) RETURNING *;

-- name: CountValidationCodes :one
SELECT COUNT(*) FROM validation_codes WHERE email = $1;

-- name: FindValidationCode :one
SELECT * FROM validation_codes WHERE email = $1 AND code = $2 AND used_at IS NULL ORDER BY created_at DESC LIMIT 1;
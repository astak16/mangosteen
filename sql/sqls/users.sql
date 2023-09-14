-- name: CreateUser :one
INSERT INTO USERS (EMAIL) VALUES ($1) RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM USERS;

-- name: FindUserByEmail :one
SELECT * FROM users WHERE email = $1;
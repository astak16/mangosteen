-- name: CreateUser :one
INSERT INTO USERS (EMAIL) VALUES ($1) RETURNING *;

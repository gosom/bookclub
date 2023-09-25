-- name: CreateUser :one
INSERT INTO users
    (email, passwd) 
VALUES 
    ($1, $2)
RETURNING *;

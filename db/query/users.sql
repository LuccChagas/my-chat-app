-- name: CreateUsers :one
INSERT INTO users
(id, "password", cpf, email, phone, name, first_name, last_name, nick_name, created_at)
VALUES( $1, $2, $3, $4, $5, $6, $7, $8, $9, now())
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
where users.id = $1;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUserByNickname :one
SELECT * FROM users
WHERE users.nick_name = $1;
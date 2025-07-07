-- name: SignUpUser :exec
INSERT INTO users (login, password)
VALUES ($1, $2);

-- name: GetUser :one
SELECT login, password
FROM users
WHERE login = $1;
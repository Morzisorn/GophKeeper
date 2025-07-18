-- name: SignUpUser :exec
INSERT INTO users (login, password, salt)
VALUES ($1, $2, $3);

-- name: GetUser :one
SELECT login, password, salt
FROM users
WHERE login = $1;

-- name: GetAllUserItems :many
SELECT 
    i.id,
    i.name,
    i.type,
    i.encrypted_data_content,
    i.encrypted_data_nonce,
    i.meta,
    i.created_at,
    i.updated_at
FROM items i
WHERE i.user_login = $1
ORDER BY i.created_at DESC;

-- name: GetUserItemsWithType :many
SELECT 
    i.id,
    i.name,
    i.type,
    i.encrypted_data_content,
    i.encrypted_data_nonce,
    i.meta,
    i.created_at,
    i.updated_at
FROM items i
WHERE i.user_login = $1 AND i.type = $2
ORDER BY created_at DESC;

-- name: GetTypesCounts :many
SELECT 
    type, 
    COUNT(*) as count
FROM items
WHERE user_login = $1
GROUP BY type;

-- name: AddItem :one
INSERT INTO items (user_login, name, type, encrypted_data_content, encrypted_data_nonce, meta)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: EditItem :exec
UPDATE items
SET name = $2, encrypted_data_content = $3, encrypted_data_nonce = $4, meta = $5, updated_at =  NOW()
WHERE id = $1;

-- name: DeleteItem :exec
DELETE FROM items
WHERE user_login = $1 AND id = $2;
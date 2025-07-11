-- name: SignUpUser :exec
INSERT INTO users (login, password)
VALUES ($1, $2);

-- name: GetUser :one
SELECT login, password
FROM users
WHERE login = $1;

-- name: GetAllUserItems :many
SELECT 
    i.id,
    i.name,
    i.type,
    i.meta,
    i.created_at,
    i.updated_at,
    -- Credentials
    cr.login,
    cr.password,
    -- Card data
    cd.number,
    cd.expiry_date,
    cd.security_code,
    cd.cardholder_name,
    -- Text data
    td.content as text_content,
    -- Binary data
    bd.content as binary_content
FROM items i
LEFT JOIN credentials_data cr ON i.id = cr.item_id AND i.type = 'CREDENTIALS'
LEFT JOIN card_data cd ON i.id = cd.item_id AND i.type = 'CARD'
LEFT JOIN text_data td ON i.id = td.item_id AND i.type = 'TEXT'
LEFT JOIN binary_data bd ON i.id = bd.item_id AND i.type = 'BINARY'
WHERE i.user_login = $1
ORDER BY i.created_at DESC;

-- name: GetCredentials :many
SELECT 
    i.id,
    i.name,
    i.meta,
    i.created_at,
    i.updated_at,
    cr.login,
    cr.password
FROM items i
LEFT JOIN credentials_data cr ON i.id = cr.item_id 
WHERE i.type = 'CREDENTIALS' AND i.user_login = $1
ORDER BY created_at DESC;

-- name: GetTexts :many
SELECT 
    i.id,
    i.name,
    i.meta,
    i.created_at,
    i.updated_at,
    t.content
FROM items i
LEFT JOIN text_data t ON i.id = t.item_id 
WHERE i.type = 'TEXT' AND i.user_login = $1
ORDER BY created_at DESC;

-- name: GetBinaries :many
SELECT 
    i.id,
    i.name,
    i.meta,
    i.created_at,
    i.updated_at,
    b.content
FROM items i
LEFT JOIN binary_data b ON i.id = b.item_id 
WHERE i.type = 'BINARY' AND i.user_login = $1
ORDER BY created_at DESC;

-- name: GetCards :many
SELECT 
    i.id,
    i.name,
    i.meta,
    i.created_at,
    i.updated_at,
    c.number,
    c.expiry_date,
    c.security_code,
    c.cardholder_name
FROM items i
LEFT JOIN card_data c ON i.id = c.item_id 
WHERE i.type = 'BINARY' AND i.user_login = $1
ORDER BY created_at DESC;

-- name: GetTypesCounts :many
SELECT 
    type, 
    COUNT(*) as count
FROM items
WHERE user_login = $1
GROUP BY type;

-- name: AddItem :one
INSERT INTO items (user_login, name, type, meta)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: AddCredentials :exec
INSERT INTO credentials_data (item_id, login, password)
VALUES ($1, $2, $3);

-- name: AddText :exec
INSERT INTO text_data (item_id, content)
VALUES ($1, $2);

-- name: AddBinary :exec
INSERT INTO binary_data (item_id, content)
VALUES ($1, $2);

-- name: AddCard :exec
INSERT INTO card_data (item_id, number, expiry_date, security_code, cardholder_name)
VALUES ($1, $2, $3, $4, $5);

-- name: EditItem :exec
UPDATE items
SET name = $2, meta = $3, updated_at =  NOW()
WHERE id = $1;

-- name: EditCredentials :exec
UPDATE credentials_data
SET login = $2, password = $3 
WHERE item_id = $1;

-- name: EditText :exec
UPDATE text_data
SET content = $2 
WHERE item_id = $1;

-- name: EditBinary :exec
UPDATE binary_data
SET content = $2
WHERE item_id = $1;

-- name: EditCard :exec
UPDATE card_data
SET number = $2, expiry_date = $3, security_code = $4, cardholder_name = $5 
WHERE item_id = $1;

-- name: DeleteItem :exec
DELETE FROM items
WHERE user_login = $1 AND id = $2;
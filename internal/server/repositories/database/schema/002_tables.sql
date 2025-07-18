CREATE TABLE IF NOT EXISTS users (
    login VARCHAR(50) NOT NULL PRIMARY KEY,
    password BYTEA NOT NULL,
    salt  TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_login VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    encrypted_data_content TEXT NOT NULL,
    encrypted_data_nonce VARCHAR(50) NOT NULL,
    type item_type NOT NULL,
    meta JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_login) REFERENCES users(login)
);

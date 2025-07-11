CREATE TABLE IF NOT EXISTS users (
    login VARCHAR(50) NOT NULL PRIMARY KEY,
    password BYTEA NOT NULL
);

CREATE TYPE item_type AS ENUM ('CREDENTIALS', 'TEXT', 'BINARY', 'CARD');

CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_login VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    type item_type NOT NULL,
    meta JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_login) REFERENCES users(login)
);

CREATE TABLE IF NOT EXISTS credentials_data (
    item_id UUID PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,
    login VARCHAR(255) NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE text_data (
    item_id UUID PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,
    content TEXT NOT NULL 
);

CREATE TABLE binary_data (
    item_id UUID PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,
    content BYTEA NOT NULL 
);

CREATE TABLE card_data (
    item_id UUID PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,
    number VARCHAR(255) NOT NULL, 
    expiry_date VARCHAR(10) NOT NULL,
    security_code VARCHAR(10) NOT NULL,
    cardholder_name VARCHAR(255) NOT NULL 
);
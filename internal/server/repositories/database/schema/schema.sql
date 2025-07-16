CREATE TABLE IF NOT EXISTS users (
    login VARCHAR(50) NOT NULL PRIMARY KEY,
    password BYTEA NOT NULL,
    salt  TEXT NOT NULL
);

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'item_type') THEN
        CREATE TYPE item_type AS ENUM ('CREDENTIALS', 'TEXT', 'BINARY', 'CARD');
    END IF;
END $$;

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

CREATE TABLE IF NOT EXISTS text_data (
    item_id UUID PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,
    content TEXT NOT NULL 
);

CREATE TABLE IF NOT EXISTS binary_data (
    item_id UUID PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,
    content BYTEA NOT NULL 
);

CREATE TABLE IF NOT EXISTS card_data (
    item_id UUID PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,
    number VARCHAR(255) NOT NULL, 
    expiry_date VARCHAR(10) NOT NULL,
    security_code VARCHAR(10) NOT NULL,
    cardholder_name VARCHAR(255) NOT NULL 
);
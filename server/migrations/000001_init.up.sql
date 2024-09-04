DO $$ BEGIN
    CREATE TYPE data_type AS ENUM ('login_password', 'text', 'binary', 'bank_card');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

CREATE TABLE IF NOT EXISTS data_items (
                            id TEXT NOT NULL PRIMARY KEY,
                            user_id INT REFERENCES users(id) ON DELETE CASCADE,
                            type data_type NOT NULL,
                            data BYTEA NOT NULL,
                            meta TEXT DEFAULT '',
                            url VARCHAR(255) NOT NULL DEFAULT '',
                            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_data_items_user_id ON data_items(user_id);
CREATE INDEX IF NOT EXISTS idx_data_items_type ON data_items(type);

CREATE TYPE data_type AS ENUM ('login_password', 'text', 'binary', 'bank_card');

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);

CREATE TABLE data_items (
                            id SERIAL PRIMARY KEY,
                            user_id INT REFERENCES users(id) ON DELETE CASCADE,
                            type data_type NOT NULL,  -- Использование ENUM типа данных
                            data BYTEA NOT NULL,  -- Данные хранятся в бинарном формате
                            meta TEXT,
                            url VARCHAR(255) NOT NULL DEFAULT '',
                            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_data_items_user_id ON data_items(user_id);
CREATE INDEX idx_data_items_type ON data_items(type);

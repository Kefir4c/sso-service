CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    pass_hash BYTEA NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email ON users(email);

COMMENT ON TABLE users IS 'Пользователи сервиса аутентификации';
COMMENT ON COLUMN users.id IS 'Уникальный идентификатор пользователя';
COMMENT ON COLUMN users.email IS 'Email пользователя (используется как логин)';
COMMENT ON COLUMN users.pass_hash IS 'Bcrypt хеш пароля';
COMMENT ON COLUMN users.is_admin IS 'Флаг администратора';
COMMENT ON COLUMN users.created_at IS 'Время создания записи';
COMMENT ON COLUMN users.updated_at IS 'Время последнего обновления';

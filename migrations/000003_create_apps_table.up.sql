CREATE TABLE IF NOT EXISTS apps (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    secret VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE apps IS 'Приложения, использующие SSO';
COMMENT ON COLUMN apps.id IS 'ID приложения';
COMMENT ON COLUMN apps.name IS 'Название приложения';
COMMENT ON COLUMN apps.secret IS 'Секретный ключ для подписи JWT';
COMMENT ON COLUMN apps.created_at IS 'Время регистрации приложения';
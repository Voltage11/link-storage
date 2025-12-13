-- ==================== TABLE: registrations ====================
CREATE TABLE registrations(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hashed VARCHAR(255) NOT NULL,
    ip_address VARCHAR(50) NOT NULL,
    user_agent TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    expired_at TIMESTAMPTZ NOT NULL,
    activated_at TIMESTAMPTZ,
    token VARCHAR(255) NOT NULL UNIQUE,
    verify_code varchar(5)
);
COMMENT ON TABLE registrations IS 'Заявки на регистрацию';
CREATE INDEX idx_registrations_email ON registrations (email);
CREATE INDEX idx_registrations_token ON registrations (token);

-- ==================== TABLE: users ====================
CREATE TABLE users
(
    id             SERIAL PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    email          VARCHAR(255) NOT NULL UNIQUE,
    password_hashed  VARCHAR(255) NOT NULL,
    is_active      BOOLEAN      NOT NULL DEFAULT false,
    is_admin       BOOLEAN      NOT NULL DEFAULT false,
    last_login_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE users IS 'Пользователи системы';

CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_is_active ON users (is_active) WHERE is_active = true;

-- ==================== TABLE: sessions ====================
CREATE TABLE sessions
(
    id            SERIAL PRIMARY KEY,
    user_id       INTEGER             NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    refresh_token VARCHAR(1000) UNIQUE NOT NULL,
    user_agent    TEXT DEFAULT '',
    ip_address    VARCHAR(50),
    expired_at    TIMESTAMPTZ         NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE sessions IS 'Сессии пользователей с refresh tokens';

CREATE INDEX idx_sessions_refresh_token ON sessions (refresh_token);
CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE INDEX idx_sessions_expired_at ON sessions (expired_at);
CREATE INDEX idx_sessions_user_id_expired_at ON sessions (user_id, expired_at);

-- ===================== TABLE: link_groups ===================
CREATE TABLE link_groups(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    description TEXT DEFAULT '',
    position INTEGER DEFAULT 0,
    color VARCHAR(7) DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE link_groups IS 'Группы ссылок';
CREATE INDEX idx_link_groups_user_id ON link_groups(user_id);

-- ===================== TABLE: links ===================
CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    link_group_id INTEGER REFERENCES link_groups(id) ON DELETE SET NULL,
    url TEXT NOT NULL,
    title VARCHAR(500) DEFAULT '',
    description TEXT DEFAULT '',
    favicon_url TEXT DEFAULT '',
    preview_image TEXT DEFAULT '',
    is_archived BOOLEAN DEFAULT FALSE,
    is_favorite BOOLEAN DEFAULT FALSE,
    click_count INTEGER DEFAULT 0,
    last_visited TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE links IS 'Ссылки';
CREATE INDEX idx_links_user_id ON links(user_id);
CREATE INDEX idx_links_link_group_id ON links(link_group_id);

-- ===================== TABLE: tags ===================
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    color VARCHAR(7) DEFAULT '',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE tags IS 'Теги';
CREATE INDEX idx_tags_user_id ON tags(user_id);
CREATE INDEX idx_tags_name ON tags(name);

-- ===================== TABLE: link_tags ===================
CREATE TABLE link_tags (
    link_id INTEGER NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE
);
COMMENT ON TABLE link_tags IS 'Теги для ссылок';
CREATE INDEX idx_link_tags_link_id ON link_tags(link_id);
CREATE INDEX idx_link_tags_tag_id ON link_tags(tag_id);
CREATE UNIQUE INDEX idx_link_tags_link_id_tag_id ON link_tags(link_id, tag_id);

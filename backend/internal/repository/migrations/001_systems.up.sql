CREATE TABLE IF NOT EXISTS systems (
    id UUID PRIMARY KEY,
    name VARCHAR(120) NOT NULL,
    slug VARCHAR(120) NOT NULL UNIQUE,
    edition VARCHAR(255),
    publisher VARCHAR(255),
    description TEXT,
    license_type VARCHAR(32),
    version VARCHAR(32) NOT NULL DEFAULT '0.1.0',
    playstyle VARCHAR(32),
    complexity INT CHECK (complexity >= 1 AND complexity <= 5),
    measurement_unit VARCHAR(16),
    currency_symbol VARCHAR(16),
    status VARCHAR(16) NOT NULL DEFAULT 'draft',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    system_family VARCHAR(255),
    player_count_min INT,
    player_count_max INT,
    official_links JSONB NOT NULL DEFAULT '[]',
    tags JSONB NOT NULL DEFAULT '[]',
    core_rulebooks JSONB NOT NULL DEFAULT '[]',
    icon_url VARCHAR(512),
    cover_url VARCHAR(512),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_systems_status ON systems (status);
CREATE INDEX idx_systems_is_active ON systems (is_active);
CREATE INDEX idx_systems_slug ON systems (slug);

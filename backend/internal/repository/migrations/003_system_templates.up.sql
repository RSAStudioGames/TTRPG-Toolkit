CREATE TABLE IF NOT EXISTS system_templates (
    id UUID PRIMARY KEY,
    name VARCHAR(120) NOT NULL,
    description TEXT,
    schema_json JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS system_tags (
    system_id UUID NOT NULL REFERENCES systems (id) ON DELETE CASCADE,
    tag VARCHAR(64) NOT NULL,
    PRIMARY KEY (system_id, tag)
);

CREATE TABLE IF NOT EXISTS system_rulebooks (
    system_id UUID NOT NULL REFERENCES systems (id) ON DELETE CASCADE,
    rulebook VARCHAR(255) NOT NULL,
    PRIMARY KEY (system_id, rulebook)
);

CREATE TABLE IF NOT EXISTS system_links (
    system_id UUID NOT NULL REFERENCES systems (id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    PRIMARY KEY (system_id, url)
);

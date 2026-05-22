CREATE TABLE IF NOT EXISTS system_mechanics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    system_id UUID NOT NULL REFERENCES systems (id) ON DELETE CASCADE,
    resolution_config JSONB NOT NULL DEFAULT '{}',
    progression_config JSONB NOT NULL DEFAULT '{}',
    action_economy_config JSONB NOT NULL DEFAULT '{}',
    UNIQUE (system_id)
);

CREATE INDEX IF NOT EXISTS idx_system_mechanics_system_id ON system_mechanics (system_id);

CREATE TABLE IF NOT EXISTS system_attributes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    system_id UUID NOT NULL REFERENCES systems (id) ON DELETE CASCADE,
    group_name VARCHAR(255),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(32) NOT NULL,
    config JSONB NOT NULL DEFAULT '{}',
    sort_order INT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_system_attributes_system_id ON system_attributes (system_id);
CREATE INDEX IF NOT EXISTS idx_system_attributes_system_sort ON system_attributes (system_id, sort_order);

CREATE TABLE IF NOT EXISTS system_skills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    system_id UUID NOT NULL REFERENCES systems (id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    linked_attribute_id UUID REFERENCES system_attributes (id) ON DELETE SET NULL,
    type VARCHAR(32) NOT NULL,
    category VARCHAR(255),
    config JSONB NOT NULL DEFAULT '{}',
    sort_order INT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_system_skills_system_id ON system_skills (system_id);
CREATE INDEX IF NOT EXISTS idx_system_skills_linked_attribute ON system_skills (linked_attribute_id);
CREATE INDEX IF NOT EXISTS idx_system_skills_system_sort ON system_skills (system_id, sort_order);

CREATE TABLE IF NOT EXISTS system_resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    system_id UUID NOT NULL REFERENCES systems (id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(32) NOT NULL,
    config JSONB NOT NULL DEFAULT '{}',
    sort_order INT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_system_resources_system_id ON system_resources (system_id);
CREATE INDEX IF NOT EXISTS idx_system_resources_system_sort ON system_resources (system_id, sort_order);

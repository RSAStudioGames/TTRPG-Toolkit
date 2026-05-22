ALTER TABLE system_mechanics
  ADD COLUMN IF NOT EXISTS attributes_config JSONB NOT NULL DEFAULT '{}';

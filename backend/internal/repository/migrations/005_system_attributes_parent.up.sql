ALTER TABLE system_attributes
  ADD COLUMN IF NOT EXISTS parent_attribute_id UUID
  REFERENCES system_attributes (id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_system_attributes_parent
  ON system_attributes (parent_attribute_id);

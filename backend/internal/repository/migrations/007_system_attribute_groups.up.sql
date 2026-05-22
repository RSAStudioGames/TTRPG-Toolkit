CREATE TABLE system_attribute_groups (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  system_id UUID NOT NULL REFERENCES systems (id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  sort_order INT NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX idx_system_attribute_groups_system_name
  ON system_attribute_groups (system_id, lower(name));

ALTER TABLE system_attributes
  ADD COLUMN attribute_group_id UUID
  REFERENCES system_attribute_groups (id) ON DELETE SET NULL;

CREATE INDEX idx_system_attributes_group ON system_attributes (attribute_group_id);

INSERT INTO system_attribute_groups (id, system_id, name, sort_order)
SELECT gen_random_uuid(), system_id, trim(group_name),
       ROW_NUMBER() OVER (PARTITION BY system_id ORDER BY trim(group_name)) - 1
FROM system_attributes
WHERE group_name IS NOT NULL AND trim(group_name) <> ''
GROUP BY system_id, trim(group_name);

UPDATE system_attributes a
SET attribute_group_id = g.id,
    group_name = g.name
FROM system_attribute_groups g
WHERE a.system_id = g.system_id
  AND a.group_name IS NOT NULL AND trim(a.group_name) <> ''
  AND lower(trim(a.group_name)) = lower(g.name);

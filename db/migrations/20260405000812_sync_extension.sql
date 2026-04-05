ALTER TABLE armor_pieces
    ADD COLUMN IF NOT EXISTS armor_set_name text NULL;
CREATE INDEX IF NOT EXISTS idx_armor_pieces_armor_set_name ON armor_pieces (armor_set_name);

ALTER TABLE kinsects
    ADD COLUMN IF NOT EXISTS kinsect_bonus_primary text NULL,
    ADD COLUMN IF NOT EXISTS kinsect_bonus_secondary text NULL;
CREATE INDEX IF NOT EXISTS idx_kinsects_kinsect_bonus_primary ON kinsects (kinsect_bonus_primary);
CREATE INDEX IF NOT EXISTS idx_kinsects_kinsect_bonus_secondary ON kinsects (kinsect_bonus_secondary);

CREATE TABLE IF NOT EXISTS skill_levels (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    skill_id uuid NOT NULL,
    level smallint NOT NULL DEFAULT 1,
    effect_value_text text NULL,
    description text NULL,
    details_json jsonb NOT NULL DEFAULT '{}'::jsonb,
    CONSTRAINT uq_skill_levels_skill_level UNIQUE (skill_id, level),
    CONSTRAINT fk_skill_levels_skill_id FOREIGN KEY (skill_id) REFERENCES skills (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_skill_levels_skill_id ON skill_levels (skill_id);

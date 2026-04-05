CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS languages (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    code varchar(16) NOT NULL,
    label varchar(64) NOT NULL,
    is_active boolean NOT NULL DEFAULT true,
    sort_order smallint NOT NULL DEFAULT 0,
    CONSTRAINT uq_languages_code UNIQUE (code)
);
CREATE INDEX IF NOT EXISTS idx_languages_is_active ON languages (is_active);

CREATE TABLE IF NOT EXISTS source_sync_runs (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    provider text NOT NULL,
    category text NOT NULL,
    status text NOT NULL,
    triggered_by text NOT NULL DEFAULT 'manual',
    started_at timestamptz NOT NULL DEFAULT now(),
    finished_at timestamptz NULL,
    error_message text NULL,
    metadata_json jsonb NOT NULL DEFAULT '{}'::jsonb
);
CREATE INDEX IF NOT EXISTS idx_source_sync_runs_provider ON source_sync_runs (provider);
CREATE INDEX IF NOT EXISTS idx_source_sync_runs_category ON source_sync_runs (category);
CREATE INDEX IF NOT EXISTS idx_source_sync_runs_status ON source_sync_runs (status);

CREATE TABLE IF NOT EXISTS skills (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz NULL,
    provider text NOT NULL,
    source_id text NOT NULL,
    source_url text NULL,
    source_slug text NULL,
    fetched_at timestamptz NOT NULL DEFAULT now(),
    sync_run_id uuid NULL,
    skill_kind text NOT NULL,
    max_level smallint NOT NULL DEFAULT 1,
    is_binary boolean NOT NULL DEFAULT false,
    is_set_bonus_skill boolean NOT NULL DEFAULT false,
    sort_order integer NOT NULL DEFAULT 0,
    details_json jsonb NOT NULL DEFAULT '{}'::jsonb,
    CONSTRAINT uq_skills_provider_source_id UNIQUE (provider, source_id),
    CONSTRAINT fk_skills_sync_run_id FOREIGN KEY (sync_run_id) REFERENCES source_sync_runs (id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_skills_deleted_at ON skills (deleted_at);
CREATE INDEX IF NOT EXISTS idx_skills_source_slug ON skills (source_slug);
CREATE INDEX IF NOT EXISTS idx_skills_sync_run_id ON skills (sync_run_id);
CREATE INDEX IF NOT EXISTS idx_skills_skill_kind ON skills (skill_kind);
CREATE INDEX IF NOT EXISTS idx_skills_sort_order ON skills (sort_order);

CREATE TABLE IF NOT EXISTS items (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz NULL,
    provider text NOT NULL,
    source_id text NOT NULL,
    source_url text NULL,
    source_slug text NULL,
    fetched_at timestamptz NOT NULL DEFAULT now(),
    sync_run_id uuid NULL,
    category text NOT NULL,
    subcategory text NULL,
    rarity smallint NOT NULL DEFAULT 0,
    carry_limit smallint NULL,
    buy_price integer NULL,
    sell_price integer NULL,
    points integer NULL,
    is_crafting_material boolean NOT NULL DEFAULT false,
    is_consumable boolean NOT NULL DEFAULT false,
    icon_key text NULL,
    sort_order integer NOT NULL DEFAULT 0,
    details_json jsonb NOT NULL DEFAULT '{}'::jsonb,
    CONSTRAINT uq_items_provider_source_id UNIQUE (provider, source_id),
    CONSTRAINT fk_items_sync_run_id FOREIGN KEY (sync_run_id) REFERENCES source_sync_runs (id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_items_deleted_at ON items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_items_source_slug ON items (source_slug);
CREATE INDEX IF NOT EXISTS idx_items_sync_run_id ON items (sync_run_id);
CREATE INDEX IF NOT EXISTS idx_items_category ON items (category);
CREATE INDEX IF NOT EXISTS idx_items_subcategory ON items (subcategory);
CREATE INDEX IF NOT EXISTS idx_items_rarity ON items (rarity);
CREATE INDEX IF NOT EXISTS idx_items_sort_order ON items (sort_order);

CREATE TABLE IF NOT EXISTS weapons (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz NULL,
    provider text NOT NULL,
    source_id text NOT NULL,
    source_url text NULL,
    source_slug text NULL,
    fetched_at timestamptz NOT NULL DEFAULT now(),
    sync_run_id uuid NULL,
    weapon_type text NOT NULL,
    rarity smallint NOT NULL DEFAULT 0,
    attack integer NOT NULL DEFAULT 0,
    affinity_percent smallint NOT NULL DEFAULT 0,
    defense_bonus smallint NOT NULL DEFAULT 0,
    slot1_level smallint NOT NULL DEFAULT 0,
    slot2_level smallint NOT NULL DEFAULT 0,
    slot3_level smallint NOT NULL DEFAULT 0,
    element_type text NULL,
    element_value integer NOT NULL DEFAULT 0,
    ailment_type text NULL,
    ailment_value integer NOT NULL DEFAULT 0,
    sharpness_red smallint NOT NULL DEFAULT 0,
    sharpness_orange smallint NOT NULL DEFAULT 0,
    sharpness_yellow smallint NOT NULL DEFAULT 0,
    sharpness_green smallint NOT NULL DEFAULT 0,
    sharpness_blue smallint NOT NULL DEFAULT 0,
    sharpness_white smallint NOT NULL DEFAULT 0,
    sharpness_purple smallint NOT NULL DEFAULT 0,
    craft_cost integer NULL,
    upgrade_cost integer NULL,
    tree_depth smallint NOT NULL DEFAULT 0,
    is_final_upgrade boolean NOT NULL DEFAULT false,
    sort_order integer NOT NULL DEFAULT 0,
    details_json jsonb NOT NULL DEFAULT '{}'::jsonb,
    CONSTRAINT uq_weapons_provider_source_id UNIQUE (provider, source_id),
    CONSTRAINT fk_weapons_sync_run_id FOREIGN KEY (sync_run_id) REFERENCES source_sync_runs (id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_weapons_deleted_at ON weapons (deleted_at);
CREATE INDEX IF NOT EXISTS idx_weapons_source_slug ON weapons (source_slug);
CREATE INDEX IF NOT EXISTS idx_weapons_sync_run_id ON weapons (sync_run_id);
CREATE INDEX IF NOT EXISTS idx_weapons_weapon_type ON weapons (weapon_type);
CREATE INDEX IF NOT EXISTS idx_weapons_rarity ON weapons (rarity);
CREATE INDEX IF NOT EXISTS idx_weapons_element_type ON weapons (element_type);
CREATE INDEX IF NOT EXISTS idx_weapons_ailment_type ON weapons (ailment_type);
CREATE INDEX IF NOT EXISTS idx_weapons_sort_order ON weapons (sort_order);

CREATE TABLE IF NOT EXISTS armor_pieces (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz NULL,
    provider text NOT NULL,
    source_id text NOT NULL,
    source_url text NULL,
    source_slug text NULL,
    fetched_at timestamptz NOT NULL DEFAULT now(),
    sync_run_id uuid NULL,
    armor_set_key text NOT NULL,
    armor_set_variant text NULL,
    piece_type text NOT NULL,
    rank text NOT NULL,
    gender text NULL,
    rarity smallint NOT NULL DEFAULT 0,
    defense_base smallint NOT NULL DEFAULT 0,
    defense_max smallint NOT NULL DEFAULT 0,
    defense_augmented_max smallint NOT NULL DEFAULT 0,
    fire_res smallint NOT NULL DEFAULT 0,
    water_res smallint NOT NULL DEFAULT 0,
    thunder_res smallint NOT NULL DEFAULT 0,
    ice_res smallint NOT NULL DEFAULT 0,
    dragon_res smallint NOT NULL DEFAULT 0,
    slot1_level smallint NOT NULL DEFAULT 0,
    slot2_level smallint NOT NULL DEFAULT 0,
    slot3_level smallint NOT NULL DEFAULT 0,
    is_layered boolean NOT NULL DEFAULT false,
    sort_order integer NOT NULL DEFAULT 0,
    details_json jsonb NOT NULL DEFAULT '{}'::jsonb,
    CONSTRAINT uq_armor_pieces_provider_source_id UNIQUE (provider, source_id),
    CONSTRAINT fk_armor_pieces_sync_run_id FOREIGN KEY (sync_run_id) REFERENCES source_sync_runs (id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_armor_pieces_deleted_at ON armor_pieces (deleted_at);
CREATE INDEX IF NOT EXISTS idx_armor_pieces_source_slug ON armor_pieces (source_slug);
CREATE INDEX IF NOT EXISTS idx_armor_pieces_sync_run_id ON armor_pieces (sync_run_id);
CREATE INDEX IF NOT EXISTS idx_armor_pieces_armor_set_key ON armor_pieces (armor_set_key);
CREATE INDEX IF NOT EXISTS idx_armor_pieces_piece_type ON armor_pieces (piece_type);
CREATE INDEX IF NOT EXISTS idx_armor_pieces_rank ON armor_pieces (rank);
CREATE INDEX IF NOT EXISTS idx_armor_pieces_gender ON armor_pieces (gender);
CREATE INDEX IF NOT EXISTS idx_armor_pieces_rarity ON armor_pieces (rarity);
CREATE INDEX IF NOT EXISTS idx_armor_pieces_sort_order ON armor_pieces (sort_order);

CREATE TABLE IF NOT EXISTS decorations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz NULL,
    provider text NOT NULL,
    source_id text NOT NULL,
    source_url text NULL,
    source_slug text NULL,
    fetched_at timestamptz NOT NULL DEFAULT now(),
    sync_run_id uuid NULL,
    slot_size smallint NOT NULL DEFAULT 1,
    rarity smallint NOT NULL DEFAULT 0,
    is_craftable boolean NOT NULL DEFAULT true,
    craft_cost integer NULL,
    sort_order integer NOT NULL DEFAULT 0,
    details_json jsonb NOT NULL DEFAULT '{}'::jsonb,
    CONSTRAINT uq_decorations_provider_source_id UNIQUE (provider, source_id),
    CONSTRAINT fk_decorations_sync_run_id FOREIGN KEY (sync_run_id) REFERENCES source_sync_runs (id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_decorations_deleted_at ON decorations (deleted_at);
CREATE INDEX IF NOT EXISTS idx_decorations_source_slug ON decorations (source_slug);
CREATE INDEX IF NOT EXISTS idx_decorations_sync_run_id ON decorations (sync_run_id);
CREATE INDEX IF NOT EXISTS idx_decorations_slot_size ON decorations (slot_size);
CREATE INDEX IF NOT EXISTS idx_decorations_rarity ON decorations (rarity);
CREATE INDEX IF NOT EXISTS idx_decorations_sort_order ON decorations (sort_order);

CREATE TABLE IF NOT EXISTS charms (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz NULL,
    provider text NOT NULL,
    source_id text NOT NULL,
    source_url text NULL,
    source_slug text NULL,
    fetched_at timestamptz NOT NULL DEFAULT now(),
    sync_run_id uuid NULL,
    rarity smallint NOT NULL DEFAULT 0,
    max_rank smallint NOT NULL DEFAULT 1,
    craft_cost integer NULL,
    sort_order integer NOT NULL DEFAULT 0,
    details_json jsonb NOT NULL DEFAULT '{}'::jsonb,
    CONSTRAINT uq_charms_provider_source_id UNIQUE (provider, source_id),
    CONSTRAINT fk_charms_sync_run_id FOREIGN KEY (sync_run_id) REFERENCES source_sync_runs (id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_charms_deleted_at ON charms (deleted_at);
CREATE INDEX IF NOT EXISTS idx_charms_source_slug ON charms (source_slug);
CREATE INDEX IF NOT EXISTS idx_charms_sync_run_id ON charms (sync_run_id);
CREATE INDEX IF NOT EXISTS idx_charms_rarity ON charms (rarity);
CREATE INDEX IF NOT EXISTS idx_charms_sort_order ON charms (sort_order);

CREATE TABLE IF NOT EXISTS food_skills (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz NULL,
    provider text NOT NULL,
    source_id text NOT NULL,
    source_url text NULL,
    source_slug text NULL,
    fetched_at timestamptz NOT NULL DEFAULT now(),
    sync_run_id uuid NULL,
    food_category text NOT NULL,
    max_level smallint NOT NULL DEFAULT 1,
    base_duration_seconds integer NULL,
    base_activation_percent smallint NULL,
    sort_order integer NOT NULL DEFAULT 0,
    details_json jsonb NOT NULL DEFAULT '{}'::jsonb,
    CONSTRAINT uq_food_skills_provider_source_id UNIQUE (provider, source_id),
    CONSTRAINT fk_food_skills_sync_run_id FOREIGN KEY (sync_run_id) REFERENCES source_sync_runs (id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_food_skills_deleted_at ON food_skills (deleted_at);
CREATE INDEX IF NOT EXISTS idx_food_skills_source_slug ON food_skills (source_slug);
CREATE INDEX IF NOT EXISTS idx_food_skills_sync_run_id ON food_skills (sync_run_id);
CREATE INDEX IF NOT EXISTS idx_food_skills_food_category ON food_skills (food_category);
CREATE INDEX IF NOT EXISTS idx_food_skills_sort_order ON food_skills (sort_order);

CREATE TABLE IF NOT EXISTS kinsects (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz NULL,
    provider text NOT NULL,
    source_id text NOT NULL,
    source_url text NULL,
    source_slug text NULL,
    fetched_at timestamptz NOT NULL DEFAULT now(),
    sync_run_id uuid NULL,
    kinsect_type text NOT NULL,
    attack_type text NULL,
    powder_type text NULL,
    rarity smallint NOT NULL DEFAULT 0,
    power_value smallint NOT NULL DEFAULT 0,
    speed_value smallint NOT NULL DEFAULT 0,
    heal_value smallint NOT NULL DEFAULT 0,
    stamina_value smallint NOT NULL DEFAULT 0,
    element_type text NULL,
    element_value integer NOT NULL DEFAULT 0,
    sort_order integer NOT NULL DEFAULT 0,
    details_json jsonb NOT NULL DEFAULT '{}'::jsonb,
    CONSTRAINT uq_kinsects_provider_source_id UNIQUE (provider, source_id),
    CONSTRAINT fk_kinsects_sync_run_id FOREIGN KEY (sync_run_id) REFERENCES source_sync_runs (id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_kinsects_deleted_at ON kinsects (deleted_at);
CREATE INDEX IF NOT EXISTS idx_kinsects_source_slug ON kinsects (source_slug);
CREATE INDEX IF NOT EXISTS idx_kinsects_sync_run_id ON kinsects (sync_run_id);
CREATE INDEX IF NOT EXISTS idx_kinsects_kinsect_type ON kinsects (kinsect_type);
CREATE INDEX IF NOT EXISTS idx_kinsects_attack_type ON kinsects (attack_type);
CREATE INDEX IF NOT EXISTS idx_kinsects_powder_type ON kinsects (powder_type);
CREATE INDEX IF NOT EXISTS idx_kinsects_rarity ON kinsects (rarity);
CREATE INDEX IF NOT EXISTS idx_kinsects_element_type ON kinsects (element_type);
CREATE INDEX IF NOT EXISTS idx_kinsects_sort_order ON kinsects (sort_order);

CREATE TABLE IF NOT EXISTS skill_translations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    skill_id uuid NOT NULL,
    language_id uuid NOT NULL,
    name text NOT NULL,
    description text NULL,
    effect_summary text NULL,
    slug text NULL,
    CONSTRAINT uq_skill_translations_skill_language UNIQUE (skill_id, language_id),
    CONSTRAINT fk_skill_translations_skill_id FOREIGN KEY (skill_id) REFERENCES skills (id) ON DELETE CASCADE,
    CONSTRAINT fk_skill_translations_language_id FOREIGN KEY (language_id) REFERENCES languages (id) ON DELETE RESTRICT
);
CREATE INDEX IF NOT EXISTS idx_skill_translations_name ON skill_translations (name);
CREATE INDEX IF NOT EXISTS idx_skill_translations_slug ON skill_translations (slug);
CREATE INDEX IF NOT EXISTS idx_skill_translations_skill_id ON skill_translations (skill_id);
CREATE INDEX IF NOT EXISTS idx_skill_translations_language_id ON skill_translations (language_id);

CREATE TABLE IF NOT EXISTS item_translations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    item_id uuid NOT NULL,
    language_id uuid NOT NULL,
    name text NOT NULL,
    description text NULL,
    flavor_text text NULL,
    slug text NULL,
    CONSTRAINT uq_item_translations_item_language UNIQUE (item_id, language_id),
    CONSTRAINT fk_item_translations_item_id FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE CASCADE,
    CONSTRAINT fk_item_translations_language_id FOREIGN KEY (language_id) REFERENCES languages (id) ON DELETE RESTRICT
);
CREATE INDEX IF NOT EXISTS idx_item_translations_name ON item_translations (name);
CREATE INDEX IF NOT EXISTS idx_item_translations_slug ON item_translations (slug);
CREATE INDEX IF NOT EXISTS idx_item_translations_item_id ON item_translations (item_id);
CREATE INDEX IF NOT EXISTS idx_item_translations_language_id ON item_translations (language_id);

CREATE TABLE IF NOT EXISTS weapon_translations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    weapon_id uuid NOT NULL,
    language_id uuid NOT NULL,
    name text NOT NULL,
    description text NULL,
    slug text NULL,
    CONSTRAINT uq_weapon_translations_weapon_language UNIQUE (weapon_id, language_id),
    CONSTRAINT fk_weapon_translations_weapon_id FOREIGN KEY (weapon_id) REFERENCES weapons (id) ON DELETE CASCADE,
    CONSTRAINT fk_weapon_translations_language_id FOREIGN KEY (language_id) REFERENCES languages (id) ON DELETE RESTRICT
);
CREATE INDEX IF NOT EXISTS idx_weapon_translations_name ON weapon_translations (name);
CREATE INDEX IF NOT EXISTS idx_weapon_translations_slug ON weapon_translations (slug);
CREATE INDEX IF NOT EXISTS idx_weapon_translations_weapon_id ON weapon_translations (weapon_id);
CREATE INDEX IF NOT EXISTS idx_weapon_translations_language_id ON weapon_translations (language_id);

CREATE TABLE IF NOT EXISTS armor_piece_translations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    armor_piece_id uuid NOT NULL,
    language_id uuid NOT NULL,
    name text NOT NULL,
    description text NULL,
    slug text NULL,
    CONSTRAINT uq_armor_piece_translations_piece_language UNIQUE (armor_piece_id, language_id),
    CONSTRAINT fk_armor_piece_translations_piece_id FOREIGN KEY (armor_piece_id) REFERENCES armor_pieces (id) ON DELETE CASCADE,
    CONSTRAINT fk_armor_piece_translations_language_id FOREIGN KEY (language_id) REFERENCES languages (id) ON DELETE RESTRICT
);
CREATE INDEX IF NOT EXISTS idx_armor_piece_translations_name ON armor_piece_translations (name);
CREATE INDEX IF NOT EXISTS idx_armor_piece_translations_slug ON armor_piece_translations (slug);
CREATE INDEX IF NOT EXISTS idx_armor_piece_translations_piece_id ON armor_piece_translations (armor_piece_id);
CREATE INDEX IF NOT EXISTS idx_armor_piece_translations_language_id ON armor_piece_translations (language_id);

CREATE TABLE IF NOT EXISTS decorations_translations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    decoration_id uuid NOT NULL,
    language_id uuid NOT NULL,
    name text NOT NULL,
    description text NULL,
    slug text NULL,
    CONSTRAINT uq_decorations_translations_decoration_language UNIQUE (decoration_id, language_id),
    CONSTRAINT fk_decorations_translations_decoration_id FOREIGN KEY (decoration_id) REFERENCES decorations (id) ON DELETE CASCADE,
    CONSTRAINT fk_decorations_translations_language_id FOREIGN KEY (language_id) REFERENCES languages (id) ON DELETE RESTRICT
);
CREATE INDEX IF NOT EXISTS idx_decorations_translations_name ON decorations_translations (name);
CREATE INDEX IF NOT EXISTS idx_decorations_translations_slug ON decorations_translations (slug);
CREATE INDEX IF NOT EXISTS idx_decorations_translations_decoration_id ON decorations_translations (decoration_id);
CREATE INDEX IF NOT EXISTS idx_decorations_translations_language_id ON decorations_translations (language_id);

CREATE TABLE IF NOT EXISTS charm_translations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    charm_id uuid NOT NULL,
    language_id uuid NOT NULL,
    name text NOT NULL,
    description text NULL,
    slug text NULL,
    CONSTRAINT uq_charm_translations_charm_language UNIQUE (charm_id, language_id),
    CONSTRAINT fk_charm_translations_charm_id FOREIGN KEY (charm_id) REFERENCES charms (id) ON DELETE CASCADE,
    CONSTRAINT fk_charm_translations_language_id FOREIGN KEY (language_id) REFERENCES languages (id) ON DELETE RESTRICT
);
CREATE INDEX IF NOT EXISTS idx_charm_translations_name ON charm_translations (name);
CREATE INDEX IF NOT EXISTS idx_charm_translations_slug ON charm_translations (slug);
CREATE INDEX IF NOT EXISTS idx_charm_translations_charm_id ON charm_translations (charm_id);
CREATE INDEX IF NOT EXISTS idx_charm_translations_language_id ON charm_translations (language_id);

CREATE TABLE IF NOT EXISTS food_skill_translations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    food_skill_id uuid NOT NULL,
    language_id uuid NOT NULL,
    name text NOT NULL,
    description text NULL,
    slug text NULL,
    CONSTRAINT uq_food_skill_translations_food_skill_language UNIQUE (food_skill_id, language_id),
    CONSTRAINT fk_food_skill_translations_food_skill_id FOREIGN KEY (food_skill_id) REFERENCES food_skills (id) ON DELETE CASCADE,
    CONSTRAINT fk_food_skill_translations_language_id FOREIGN KEY (language_id) REFERENCES languages (id) ON DELETE RESTRICT
);
CREATE INDEX IF NOT EXISTS idx_food_skill_translations_name ON food_skill_translations (name);
CREATE INDEX IF NOT EXISTS idx_food_skill_translations_slug ON food_skill_translations (slug);
CREATE INDEX IF NOT EXISTS idx_food_skill_translations_food_skill_id ON food_skill_translations (food_skill_id);
CREATE INDEX IF NOT EXISTS idx_food_skill_translations_language_id ON food_skill_translations (language_id);

CREATE TABLE IF NOT EXISTS kinsect_translations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    kinsect_id uuid NOT NULL,
    language_id uuid NOT NULL,
    name text NOT NULL,
    description text NULL,
    slug text NULL,
    CONSTRAINT uq_kinsect_translations_kinsect_language UNIQUE (kinsect_id, language_id),
    CONSTRAINT fk_kinsect_translations_kinsect_id FOREIGN KEY (kinsect_id) REFERENCES kinsects (id) ON DELETE CASCADE,
    CONSTRAINT fk_kinsect_translations_language_id FOREIGN KEY (language_id) REFERENCES languages (id) ON DELETE RESTRICT
);
CREATE INDEX IF NOT EXISTS idx_kinsect_translations_name ON kinsect_translations (name);
CREATE INDEX IF NOT EXISTS idx_kinsect_translations_slug ON kinsect_translations (slug);
CREATE INDEX IF NOT EXISTS idx_kinsect_translations_kinsect_id ON kinsect_translations (kinsect_id);
CREATE INDEX IF NOT EXISTS idx_kinsect_translations_language_id ON kinsect_translations (language_id);

CREATE TABLE IF NOT EXISTS weapon_skills (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    weapon_id uuid NOT NULL,
    skill_id uuid NOT NULL,
    sort_order smallint NOT NULL DEFAULT 0,
    level smallint NOT NULL DEFAULT 1,
    source_type text NOT NULL DEFAULT 'base',
    CONSTRAINT uq_weapon_skills_weapon_skill_order UNIQUE (weapon_id, skill_id, sort_order),
    CONSTRAINT fk_weapon_skills_weapon_id FOREIGN KEY (weapon_id) REFERENCES weapons (id) ON DELETE CASCADE,
    CONSTRAINT fk_weapon_skills_skill_id FOREIGN KEY (skill_id) REFERENCES skills (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_weapon_skills_weapon_id ON weapon_skills (weapon_id);
CREATE INDEX IF NOT EXISTS idx_weapon_skills_skill_id ON weapon_skills (skill_id);

CREATE TABLE IF NOT EXISTS armor_piece_skills (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    armor_piece_id uuid NOT NULL,
    skill_id uuid NOT NULL,
    sort_order smallint NOT NULL DEFAULT 0,
    level smallint NOT NULL DEFAULT 1,
    CONSTRAINT uq_armor_piece_skills_piece_skill_order UNIQUE (armor_piece_id, skill_id, sort_order),
    CONSTRAINT fk_armor_piece_skills_piece_id FOREIGN KEY (armor_piece_id) REFERENCES armor_pieces (id) ON DELETE CASCADE,
    CONSTRAINT fk_armor_piece_skills_skill_id FOREIGN KEY (skill_id) REFERENCES skills (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_armor_piece_skills_piece_id ON armor_piece_skills (armor_piece_id);
CREATE INDEX IF NOT EXISTS idx_armor_piece_skills_skill_id ON armor_piece_skills (skill_id);

CREATE TABLE IF NOT EXISTS decoration_skills (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    decoration_id uuid NOT NULL,
    skill_id uuid NOT NULL,
    sort_order smallint NOT NULL DEFAULT 0,
    level smallint NOT NULL DEFAULT 1,
    CONSTRAINT uq_decoration_skills_decoration_skill_order UNIQUE (decoration_id, skill_id, sort_order),
    CONSTRAINT fk_decoration_skills_decoration_id FOREIGN KEY (decoration_id) REFERENCES decorations (id) ON DELETE CASCADE,
    CONSTRAINT fk_decoration_skills_skill_id FOREIGN KEY (skill_id) REFERENCES skills (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_decoration_skills_decoration_id ON decoration_skills (decoration_id);
CREATE INDEX IF NOT EXISTS idx_decoration_skills_skill_id ON decoration_skills (skill_id);

CREATE TABLE IF NOT EXISTS charm_skills (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    charm_id uuid NOT NULL,
    skill_id uuid NOT NULL,
    rank_required smallint NOT NULL DEFAULT 1,
    level smallint NOT NULL DEFAULT 1,
    sort_order smallint NOT NULL DEFAULT 0,
    CONSTRAINT uq_charm_skills_charm_skill_rank UNIQUE (charm_id, skill_id, rank_required),
    CONSTRAINT fk_charm_skills_charm_id FOREIGN KEY (charm_id) REFERENCES charms (id) ON DELETE CASCADE,
    CONSTRAINT fk_charm_skills_skill_id FOREIGN KEY (skill_id) REFERENCES skills (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_charm_skills_charm_id ON charm_skills (charm_id);
CREATE INDEX IF NOT EXISTS idx_charm_skills_skill_id ON charm_skills (skill_id);

CREATE TABLE IF NOT EXISTS food_skill_levels (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    food_skill_id uuid NOT NULL,
    level smallint NOT NULL DEFAULT 1,
    duration_seconds integer NULL,
    activation_percent smallint NULL,
    effect_value_text text NULL,
    details_json jsonb NOT NULL DEFAULT '{}'::jsonb,
    CONSTRAINT uq_food_skill_levels_food_skill_level UNIQUE (food_skill_id, level),
    CONSTRAINT fk_food_skill_levels_food_skill_id FOREIGN KEY (food_skill_id) REFERENCES food_skills (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_food_skill_levels_food_skill_id ON food_skill_levels (food_skill_id);

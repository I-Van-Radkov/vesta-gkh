CREATE TABLE channels (
	id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title               TEXT NOT NULL,
    url                 TEXT NOT NULL UNIQUE,
    source_name         TEXT NOT NULL,
    description         TEXT,
    language            TEXT DEFAULT 'ru',
    last_build_date     TIMESTAMPTZ,
    generator           TEXT,
    image_url           TEXT,
    is_active           BOOLEAN DEFAULT TRUE,
    priority            INTEGER DEFAULT 10,
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_channels_active_priority ON channels(is_active, priority);
CREATE INDEX idx_channels_url ON channels(url);

CREATE TABLE news (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    channel_id          UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    guid                TEXT UNIQUE,
    title               TEXT NOT NULL,
    link                TEXT NOT NULL UNIQUE,
    description         TEXT,
    full_text           TEXT,
    pub_date            TIMESTAMPTZ NOT NULL,
    fetched_at          TIMESTAMPTZ DEFAULT NOW(),
    enclosure_url       TEXT,
    enclosure_type      TEXT,
    enclosure_length    BIGINT,
    media_group         JSONB,
    relevance_score     REAL,
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_news_channel_pubdate ON news(channel_id, pub_date DESC);
CREATE INDEX idx_news_link ON news(link);

CREATE TABLE gkh_keywords (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    keyword     TEXT NOT NULL UNIQUE,
    weight      REAL NOT NULL DEFAULT 30.0,
    category    TEXT,
    is_active   BOOLEAN DEFAULT TRUE,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE gkh_regex_rules (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pattern     TEXT NOT NULL UNIQUE,
    bonus_score REAL NOT NULL DEFAULT 55.0,
    is_active   BOOLEAN DEFAULT TRUE,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE gkh_categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL UNIQUE,
    bonus_per_hit    REAL DEFAULT 12.0,
    min_hits_for_bonus INTEGER DEFAULT 1,
    big_bonus        REAL DEFAULT 35.0,
    description      TEXT
);

CREATE TABLE gkh_category_keywords (
    category_id UUID NOT NULL REFERENCES gkh_categories(id) ON DELETE CASCADE,
    keyword_id  UUID NOT NULL REFERENCES gkh_keywords(id) ON DELETE CASCADE,
    PRIMARY KEY (category_id, keyword_id)
);

CREATE TABLE gkh_negative_regex (
    id          BIGSERIAL PRIMARY KEY,
    pattern     TEXT NOT NULL UNIQUE,
    penalty     REAL NOT NULL DEFAULT -45.0,
    description TEXT,
    is_active   BOOLEAN DEFAULT TRUE,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW()
);
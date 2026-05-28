CREATE TABLE IF NOT EXISTS games (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    title       TEXT        NOT NULL,
    release_date DATE        NOT NULL,
    price       NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    rating      INT         NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

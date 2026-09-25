CREATE TABLE IF NOT EXISTS user_plants (
    id UUID PRIMARY KEY,
    species_id INTEGER NOT NULL,
    user_id UUID NOT NULL,
    name TEXT NOT NULL,
    watering_interval_days INTEGER NOT NULL,
    status TEXT NOT NULL,
    last_watered_at TIMESTAMPTZ NOT NULL,
    next_watering_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_user_plants_user_id ON user_plants (user_id);

CREATE TABLE IF NOT EXISTS watering_events (
    id UUID PRIMARY KEY,
    user_plant_id UUID NOT NULL REFERENCES user_plants (id) ON DELETE CASCADE,
    watered_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

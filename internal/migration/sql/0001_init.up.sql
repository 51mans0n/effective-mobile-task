CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE people (
                        id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                        name            TEXT NOT NULL,
                        surname         TEXT NOT NULL,
                        patronymic      TEXT,
                        age             INT,
                        gender          TEXT,
                        country_code    CHAR(2),
                        nat_probability REAL,
                        created_at      TIMESTAMPTZ DEFAULT now(),
                        updated_at      TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_people_fulltext
    ON people USING gin (to_tsvector('simple', name || ' ' || surname));


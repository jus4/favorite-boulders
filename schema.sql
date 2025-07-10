CREATE TABLE areas (
  area_id BIGINT PRIMARY KEY NOT NULL,
  name    VARCHAR,
  country VARCHAR
);

CREATE TABLE route (
  route_id        BIGINT PRIMARY KEY NOT NULL,
  sector_id       BIGINT,
  crag_id         BIGINT,
  title           TEXT,
  grade           TEXT,
  route_type      TEXT,
  FOREIGN KEY (crag_id) REFERENCES crag(crag_id),
  FOREIGN KEY (sector_id) REFERENCES sector(sector_id)
);

CREATE TABLE sector (
  sector_id     BIGINT PRIMARY KEY NOT NULL,
  crag_id       BIGINT,
  name          TEXT,
  url           TEXT,
  latitude      DOUBLE PRECISION,
  longitude     DOUBLE PRECISION,
  FOREIGN KEY (crag_id) REFERENCES crag(crag_id)
);

CREATE TABLE crag (
  crag_id     BIGINT PRIMARY KEY NOT NULL,
  area_id     BIGINT,
  param_id    TEXT,
  premium_topo_active   BOOLEAN,
  latitude    DOUBLE PRECISION,
  longitude     DOUBLE PRECISION,
  name        TEXT,
  info        TEXT,
  url         TEXT,
  FOREIGN KEY (crag_id) REFERENCES crag(crag_id),
  FOREIGN KEY (crag_id) REFERENCES crag(crag_id)
);

-- CREATE TABLE area (
--   area_id     BIGINT PRIMARY KEY,
--   param_id    TEXT,
--   name        TEXT,
--   country     TEXT
-- );

CREATE TABLE users (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- unique UUID as primary key
  email       TEXT NOT NULL UNIQUE,                    -- required and unique
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(), -- track when the user was created
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()  -- track updates (for future use)
);

CREATE TABLE user_favourite_lists  (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name        VARCHAR,
  description TEXT,
  owner       TEXT NOT NULL,
  FOREIGN KEY (owner) REFERENCES users(email)
);

CREATE TABLE user_favourite_climb_items  (
  id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  route_id              BIGINT,
  favourite_climb_list  UUID,
  FOREIGN KEY (favourite_climb_list) REFERENCES favourite_climb_list(id) ON DELETE CASCADE,
  FOREIGN KEY (route_id) REFERENCES route(route_id)
);

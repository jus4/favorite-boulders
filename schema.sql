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

CREATE TABLE area (
  area_id     BIGINT PRIMARY KEY,
  param_id    TEXT,
  name        TEXT,
  country     TEXT
);

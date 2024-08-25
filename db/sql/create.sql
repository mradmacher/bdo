DROP TABLE IF EXISTS materials;
DROP TABLE IF EXISTS capabilities;
DROP TABLE IF EXISTS installations;

CREATE TABLE installations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  address_line1 TEXT NOT NULL,
  address_line2 TEXT,
  lat TEXT,
  lng TEXT,
  state_code VARCHAR(2) NOT NULL
);

CREATE TABLE capabilities (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  waste_code VARCHAR(6) NOT NULL,
  dangerous BOOLEAN NOT NULL DEFAULT FALSE,
  process_code VARCHAR(2) NOT NULL,
  activity_code VARCHAR(2),
  quantity INTEGER,
  installation_id INTEGER,
  FOREIGN KEY(installation_id) REFERENCES installations(id)
);
CREATE TABLE materials (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  code VARCHAR NOT NULL,
  capability_id INTEGER,
  FOREIGN KEY(capability_id) REFERENCES capabilities(id)
);
CREATE INDEX capability_waste_code ON capabilities(waste_code);
CREATE INDEX capability_process_code ON capabilities(process_code);
CREATE UNIQUE INDEX material_capability ON materials(capability_id, code);

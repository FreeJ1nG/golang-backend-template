CREATE TABLE IF NOT EXISTS Account (
  id SERIAL NOT NULL UNIQUE,
  username VARCHAR(64) NOT NULL UNIQUE,
  first_name VARCHAR(128) NOT NULL,
  last_name VARCHAR(128),
  password_hash VARCHAR(256) NOT NULL,
  PRIMARY KEY (id, username)
);
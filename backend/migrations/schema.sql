CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  name text NOT NULL,
  password text NOT NULL
);

CREATE TABLE urls (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  orginal_url text NOT NULL,
  short_url text NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);
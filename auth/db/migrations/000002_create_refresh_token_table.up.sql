CREATE TABLE IF NOT EXISTS refresh_tokens ( 
  id serial PRIMARY KEY,
  token text NOT NULL,
  user_id integer NOT NULL REFERENCES users (id),
  expires int NOT NULL
);

CREATE INDEX IF NOT EXISTS refresh_tokens_token_idx ON refresh_tokens (token);
-- Create table if it doesnt exist, refresh tokens is the table name, it contains expiry as seconds since epoch, and a foreign key referecing users user_id
CREATE TABLE IF NOT EXISTS refresh_tokens ( 
  id serial PRIMARY KEY,
  token text NOT NULL,
  user_id uuid NOT NULL REFERENCES users (id),
  expires int NOT NULL
);

CREATE INDEX IF NOT EXISTS refresh_tokens_token_idx ON refresh_tokens (token);
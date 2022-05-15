CREATE TABLE IF NOT EXISTS users(
   id serial PRIMARY KEY,
   password VARCHAR (100) NOT NULL,
   email VARCHAR (300) UNIQUE NOT NULL
);

CREATE INDEX IF NOT EXISTS users_email_idx ON users (email);
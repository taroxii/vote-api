BEGIN;
CREATE TABLE IF NOT EXISTS users(
   id serial PRIMARY KEY,
   username VARCHAR (50) UNIQUE NOT NULL,
   created_at timestamp,
   updated_at timestamp
);
COMMIT;
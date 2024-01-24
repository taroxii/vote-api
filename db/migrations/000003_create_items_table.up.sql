BEGIN;
CREATE TABLE IF NOT EXISTS items(
    id serial PRIMARY KEY, 
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    vote_count int default 0,
    created_at timestamp,
    updated_at timestamp
);

COMMIT;
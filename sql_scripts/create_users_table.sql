CREATE TABLE users (
  user_id SERIAL PRIMARY KEY,
  email TEXT NOT NULL UNIQUE
);

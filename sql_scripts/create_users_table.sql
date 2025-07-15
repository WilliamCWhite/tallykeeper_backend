CREATE TABLE users (
  user_id SERIAL PRIMARY KEY,
  username VARCHAR(32) NOT NULL,
  user_key VARCHAR(20) NOT NULL
);

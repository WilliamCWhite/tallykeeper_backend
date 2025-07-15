CREATE TABLE entries (
  entry_id SERIAL PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  score INT NOT NULL,
  list_id INT NOT NULL,
  FOREIGN KEY (list_id)
    REFERENCES lists(list_id)
      ON DELETE CASCADE
);

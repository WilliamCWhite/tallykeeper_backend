INSERT INTO users(email)
VALUES
  ('test_user_1'),
  ('test_user_2');

INSERT INTO lists(title, user_id)
VALUES
  ('test_list_1-1', 1),
  ('test_list_1-2', 1),
  ('test_list_1-3', 1),
  ('test_list_2-1', 2),
  ('test_list_2-2', 2);

INSERT INTO entries(name, score, list_id)
VALUES
  ('entry_1-1-1', 5, 1),
  ('entry_1-1-2', 10, 1),
  ('entry_1-1-3', 15, 1),
  ('entry_1-2-1', 20, 2),
  ('entry_1-2-2', 25, 2),
  ('entry_1-3-1', 30, 3),
  ('entry_2-1-1', 35, 4),
  ('entry_2-1-2', 40, 4),
  ('entry_2-2-1', 45, 5);

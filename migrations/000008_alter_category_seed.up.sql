INSERT INTO "categories" (id, category_name, is_income, category_icon)
VALUES 
  ('01ARZ3NDEKTSV4RRFFQ69G5FAV', 'Transfer IN', true, 30),
  ('01ARZ3NDEKTSV4RRFFQ69G5FAW', 'Transfer OUT', false, 31)
ON CONFLICT (id) DO NOTHING;
INSERT INTO "categories" (id, category_name, is_income, category_icon)
VALUES 
  ('00000000-0000-0000-0000-000000000000', 'Transfer IN', true, 30),
  ('00000000-0000-0000-0000-000000000001', 'Transfer OUT', false, 31)
ON CONFLICT (id) DO NOTHING;
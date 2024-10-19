INSERT INTO public.pockets (id,owner_id,editor_id,watcher_id,pocket_name,balance,currency,icon,"level","version") VALUES
	 ('00000000000000000000000000',NULL,'{}','{}','SYSTEM',0,'',0,0,0)
ON CONFLICT (id) DO NOTHING;

INSERT INTO public.categories (id,pocket_id,category_name,is_income,category_icon) VALUES
    ('01ARZ3NDEKTSV4RRFFQ69G5FAV','00000000000000000000000000','Transfer IN',true,100),
    ('01J7QVGPM105H3BTRWSMRMQQ8F','00000000000000000000000000','Salary',true,1),
    ('01J7QVGPM105H3BTRWSP3QPERD','00000000000000000000000000','Grants',true,2),
    ('01J7QVGPM105H3BTRWSR8FFT0D','00000000000000000000000000','Refunds',true,3),
    ('01J7QVGPM105H3BTRWSSBPHPJV','00000000000000000000000000','Sale',true,4),
    ('01J7QVGPM105H3BTRWSTD080EF','00000000000000000000000000','Rental',true,5),
    ('01FHT3EFT7YJ3QW0D2XVWX08TY','00000000000000000000000000','Bonus',true,6),
    ('01FHT3EFT7YJ3QW0D2XVWX08TZ','00000000000000000000000000','Investment Returns',true,7),
    ('01FHT3EFT7YJ3QW0D2XVWX08TX','00000000000000000000000000','Other Income',true,8)
  ON CONFLICT (id) DO NOTHING;

INSERT INTO public.categories (id,pocket_id,category_name,is_income,category_icon) VALUES
     ('01ARZ3NDEKTSV4RRFFQ69G5FAW','00000000000000000000000000','Transfer OUT',false,101),
     ('01J7QVGPM105H3BTRWSVNMJBKR','00000000000000000000000000','Baby',false,31),
     ('01J7QVGPM105H3BTRWSZ6J4KP4','00000000000000000000000000','Beauty',false,32),
     ('01J7QVGPM105H3BTRWT2Q28RJP','00000000000000000000000000','Bills',false,33),
     ('01J7QVGPM105H3BTRWT3BPZE57','00000000000000000000000000','Vehicle',false,34),
     ('01J7QVGPM105H3BTRWT7B4EAHD','00000000000000000000000000','Clothing',false,35),
     ('01J7QVGPM105H3BTRWT8Y3SK34','00000000000000000000000000','Education',false,36),
     ('01J7QVGPM105H3BTRWT97ZKHVP','00000000000000000000000000','Electronics',false,37),
     ('01J7QVGPM105H3BTRWTAGG7998','00000000000000000000000000','Entertainment',false,38),
     ('01J7QVGPM105H3BTRWTAWKXF0N','00000000000000000000000000','Food',false,39),
     ('01J7QVGPM105H3BTRWTEQXD8CT','00000000000000000000000000','Health',false,40),
     ('01J7QVGPM105H3BTRWTGQY7K80','00000000000000000000000000','Home',false,41),
     ('01J7QVGPM105H3BTRWTHZ1Q8T8','00000000000000000000000000','Insurance',false,42),
     ('01J7QVGPM105H3BTRWTNTGND1R','00000000000000000000000000','Shopping',false,43),
     ('01J7QVGPM105H3BTRWTSFAGG7R','00000000000000000000000000','Social',false,44),
     ('01J7QVGPM105H3BTRWTTXJADTJ','00000000000000000000000000','Sport',false,45),
     ('01J7QVGPM105H3BTRWTTYZXTNP','00000000000000000000000000','Tax',false,46),
     ('01J7QVGPM105H3BTRWTXSYV0NV','00000000000000000000000000','Telephone',false,47),
     ('01J7QVGPM105H3BTRWTYPHR2RA','00000000000000000000000000','Internet',false,48),
     ('01J7QVGPM105H3BTRWV2379ARJ','00000000000000000000000000','Transportation',false,50),
     ('01J7QVGPM105H3BTRWV3T3JEB6','00000000000000000000000000','Work',false,51),
     ('01FHT3EFT7YJ3QW0D2XVWX08TR','00000000000000000000000000','Donation',false,52),
     ('01FHT3EFT7YJ3QW0D2XVWX08TS','00000000000000000000000000','Investment',false,53),
     ('01FHT3EFT7YJ3QW0D2XVWX08TT','00000000000000000000000000','Family and Friends',false,54),
     ('01FHT3EFT7YJ3QW0D2XVWX08TV','00000000000000000000000000','Pets',false,55),
     ('01FHT3EFT7YJ3QW0D2XVWX08TW','00000000000000000000000000','Technology',false,56),
     ('01FHT3EFT7YJ3QW0D2XVWX08TU','00000000000000000000000000','Other',false,57)
ON CONFLICT (id) DO NOTHING;
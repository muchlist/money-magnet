INSERT INTO public.pockets (id,owner_id,editor_id,watcher_id,pocket_name,balance,currency,icon,"level","version") VALUES
	 ('00000000000000000000000000',NULL,'{}','{}','SYSTEM',0,'',0,0,0)
ON CONFLICT (id) DO NOTHING;

INSERT INTO public.categories (id,pocket_id,category_name,is_income,category_icon,default_spend_type) VALUES
    ('01ARZ3NDEKTSV4RRFFQ69G5FAV','00000000000000000000000000','Transfer IN',true,100,0),
    ('01ARZ3NDEKTSV4RRFFQ69G5FAX','00000000000000000000000000','Saving IN',true,102,3),
    ('01J7QVGPM105H3BTRWSMRMQQ8F','00000000000000000000000000','Salary',true,1,0),
    ('01J7QVGPM105H3BTRWSP3QPERD','00000000000000000000000000','Grants',true,2,0),
    ('01J7QVGPM105H3BTRWSR8FFT0D','00000000000000000000000000','Refunds',true,3,0),
    ('01J7QVGPM105H3BTRWSSBPHPJV','00000000000000000000000000','Sale',true,4,0),
    ('01J7QVGPM105H3BTRWSTD080EF','00000000000000000000000000','Rental',true,5,0),
    ('01FHT3EFT7YJ3QW0D2XVWX08TY','00000000000000000000000000','Bonus',true,6,0),
    ('01FHT3EFT7YJ3QW0D2XVWX08TZ','00000000000000000000000000','Investment Returns',true,7,0),
    ('01FHT3EFT7YJ3QW0D2XVWX08TX','00000000000000000000000000','Other Income',true,8,0)
    ON CONFLICT (id) DO UPDATE 
    SET pocket_id = EXCLUDED.pocket_id,
        category_name = EXCLUDED.category_name,
        is_income = EXCLUDED.is_income,
        category_icon = EXCLUDED.category_icon,
        default_spend_type = EXCLUDED.default_spend_type;


INSERT INTO public.categories (id,pocket_id,category_name,is_income,category_icon,default_spend_type) VALUES
    ('01ARZ3NDEKTSV4RRFFQ69G5FAW','00000000000000000000000000','Transfer OUT',false,101,0),
    ('01ARZ3NDEKTSV4RRFFQ69G5FAY','00000000000000000000000000','Saving OUT',false,103,3),
    ('01J7QVGPM105H3BTRWSVNMJBKR','00000000000000000000000000','Baby',false,31,1),
    ('01J7QVGPM105H3BTRWSZ6J4KP4','00000000000000000000000000','Beauty',false,32,2),
    ('01J7QVGPM105H3BTRWT2Q28RJP','00000000000000000000000000','Bills',false,33,1),
    ('01J7QVGPM105H3BTRWT3BPZE57','00000000000000000000000000','Vehicle',false,34,1),
    ('01J7QVGPM105H3BTRWT7B4EAHD','00000000000000000000000000','Clothing',false,35,1),
    ('01J7QVGPM105H3BTRWT8Y3SK34','00000000000000000000000000','Education',false,36,1),
    ('01J7QVGPM105H3BTRWT97ZKHVP','00000000000000000000000000','Electronics',false,37,2),
    ('01J7QVGPM105H3BTRWTAGG7998','00000000000000000000000000','Entertainment',false,38,2),
    ('01J7QVGPM105H3BTRWTAWKXF0N','00000000000000000000000000','Food',false,39,1),
    ('01J7QVGPM105H3BTRWTEQXD8CT','00000000000000000000000000','Health',false,40,1),
    ('01J7QVGPM105H3BTRWTGQY7K80','00000000000000000000000000','Home',false,41,1),
    ('01J7QVGPM105H3BTRWTHZ1Q8T8','00000000000000000000000000','Insurance',false,42,1),
    ('01J7QVGPM105H3BTRWTNTGND1R','00000000000000000000000000','Shopping',false,43,2),
    ('01J7QVGPM105H3BTRWTSFAGG7R','00000000000000000000000000','Social',false,44,2),
    ('01J7QVGPM105H3BTRWTTXJADTJ','00000000000000000000000000','Sport',false,45,2),
    ('01J7QVGPM105H3BTRWTTYZXTNP','00000000000000000000000000','Tax',false,46,1),
    ('01J7QVGPM105H3BTRWTXSYV0NV','00000000000000000000000000','Telephone',false,47,1),
    ('01J7QVGPM105H3BTRWTYPHR2RA','00000000000000000000000000','Internet',false,48,1),
    ('01J7QVGPM105H3BTRWV2379ARJ','00000000000000000000000000','Transportation',false,50,1),
    ('01J7QVGPM105H3BTRWV3T3JEB6','00000000000000000000000000','Work',false,51,1),
    ('01FHT3EFT7YJ3QW0D2XVWX08TR','00000000000000000000000000','Donation',false,52,2),
    ('01FHT3EFT7YJ3QW0D2XVWX08TS','00000000000000000000000000','Investment',false,53,3),
    ('01FHT3EFT7YJ3QW0D2XVWX08TT','00000000000000000000000000','Family',false,54,1),
    ('01FHT3EFT7YJ3QW0D2XVWX08TV','00000000000000000000000000','Pets',false,55,2),
    ('01FHT3EFT7YJ3QW0D2XVWX08TW','00000000000000000000000000','Technology',false,56,2),
    ('01FHT3EFT7YJ3QW0D2XVWX08TU','00000000000000000000000000','Other',false,57,2),
    ('01J7QVGPM105H3BTRWTAGG7999','00000000000000000000000000','Game',false,58,2),
    ('01FHT3EFT7YJ3QW0D2XVWX09TT','00000000000000000000000000','Friend',false,59,2)
    ON CONFLICT (id) DO UPDATE 
    SET pocket_id = EXCLUDED.pocket_id,
        category_name = EXCLUDED.category_name,
        is_income = EXCLUDED.is_income,
        category_icon = EXCLUDED.category_icon,
        default_spend_type = EXCLUDED.default_spend_type;
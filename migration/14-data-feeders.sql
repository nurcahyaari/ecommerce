INSERT INTO ecommerce.`user` (name,email,phone,password,is_actived,is_removed,created_at,created_by,updated_at,updated_by,deleted_at,deleted_by) VALUES
	 ('test','test@test.com','12345','$2a$12$YBriH.WM/BncqUxRHdGUs.K.c5/NsCTOopnszKGk6vXnPS./15Kei
',1,0,'2024-10-02 14:53:04.0',NULL,'2024-10-02 14:53:04.0',NULL,NULL,NULL); -- the password is password

INSERT INTO ecommerce.user_address (user_id,full_address,is_actived,is_removed,created_at,created_by,updated_at,updated_by,deleted_at,deleted_by) VALUES
	 (1,'Testing, Testing',1,0,'2024-10-05 01:20:13.0',NULL,'2024-10-05 01:20:13.0',NULL,NULL,NULL);

INSERT INTO ecommerce.store (name,is_actived,is_removed,created_at,created_by,updated_at,updated_by,deleted_at,deleted_by) VALUES
	 ('test',1,0,'2024-10-03 22:41:27.0',NULL,'2024-10-03 22:41:27.0',NULL,NULL,NULL);

INSERT INTO ecommerce.warehouse (store_id,name,location,is_actived,is_removed,created_at,created_by,updated_at,updated_by,deleted_at,deleted_by) VALUES
	 (1,'Jakarta','Jakarta',0,0,'2024-10-03 22:41:54.000',NULL,'2024-10-04 21:37:43.000',NULL,NULL,NULL),
	 (1,'Bandung','Bandung',1,0,'2024-10-03 22:43:58.000',NULL,'2024-10-03 22:43:58.000',NULL,NULL,NULL),
	 (1,'Yogyakarta','Yogyakarta',1,0,'2024-10-03 22:44:10.000',NULL,'2024-10-03 22:44:10.000',NULL,NULL,NULL);

INSERT INTO ecommerce.product (store_id,warehouse_id,name,price,is_actived,is_removed,created_at,created_by,updated_at,updated_by,deleted_at,deleted_by) VALUES
	 (1,2,'Produk 1',10000.00,1,0,'2024-10-03 22:42:04.0',NULL,'2024-10-05 03:49:42.0',NULL,NULL,NULL);

INSERT INTO ecommerce.product_stock (product_id,stock_reserved,stock_on_hand) VALUES
	 (1,0,100);

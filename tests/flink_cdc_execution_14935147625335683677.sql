EXECUTE STATEMENT SET
BEGIN
	INSERT INTO mysql_cdc_test_blogs SELECT * FROM databend_test_blogs
	INSERT INTO mysql_cdc_test_users SELECT * FROM databend_test_users
END
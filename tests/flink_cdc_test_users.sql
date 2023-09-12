
CREATE TABLE mysql_cdc_test_users (
	id BIGINT,
	user_id INTEGER,
	name VARCHAR,
	created_at BIGINT,
	updated_at BIGINT,
	PRIMARY KEY (id) NOT ENFORCED
) WITH (
	'database-name' = 'test',
	'table-name' = 'users',
	'connector' = 'mysql-cdc'
);

CREATE TABLE databend_test_users (
	id BIGINT,
	user_id INTEGER,
	name VARCHAR,
	created_at BIGINT,
	updated_at BIGINT,
	PRIMARY KEY (id) NOT ENFORCED
) WITH (
	'database-name' = 'test',
	'table-name' = 'users',
	'connector' = 'databend'
);

INSERT INTO databend_test_users SELECT * FROM mysql_cdc_test_users;
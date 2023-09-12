
CREATE TABLE mysql_cdc_test_users (
	id BIGINT,
	user_id INTEGER,
	name VARCHAR,
	created_at BIGINT,
	updated_at BIGINT,
	PRIMARY KEY (id,user_id) NOT ENFORCED
) WITH (
	'connector' = 'mysql-cdc'
)

CREATE TABLE databend_test_users (
	id BIGINT,
	user_id INTEGER,
	name VARCHAR,
	created_at BIGINT,
	updated_at BIGINT,
	PRIMARY KEY (id,user_id) NOT ENFORCED
) WITH (
	'connector' = 'databend'
)
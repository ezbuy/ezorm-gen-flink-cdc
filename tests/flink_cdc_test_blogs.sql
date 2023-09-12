
CREATE TABLE mysql_cdc_test_blogs (
	id BIGINT,
	user_id INTEGER,
	title VARCHAR,
	content VARCHAR,
	status INTEGER,
	readed INTEGER,
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	PRIMARY KEY (id,user_id) NOT ENFORCED
) WITH (
	'connector' = 'mysql-cdc'
)

CREATE TABLE databend_test_blogs (
	id BIGINT,
	user_id INTEGER,
	title VARCHAR,
	content VARCHAR,
	status INTEGER,
	readed INTEGER,
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	PRIMARY KEY (id,user_id) NOT ENFORCED
) WITH (
	'connector' = 'databend'
)
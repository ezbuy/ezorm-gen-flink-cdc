gen-test:
	ezorm gen -i ./tests/blog.yaml -o ./tests --goPackage tests --plugin flink-cdc --plugin-only --plugin-args 'from.connector=mysql-cdc,to.connector=databend'

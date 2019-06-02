

test:
	go test ./... -cover -v

schema:
	mysql -uroot -p < schema.sql

conn_db:
	mysql -h ${DB_HOSTNAME} -u ${DB_USERNAME} -p ${DB_SCHEMA}
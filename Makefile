

test:
	go test -cover -v

# will prompt for password
database:
	mysql -uroot -p < schema.sql

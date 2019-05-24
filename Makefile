

test:
	go test .

# will prompt for password
localdb:
	mysql -uroot -p < schema.sql

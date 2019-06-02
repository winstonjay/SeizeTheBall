

test:
	go test ./... -cover -v

# will prompt for password
database:
	mysql -uroot -p < model/schema.sql

# SELECT * FROM seizetheball.possession;
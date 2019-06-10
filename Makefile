#

USER:="karlsims_uk"
ZONE:="europe-west2-c"
INSTANCE="seizetheball-instance-1"

twlistener:
	GOOS=linux go build -o twlistener listener/main.go

clean:
	rm -f twlistener

test:
	go test ./... -cover -v

test_schema:
	mysql -uroot -p < schema.sql

shh_db:
	mysql -h $(DB_HOSTNAME) -u $(DB_USERNAME) -p $(DB_SCHEMA)
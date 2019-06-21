
USER:=karlsims_uk
ZONE:=europe-west2-c
INSTANCE:=seizetheball-instance-1

# --- listener ---

twlistener:
	GOOS=linux go build -o twlistener listener/main.go

deploy: clean twlistener
	gcloud compute scp --zone $(ZONE) twlistener twlistener.service $(USER)@$(INSTANCE):~
	gcloud compute ssh --zone $(ZONE) $(USER)@$(INSTANCE) --command \
		"sudo mv ~/twlistener.service /etc/systemd/system/"
	gcloud compute ssh --zone $(ZONE) $(USER)@$(INSTANCE) --command \
		"sudo systemctl enable twlistener && sudo systemctl start twlistener"

pingtest:
	GOOS=linux go build -o pingtest ping/main.go

deploy-pingtest: pingtest
	gcloud compute scp --zone $(ZONE) pingtest $(USER)@$(INSTANCE):~
	gcloud compute ssh --zone $(ZONE) $(USER)@$(INSTANCE) --command "./pingtest"

clean:
	rm -f twlistener
	rm -f pingtest

test:
	go test ./... -cover -v

test-schema:
	mysql -uroot -p < schema.sql
# --- frontend ---
PROJECT:=seizetheball

deploy:
	gcloud app deploy --project=$(PROJECT)

serve:
	dev_appserver.py app.yaml

browse:
	gcloud app browse --project=$(PROJECT)

create:
	gcloud app create --project=$(PROJECT) --region=$(ZONE)

viewlogs:
	gcloud app logs tail -s default
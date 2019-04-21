BINARY = slack-status
IMAGENAME = gcr.io/laputa/slack-status
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

include .env
export $(shell sed 's/=.*//' .env)

.PHONY: test clean gcloud-build build deploy

test: $(SRC)
	go test

$(BINARY): $(SRC)
	go build -o $(BINARY)

clean:
	$(RM) $(BINARY)

gcloud-build: $(SRC)
	gcloud builds submit --tag $(IMAGENAME)

build: test $(BINARY)

deploy: gcloud-build
	gcloud beta run \
		deploy slack-status \
		--region us-central1 \
		--set-env-vars SLACK_API_TOKEN="${SLACK_API_TOKEN}",ACCUWEATHER_API_KEY="${ACCUWEATHER_API_KEY}" \
		--image $(IMAGENAME)
BINARY = slack-status
IMAGENAME = laputa/slack-status
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

include secrets

.PHONY: clean docker-image build deploy

$(BINARY): $(SRC)
	GOOS=linux GOARCH=amd64 go build -o $(BINARY)

clean:
	$(RM) $(BINARY)

docker-image: $(BINARY)
	docker build -t $(IMAGENAME) .

build: docker-image

deploy: build
	docker push $(IMAGENAME)
	-ssh vultr docker stop $(BINARY)
	-ssh vultr docker rm $(BINARY)
	ssh vultr docker pull $(IMAGENAME)
	ssh vultr docker run -d --restart=always --name $(BINARY) \
		-e WUNDERGROUND_API_KEY="$(WUNDERGROUND_API_KEY)" \
		-e SLACK_API_TOKEN="$(SLACK_API_TOKEN)" \
		$(IMAGENAME)

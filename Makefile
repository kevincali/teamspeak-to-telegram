NAME := $(shell basename $(CURDIR))
IMAGE=ghcr.io/kevincali/$(NAME)
COMMIT_HASH := $(shell git rev-parse --short HEAD)

build-image:
	CGO_ENABLED=0 go build .
	docker build --tag $(IMAGE):$(COMMIT_HASH) --tag $(IMAGE):latest .

run-image: build-image
	docker run --volume ./config.yaml:/config.yaml --env CONFIG_PATH=/config.yaml $(IMAGE):latest

push-image:
	docker push $(IMAGE):$(COMMIT_HASH)
	docker push $(IMAGE):latest

build:
	go build .

run:
	go run .

ts3:
	docker run \
		-p 9987:9987/udp \
		-p 30033:30033 \
		-p 10011:10011 \
		-e TS3SERVER_LICENSE=accept \
		teamspeak

ts6:
	docker run \
		-p 9987:9987/udp \
		-p 30033:30033 \
		-p 10080:10080 \
		-e TSSERVER_LICENSE_ACCEPTED=accept \
		-e TSSERVER_QUERY_HTTP_ENABLED=1 \
		teamspeaksystems/teamspeak6-server

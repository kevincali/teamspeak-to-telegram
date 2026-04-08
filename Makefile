NAME := $(shell basename $(CURDIR))
IMAGE=ghcr.io/kevincali/$(NAME)
COMMIT_HASH := $(shell git rev-parse --short HEAD)
PLATFORMS=linux/amd64,linux/arm64

build-image:
	docker buildx build --load \
		--tag $(IMAGE):$(COMMIT_HASH) \
		--tag $(IMAGE):latest .

run-image: build-image
	docker run $(IMAGE):latest

push-image:
	docker buildx build --platform $(PLATFORMS) --push \
		--tag $(IMAGE):$(COMMIT_HASH) \
		--tag $(IMAGE):latest .

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

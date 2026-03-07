COMMIT_HASH := $(shell git --no-pager describe --tags --always --dirty)
build-image:
	CGO_ENABLED=0 go build .
	docker build --tag kevincali/teamspeak-to-telegram:$(COMMIT_HASH) --tag kevincali/teamspeak-to-telegram:latest .

run-image: build-image
	docker run --volume ./config.yaml:/config.yaml --env CONFIG_PATH=/config.yaml kevincali/teamspeak-to-telegram:latest

push-image:
	docker push kevincali/teamspeak-to-telegram:$(COMMIT_HASH)
	docker push kevincali/teamspeak-to-telegram:latest

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

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

run: build
	CONFIG_PATH=config.yaml ./teamspeak-to-telegram

ts:
	docker run -p 9987:9987/udp -p 10011:10011 -p 30033:30033 -e TS3SERVER_LICENSE=accept teamspeak


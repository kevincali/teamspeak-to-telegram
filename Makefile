COMMIT_HASH := $(shell git --no-pager describe --tags --always --dirty)
build-image:
	docker build --tag teamspeak-to-telegram:$(COMMIT_HASH) --tag teamspeak-to-telegram:latest .

run-image: build-image
	 docker run --env CONFIG_PATH=config.yaml teamspeak-to-telegram:latest

build:
	go build .

run: build
	CONFIG_PATH=config.yaml ./teamspeak-to-telegram

ts:
	docker run -p 9987:9987/udp -p 10011:10011 -p 30033:30033 -e TS3SERVER_LICENSE=accept teamspeak


.PHONY: build run stop remove

build:
	docker build -t go-tele-bot .

run:
	docker run -d --name go-tele-bot go-tele-bot

stop:
	docker stop go-tele-bot

remove:
	docker rm go-tele-bot

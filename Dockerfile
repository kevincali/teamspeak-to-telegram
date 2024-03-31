FROM golang:alpine

COPY . .

RUN go build -o /teamspeak-to-telegram

ENTRYPOINT ["/teamspeak-to-telegram"]

FROM gcr.io/distroless/static:nonroot

COPY ./teamspeak-to-telegram /bin/teamspeak-to-telegram

ENTRYPOINT ["teamspeak-to-telegram"]

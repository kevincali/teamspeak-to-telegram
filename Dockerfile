FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .

ARG TARGETOS TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/teamspeak-to-telegram .

FROM gcr.io/distroless/static:nonroot

WORKDIR /home/nonroot

COPY --from=builder /bin/teamspeak-to-telegram /bin/teamspeak-to-telegram

ENTRYPOINT ["teamspeak-to-telegram"]

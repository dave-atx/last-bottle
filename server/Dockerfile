ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app .


FROM gcr.io/distroless/base-debian12:nonroot

COPY --from=builder /run-app /usr/local/bin/
COPY --from=builder /usr/src/app/static/ /srv/http/
CMD ["run-app"]

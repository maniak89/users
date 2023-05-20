FROM golang:1.20 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go install ./...
RUN CGO_ENABLED=0 go get github.com/go-delve/delve/cmd/dlv && \
    CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv && \
    CGO_ENABLED=0 go install -tags 'postgres sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate

FROM alpine:latest

COPY migrations /srv/migrations
COPY --from=build /go/bin /srv


ENTRYPOINT ["/srv/users"]
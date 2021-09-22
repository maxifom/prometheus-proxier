FROM golang:1.17-alpine as build-env
WORKDIR /goman/
ADD . .

RUN CGO_ENABLED=0 GOFLAGS=-mod=vendor go test -v ./cmd/...
RUN CGO_ENABLED=0 \
    GOBIN=/goman/apps/ \
    GOOS=linux \
    GOARCH=amd64 \
    go install -v -a -tags netgo -mod vendor ./cmd/...

FROM alpine
WORKDIR /apps/

RUN apk add --no-cache ca-certificates tzdata
COPY --from=build-env /goman/apps/ /apps/

EXPOSE 8080
ENTRYPOINT "/apps/proxy"

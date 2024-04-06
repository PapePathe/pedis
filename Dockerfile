# syntax=docker/dockerfile:1

# Build the pedislication from source
FROM golang:1.21.2 AS build-stage

WORKDIR /pedis

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /pedis

FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM busybox AS build-release-stage
RUN mkdir -p /tmp/

WORKDIR /

COPY --from=build-stage /pedis/pedis /pedis

EXPOSE 6379
EXPOSE 12380
EXPOSE 1237

# USER nonroot:nonroot

ENTRYPOINT /pedis -id $ID -pedis $PEDIS -cluster $CLUSTER -port $PORT

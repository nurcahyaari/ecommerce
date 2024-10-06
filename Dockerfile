# Use the specified Go version
ARG GO_VERSION=1.22.3

# Builder
FROM golang:${GO_VERSION}-alpine as builder

RUN apk update && \
    apk --update add git make build-base

WORKDIR /app

COPY . .

RUN go generate ./...
RUN go build -o main .

# Distribution
FROM alpine:latest

RUN apk update && apk --no-cache add ca-certificates && \
    apk --update --no-cache add tzdata

WORKDIR /app 

EXPOSE 3000

COPY --from=builder /app/main /app

CMD /app/main

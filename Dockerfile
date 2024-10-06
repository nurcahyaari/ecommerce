# Use the specified Go version
ARG GO_VERSION=1.22.3
# Builder stage
FROM golang:${GO_VERSION}-alpine as builder

# Install dependencies
RUN apk update && apk add --no-cache git make build-base

# Set the working directory
WORKDIR /app

COPY . .

RUN GOFLAGS="-buildvcs=false" go build -o goBinary ./main.go

# Final stage
FROM alpine:latest

# Install necessary packages for the final image
RUN apk update && apk add --no-cache ca-certificates tzdata

# Set the working directory in the final image
WORKDIR /app 

# Expose the application port
EXPOSE 9090

# Copy the built binary from the builder stage
COPY --from=builder /app/goBinary /app/

# Command to run the application
CMD ["/app/goBinary"]

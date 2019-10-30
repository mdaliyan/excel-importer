FROM golang:1.12 as builder

# Add Maintainer Info
LABEL maintainer="mohammad alian <md.aliyan@gmail.com>"

# Copy .gitconfig

# Set the Current Working Directory inside the container
WORKDIR /app/

# Copy go mod and sum files
# Because of how the layer caching system works in Docker, the go mod download
# command will _ only_ be re-run when the go.mod or go.sum file change
# (or when we add another docker instruction this line)
COPY ./go.mod .
COPY ./go.sum .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY . .

# Build the Go app
RUN GOOS=linux GOARCH=amd64 go build -o main ./cmd

######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main ./

# Copy the env file for working with binary file
COPY ./env-docker.yml ./env.yml

# Command to run the executable
#CMD ["/main"]

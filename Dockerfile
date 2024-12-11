FROM golang:1.21.1-alpine as builder

ARG GO_BUILD_COMMAND="go build -tags static_all -o movie-fest-skilltest"

# Install build dependencies
RUN apk update && apk --no-cache add build-base git bash coreutils openssh openssl curl

# Create the directory where the application will reside
RUN mkdir -p /go/src/movie-fest-skilltest

WORKDIR /go/src/movie-fest-skilltest

# Copy the main application files
COPY . .

# Application builder step
RUN go mod tidy && go mod download && go mod vendor
RUN eval $GO_BUILD_COMMAND

FROM alpine:3.18.0

# Setup package dependencies
RUN apk --no-cache update && apk --no-cache  add  \
    openssh-client \
    ca-certificates  \
    bash  \
    jq  \
    curl  \
    git

# Create project directory
ENV PROJECT_DIR=/opt/movie-fest-skilltest
RUN mkdir -p $PROJECT_DIR/database/migrations

WORKDIR $PROJECT_DIR

# Copy the built application
COPY --from=builder /go/src/movie-fest-skilltest/movie-fest-skilltest movie-fest-skilltest
COPY --from=builder /go/src/movie-fest-skilltest/database/migration $PROJECT_DIR/database/migration

CMD ["sh", "-c", "/opt/movie-fest-skilltest/movie-fest-skilltest all"]
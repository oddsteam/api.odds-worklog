# Base Image
FROM golang:1.11-alpine

# Install Git
RUN apk update && apk upgrade && \
    apk add --no-cache git

# Download and install the latest release of dep for manage dependency
ADD https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

# Define work directory 
WORKDIR $GOPATH/src/gitlab.odds.team/worklog/api.odds-worklog

# Install dependency
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

# Copy code from host to docker
COPY . ./

# Run Test
RUN CGO_ENABLED=0 GOOS=linux go test ./...

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /go/bin/api

# Start API
ENTRYPOINT ["/go/bin/api"]
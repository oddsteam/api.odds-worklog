FROM golang:1.11-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache git

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/gitlab.odds.team/worklog/api.odds-worklog
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /go/bin/api

ENTRYPOINT ["/go/bin/api"]
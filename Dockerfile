# Build State
# Base Image
FROM golang:1.11-alpine AS build-state

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
RUN ls

# Run Test
RUN CGO_ENABLED=0 GOOS=linux go test ./...

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /go/bin/api


# Deploy State
FROM alpine
RUN apk update && apk upgrade
RUN apk add curl
RUN apk add python
WORKDIR /app
COPY --from=build-state /go/bin/api /app
ADD .env  /app
RUN mkdir -p files/tavi50 && mkdir image && mkdir font
ADD image /app/image
ADD font /app/font

# Start API
ENTRYPOINT ["/app/api"]

# crontab container
FROM alpine:latest

RUN apk add curl && \
    apk add python

RUN mkdir /app

WORKDIR /app

ADD ./callApi.sh /app
ADD ./updateCrontab.sh /app

RUN /bin/sh updateCrontab.sh

CMD crond -l 2 -f && /bin/sh

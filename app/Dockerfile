FROM golang:1.11-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache git
WORKDIR /go/src/api.odds-worklog/app
COPY . .

RUN go get -d -v

RUN pwd
#build the binary
RUN go build -o /go/bin/api
RUN ls
ENTRYPOINT ["/go/bin/api"]

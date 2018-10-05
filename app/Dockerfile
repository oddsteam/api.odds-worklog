FROM golang:1.11-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache git
WORKDIR /go/src/api.odds-worklog/app
COPY . .

RUN go get -d -v

#build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/api

ENTRYPOINT ["/go/bin/api"]

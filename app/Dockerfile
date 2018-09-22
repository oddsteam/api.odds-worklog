FROM golang:1.11-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache git

# RUN go get -u github.com/golang/dep/cmd/dep && \
#     dep init

WORKDIR /go/src/app
COPY . .

#get dependancies
RUN go get -d -v

#build the binary
RUN go build -o /go/bin/api

ENTRYPOINT ["/go/bin/api"]
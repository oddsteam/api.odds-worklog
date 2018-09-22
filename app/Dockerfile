FROM golang:1.11-alpine

WORKDIR /go/src/app
COPY . .

#get dependancies
RUN go get -d -v

#build the binary
RUN go build -o /go/bin/api

ENTRYPOINT ["/go/bin/api"]
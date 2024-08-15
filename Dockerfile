# Build State
# Base Image
FROM golang:1.22-alpine AS build-state

# Install Git
RUN apk update && apk upgrade && \
    apk add --no-cache git

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

# Run Test
RUN CGO_ENABLED=0 GOOS=linux go test ./...

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /go/bin/api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /go/bin/get_student_loan scripts/get_student_loan.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /go/bin/subscribe_friendslog scripts/subscribe_friendslog/main.go


# Deploy State
FROM alpine
RUN apk update && apk upgrade
RUN apk add curl
RUN apk add python3
WORKDIR /app
RUN mkdir -p files/tavi50 && mkdir image && mkdir font
ADD image /app/image
ADD font /app/font
COPY --from=build-state /go/bin/api /app
COPY --from=build-state /go/bin/get_student_loan /app
COPY --from=build-state /go/bin/subscribe_friendslog /app

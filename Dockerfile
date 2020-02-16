# Build executable binary
FROM golang:1.13.0-alpine AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ./build/barometer .


# Build minimal image with entrypoint
FROM alpine

WORKDIR /app

COPY --from=builder /app/build/barometer . 
COPY --from=builder /app/docker-entrypoint.sh .

# Setup environment variables
ENV PORT=8080  \
    BB_PROJECT_ID= \
    BB_TABLE= \
    BB_SLACK_TOKEN= \
    BB_AREA= \
    BB_TWITTER_CONSUMER_KEY= \
    BB_TWITTER_CONSUMER_SECRET= \
    BB_TWITTER_ACCESS_KEY= \
    BB_TWITTER_ACCESS_SECRET= 

ENTRYPOINT ["./docker-entrypoint.sh"]

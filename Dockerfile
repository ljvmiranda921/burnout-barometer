# STEP 1: Build executable binary
FROM golang:1.13.0-alpine AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ./build/barometer .


# STEP 2: Build minimal image with entrypoint
FROM alpine
COPY --from=builder /app/build/barometer /app

WORKDIR /app

ENV PORT=8080 
ENV BB_PROJECT_ID= 
ENV BB_TABLE=
ENV BB_SLACK_TOKEN=
ENV BB_AREA= 

ENTRYPOINT ["./docker-entrypoint.sh"]

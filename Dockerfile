# STEP 0: Declare global arguments 
ARG PORT=8080
ARG BB_PROJECT_ID
ARG BB_TABLE
ARG BB_SLACK_TOKEN
ARG BB_AREA

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
FROM scratch
COPY --from=builder /app/build/barometer /app

WORKDIR /app

ENV PORT $PORT
ENV BB_PROJECT_ID $BB_PROJECT_ID
ENV BB_TABLE $BB_TABLE
ENV BB_SLACK_TOKEN $BB_SLACK_TOKEN
ENV BB_AREA $BB_AREA

RUN /app/build/barometer init --use-env-vars
CMD ["/app/barometer", "serve", "--port=${PORT}"]

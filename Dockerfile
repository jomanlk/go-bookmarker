FROM golang:1.24-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .

RUN go build -o bookmarker ./cmd/bookmarker

FROM alpine:latest

RUN apk --no-cache add ca-certificates
RUN apk add postgresql-client
WORKDIR /app
COPY --from=builder /app/bookmarker /app/bookmarker
RUN mkdir -p /app/data/backup

# Expose the port that the application listens on
EXPOSE 8080

# Command to run the application
CMD ["./bookmarker", "start-server"]
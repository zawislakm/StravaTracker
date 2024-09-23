FROM golang:1.22 AS builder

WORKDIR /app

# Copy the source code into the container
COPY ./src /app/src
COPY go.mod go.sum main.go Makefile ./

RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@latest

RUN make build-prod

FROM ubuntu:latest

WORKDIR /app

RUN apt-get update && apt-get install -y curl

COPY --from=builder /app/tmp .

COPY .env .

EXPOSE 8080

CMD ["./main"]
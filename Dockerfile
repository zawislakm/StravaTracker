FROM golang:1.22 AS go-base

FROM go-base AS templ-builder

WORKDIR /app

COPY go.mod .

RUN go install github.com/a-h/templ/cmd/templ@$(go list -m -f '{{ .Version }}' github.com/a-h/templ)

COPY src/Templates /app/src/Templates

RUN templ generate --path=/app/src/Templates

FROM go-base AS build-stage

WORKDIR /app

COPY ./src /app/src
COPY go.mod go.sum main.go Makefile ./

COPY --from=templ-builder /app/src/Templates /app/src/Templates

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tmp/main .

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /app

COPY ./static /app/static
COPY .env .

COPY --from=build-stage /app/tmp/main /app/main

EXPOSE 8080

CMD ["/app/main"]
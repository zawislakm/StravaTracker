FROM golang:1.24.1 AS go-base

FROM go-base AS templ-builder

WORKDIR /app

COPY go.mod .

RUN go install github.com/a-h/templ/cmd/templ@$(go list -m -f '{{ .Version }}' github.com/a-h/templ)

COPY ./cmd/web/templates /app/copied-templates

RUN templ generate --path=/app/copied-templates

FROM go-base AS build-stage

WORKDIR /app

COPY ./cmd /app/cmd
COPY ./internal /app/internal
COPY go.mod go.sum ./

COPY --from=templ-builder /app/copied-templates /app/cmd/web/templates

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tmp/main cmd/api/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /app

COPY cmd/web/static /app/cmd/web/static
COPY .env .

COPY --from=build-stage /app/tmp/main /app/main

EXPOSE 8080

CMD ["/app/main"]
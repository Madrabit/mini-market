FROM golang:1.24.3-alpine AS build

WORKDIR /app
COPY . .
RUN go mod tidy || true
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/app /app/app
EXPOSE 8080
ENTRYPOINT ["/app/app"]

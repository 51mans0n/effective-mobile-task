FROM golang:1.23-bookworm AS builder
WORKDIR /app

COPY go.* ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o app ./cmd/server

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /app/app .
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["./app"]

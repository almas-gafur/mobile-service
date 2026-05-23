FROM golang:1.21-alpine AS builder

WORKDIR /src
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/repair-crm ./cmd/app

FROM alpine:3.20

RUN apk add --no-cache ca-certificates && adduser -D -H appuser
WORKDIR /app

COPY --from=builder /out/repair-crm /app/repair-crm

USER appuser
EXPOSE 8080

ENTRYPOINT ["/app/repair-crm"]

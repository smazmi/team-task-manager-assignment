FROM golang:1.26.2-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
	&& addgroup -S appgroup \
	&& adduser -S appuser -G appgroup

COPY --from=builder /out/server /app/server

USER appuser

EXPOSE 8080

CMD ["/app/server"]

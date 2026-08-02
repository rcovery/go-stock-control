  FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -tags musl -o /out/app ./cmd/api/main.go

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates libgcc

COPY --from=builder /out/app /app/app

EXPOSE 9000

ENTRYPOINT ["./app"]

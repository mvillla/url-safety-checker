FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o urlinfo ./cmd/urlinfo

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/urlinfo .
COPY data/ data/

EXPOSE 8080

CMD ["./urlinfo"]

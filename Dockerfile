FROM golang:1.25.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ssi-calendar .

FROM alpine:latest
WORKDIR /app
RUN mkdir -p /app/data
COPY --from=builder /app/ssi-calendar .
RUN addgroup -S ssical && adduser -S ssical -G ssical && chown -R ssical:ssical /app
USER ssical:ssical
CMD ["./ssi-calendar"]

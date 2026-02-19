FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY ["../go.mod", "../go.sum", "./"]
RUN go mod download

COPY . .

RUN go build -o receiver ./cmd/receiver

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/receiver /app/receiver

EXPOSE 8080

ENV ANALYZER_URL=http://analyzer:8081

ENTRYPOINT ["/app/receiver"]


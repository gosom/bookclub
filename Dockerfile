FROM golang:1.21.1-alpine as builder

WORKDIR /app

COPY go.mod ./
RUN go mod download
COPY . .

RUN go build -o /usr/bin/rest-api ./cmd/rest-api

FROM scratch

COPY --from=builder /usr/bin/rest-api /rest-api

ENTRYPOINT ["/rest-api"]


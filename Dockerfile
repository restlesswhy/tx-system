FROM golang:1.18 AS builder
WORKDIR /go/src/github.com/restlesswhy/tx-system/

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/main.go

FROM alpine:latest AS app
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /go/src/github.com/restlesswhy/tx-system/app ./
ENTRYPOINT [ "/app/app" ]  
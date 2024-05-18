# First stage: Build the binary
FROM golang:1.21.10 AS builder

LABEL stage=builder
WORKDIR /build

COPY ./pkg ./pkg
COPY ./go.mod .
COPY ./go.sum .
COPY ./main.go .
RUN go mod download && CGO_ENABLED=0 go build -trimpath -o ./server ./main.go

# Second stage: Copy the binary from the builder stage and run it
FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/server .
COPY ./configs/container.json ./configs/container.json

ENV CONF_PATH="/app/configs/container.yaml"

EXPOSE 8168
CMD ["./server"]
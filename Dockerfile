FROM golang:1.24-bookworm AS builder

WORKDIR /build

COPY go.* ./
RUN go mod download
RUN go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o super-cp

FROM buildpack-deps:bookworm-curl

COPY --from=builder /build/super-cp /usr/local/bin/super-cp

ENTRYPOINT ["/usr/local/bin/super-cp"]

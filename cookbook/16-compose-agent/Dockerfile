FROM golang:1.24.0-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY main.go .

RUN <<EOF
go mod tidy 
go build
EOF

#FROM debian:bookworm-slim
#RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
FROM scratch
WORKDIR /app
COPY --from=builder /app/compose-agent .
ENTRYPOINT ["./compose-agent"]

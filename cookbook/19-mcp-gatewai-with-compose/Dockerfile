FROM golang:1.24.0-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY main.go .

RUN <<EOF
go mod tidy 
go build
EOF


FROM debian:12-slim
# ------------------------------------
# Install Socat to use MCP Toolkit
# ------------------------------------
RUN <<EOF
apt-get update
apt-get install -y socat
apt-get clean
rm -rf /var/lib/apt/lists/*
EOF
#FROM scratch
WORKDIR /app
COPY --from=builder /app/tool-calls-mcp-stdio .
ENTRYPOINT ["./tool-calls-mcp-stdio"]

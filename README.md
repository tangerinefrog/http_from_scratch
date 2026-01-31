# HTTP Server from scratch

A minimal HTTP/1.1 server implementation in Go built without using `net/http` or any web frameworks.

## What it does

Implements core HTTP server functionality using only TCP sockets:
- Parses HTTP requests (request line, headers and body)
- Routes requests to handlers
- Builds and sends HTTP responses
- Handles concurrent connections
- Supports basic HTTP methods (GET, POST, PUT, DELETE)

## What This Doesn't Do

- HTTPS/TLS
- HTTP/2 or HTTP/3
- Complex routing patterns (regex, params)
- Middleware chains
- Request/response helpers beyond basics

This is a learning project, not production-ready.

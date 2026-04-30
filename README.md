# My App

A simple Go web application that serves a static HTML page.

## Prerequisites

- [Go](https://go.dev/dl/) 1.24 or later

## Running the application

```bash
go run .
```

The server starts on port `8080` by default. Open <http://localhost:8080> in
your browser to see the welcome page.

To use a different port, set the `PORT` environment variable:

```bash
PORT=3000 go run .
```

## Development

### Project structure

```
main.go        - HTTP server entry point
db.go          - SQLite database helpers
main_test.go   - Tests for the server
index.html     - Static HTML page served at /
```

### Running tests

```bash
go test ./...
```

### Building

```bash
go build -o my-app .
./my-app
```

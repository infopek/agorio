# Agorio

A small agar.io clone with a Go server and a browser client.

## Requirements

- Go 1.25.4 or newer

## Run

From the repository root:

```bash
cd server
go run .
```

Then open `http://localhost:8080` in your browser.
The server also serves the WebSocket endpoint at `ws://localhost:8080/ws`.

# Panel UI Assets

This directory contains the built frontend assets that are embedded into the Go binary.

## Development

In development mode, the application serves files from this directory directly.

## Production

In production mode, these files are embedded into the Go binary using `go:embed` directive in `assets.go`.

## Building UI

To rebuild the UI assets from the web source:

```bash
make build-ui
```

This will:
1. Build the frontend from `web/` directory
2. Copy the built assets to `pkg/panel/ui/`
3. The assets will be automatically embedded on next Go build

## Note

These files are committed to the repository so that users can clone and build the project without needing Node.js/Bun installed.

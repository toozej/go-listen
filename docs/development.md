# Development

## Using the Makefile (Recommended)

The Makefile provides the primary method for developing and installing the application locally. It includes all necessary tools and checks:

```bash
# Install all development dependencies and run the full workflow
make local

# Individual commands:
make local-build          # Build the binary locally
make local-test           # Run all tests with coverage
make local-run            # Run the built binary with .env file
make pre-commit           # Run all code quality checks
make local-cover          # View test coverage in browser

# Development iteration (rebuilds and restarts on file changes)
make local-iterate

# Clean up
make clean
```

## Manual Building (Alternative)

```bash
# Build binary
go build -o go-listen .

# Run tests
go test ./...

# Run with development settings
cp configs/development.env .env
./go-listen serve --debug
```

## Project Structure

```
├── cmd/go-listen/          # CLI commands
├── internal/
│   ├── middleware/         # HTTP middleware (security, logging, rate limiting)
│   ├── server/            # HTTP server and handlers
│   ├── services/          # Business logic services
│   │   ├── duplicate/     # Duplicate detection
│   │   ├── playlist/      # Playlist management
│   │   ├── search/        # Fuzzy artist search
│   │   └── spotify/       # Spotify API integration
│   └── types/             # Type definitions and interfaces
├── pkg/
│   ├── config/            # Configuration management
│   └── logging/           # Structured logging
├── docs/                  # Documentation
├── configs/               # Example configurations
└── static/                # Web interface assets
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## Update Golang Version

```bash
make update-golang-version
```
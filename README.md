# go-listen

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/toozej/go-listen)
[![Go Report Card](https://goreportcard.com/badge/github.com/toozej/go-listen)](https://goreportcard.com/report/github.com/toozej/go-listen)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/toozej/go-listen/cicd.yaml)
![Docker Pulls](https://img.shields.io/docker/pulls/toozej/go-listen)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/toozej/go-listen/total)

<img src="img/avatar.png" alt="go-listen avatar" style="background-color: #FFFFFF;" />

A web application that allows users to search for artists and automatically add their top 5 songs to designated "incoming" playlists on Spotify. Built with Go, featuring a responsive web interface, REST API, and comprehensive security features.

## Features

- **Artist Search**: Fuzzy matching to find artists even with typos or variations
- **Automatic Track Addition**: Adds top 5 tracks from found artists to selected playlists
- **Duplicate Detection**: Prevents adding the same artist's tracks multiple times with override option
- **Playlist Management**: Works with playlists in your "Incoming" folder on Spotify
- **Responsive Web Interface**: Works seamlessly on desktop, tablet, and mobile devices
- **Embedded Spotify Player**: Listen to your playlists directly in the web interface
- **REST API**: Programmatic access for automation and integration
- **Security Features**: CSRF protection, rate limiting, input validation, and security headers
- **Comprehensive Logging**: Structured logging with multiple levels and HTTP request tracking

## Quick Start

### Prerequisites

1. **Spotify Developer Account**: Create an app at [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. **Go 1.21+**: For building from source

### Installation

#### Option 1: Install via Make (recommended)

```bash
make install
```

#### Option 2: Download Binary from GitHub releases

```bash
# Download the latest release for your platform
wget https://github.com/toozej/go-listen/releases/latest/download/go-listen-linux-amd64
chmod +x go-listen-linux-amd64
mv go-listen-linux-amd64 go-listen
```

#### Option 3: Build from Source

```bash
git clone https://github.com/toozej/go-listen.git
cd go-listen
go build -o go-listen .
```

### Configuration

1. Copy the sample environment file:
   ```bash
   cp .env.sample .env
   ```

2. Edit `.env` with your Spotify credentials:
   ```bash
   # Required: Add your Spotify API credentials
   SPOTIFY_CLIENT_ID=your_spotify_client_id_here
   SPOTIFY_CLIENT_SECRET=your_spotify_client_secret_here
   
   # Optional: Customize server settings
   SERVER_HOST=localhost
   SERVER_PORT=8080
   ```

3. Start the application:
   ```bash
   ./go-listen serve
   ```

4. Open your browser to `http://localhost:8080`

## Usage

### Web Interface

1. **Select a Playlist**: Choose from your "Incoming" folder playlists using the searchable dropdown
2. **Search for an Artist**: Enter an artist name (fuzzy matching handles typos)
3. **Add Tracks**: Click "Add Artist" to add their top 5 tracks to the selected playlist
4. **Handle Duplicates**: If tracks already exist, you'll get an option to add anyway
5. **Listen**: Use the embedded Spotify player to listen to your updated playlist

### REST API

See REST [API Documentation](docs/api.md)

## Configuration Options

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SPOTIFY_CLIENT_ID` | - | **Required**: Your Spotify app client ID |
| `SPOTIFY_CLIENT_SECRET` | - | **Required**: Your Spotify app client secret |
| `SPOTIFY_REDIRECT_URL` | `http://127.0.0.1:8080/callback` | Spotify OAuth redirect URL |
| `SERVER_HOST` | `localhost` | Server bind address |
| `SERVER_PORT` | `8080` | Server port |
| `SECURITY_RATE_LIMIT_REQUESTS_PER_SECOND` | `10` | Rate limit per IP |
| `SECURITY_RATE_LIMIT_BURST` | `20` | Rate limit burst capacity |
| `LOGGING_LEVEL` | `info` | Log level (debug, info, warn, error) |
| `LOGGING_FORMAT` | `json` | Log format (json, text) |
| `LOGGING_OUTPUT` | `stdout` | Log output (stdout, stderr, file path) |
| `LOGGING_ENABLE_HTTP` | `true` | Enable HTTP request logging |

### Configuration Files

Example configurations are provided in the `configs/` directory:
- `configs/development.env` - Development settings with verbose logging
- `configs/production.env` - Production settings with optimized logging

## Deployment

### Docker

```bash
docker run -d \
  --name go-listen \
  -p 8080:8080 \
  -e SPOTIFY_CLIENT_ID=your_client_id \
  -e SPOTIFY_CLIENT_SECRET=your_client_secret \
  --restart unless-stopped \
  toozej/go-listen:latest
```

### Docker Compose

See [docker-compose.yml](./docker-compose.yml)

### Systemd Service

See `docs/deployment.md` for detailed deployment instructions including systemd service setup, Kubernetes deployment, and reverse proxy configuration.

## Documentation

- [Configuration Guide](docs/configuration.md) - Detailed configuration options and examples
- [Deployment Guide](docs/deployment.md) - Production deployment instructions
- [Development Guide](docs/development.md) - Local development instructions
- [API Documentation](docs/api.md) - REST API reference

## Security

The application includes several security features:
- **CSRF Protection**: Prevents cross-site request forgery attacks
- **Rate Limiting**: Prevents abuse with configurable limits per IP
- **Input Validation**: Sanitizes and validates all user input
- **Security Headers**: Implements security headers (HSTS, CSP, etc.)
- **Structured Logging**: Comprehensive logging for security monitoring

For production deployment, consider:
- Running behind a reverse proxy with HTTPS
- Implementing additional authentication if needed
- Monitoring logs for security events
- Adjusting rate limits based on usage patterns

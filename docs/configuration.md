# Configuration Guide

This document describes how to configure the go-listen application for deployment and operation.

## Configuration Methods

The application supports configuration through:
1. Environment variables (recommended for production)
2. `.env` file (recommended for development)
3. Command-line flags (limited options)

## Environment Variables

### Required Configuration

#### Spotify API Configuration
These are required for the application to function:

```bash
# Spotify API credentials (required)
SPOTIFY_CLIENT_ID=your_spotify_client_id_here
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret_here
SPOTIFY_REDIRECT_URL=http://localhost:8080/callback
```

**How to get Spotify credentials:**
1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Create a new app
3. Copy the Client ID and Client Secret
4. Add `http://localhost:8080/callback` to the Redirect URIs
   - See [Spotify redirect URL documentation](https://developer.spotify.com/documentation/web-api/concepts/redirect_uri) for more info

### Optional Configuration

#### Server Configuration
```bash
# Server settings (optional, defaults shown)
SERVER_HOST=127.0.0.1t          # Server bind address
SERVER_PORT=8080              # Server port
```

#### Security Configuration
```bash
# Rate limiting (optional, defaults shown)
SECURITY_RATE_LIMIT_REQUESTS_PER_SECOND=10  # Requests per second per IP
SECURITY_RATE_LIMIT_BURST=20                # Burst capacity per IP
```

#### Logging Configuration
```bash
# Logging settings (optional, defaults shown)
LOGGING_LEVEL=info            # Log level: debug, info, warn, error
LOGGING_FORMAT=json           # Log format: json, text
LOGGING_OUTPUT=stdout         # Log output: stdout, stderr, file path
LOGGING_ENABLE_HTTP=true      # Enable HTTP request logging
```


## Troubleshooting

### Common Issues

1. **"spotify client ID and secret are required"**
   - Set `SPOTIFY_CLIENT_ID` and `SPOTIFY_CLIENT_SECRET` environment variables
   - Verify credentials are correct in Spotify Developer Dashboard

2. **"Server failed to start: bind: address already in use"**
   - Change `SERVER_PORT` to an available port
   - Check if another service is using the port: `lsof -i :8080`

3. **Rate limiting too aggressive**
   - Increase `SECURITY_RATE_LIMIT_REQUESTS_PER_SECOND`
   - Increase `SECURITY_RATE_LIMIT_BURST` for burst capacity

4. **Too much/little logging**
   - Adjust `LOGGING_LEVEL` (debug, info, warn, error)
   - Set `LOGGING_ENABLE_HTTP=false` to disable HTTP request logging

### Debug Mode

Enable debug logging for troubleshooting:

```bash
# Via environment variable
LOGGING_LEVEL=debug ./go-listen serve

# Via command line flag
./go-listen serve --debug
```

## Security Considerations

### Production Security

1. **Never commit credentials**: Keep `.env` files out of version control
2. **Use environment variables**: Set credentials via environment in production
3. **Restrict network access**: Bind to specific interfaces if needed
4. **Monitor rate limits**: Adjust based on expected usage patterns
5. **Enable HTTPS**: Use a reverse proxy (nginx, Caddy) for TLS termination

### Network Security

The application includes built-in security features:
- CSRF protection for state-changing operations
- Rate limiting per IP address
- Input validation and sanitization
- Security headers (HSTS, CSP, etc.)

For production deployment, consider:
- Running behind a reverse proxy
- Using TLS/HTTPS
- Implementing additional authentication if needed
- Monitoring and alerting on security events

## Performance Tuning

### Rate Limiting

Adjust based on your usage patterns:

```bash
# For high-traffic scenarios
SECURITY_RATE_LIMIT_REQUESTS_PER_SECOND=50
SECURITY_RATE_LIMIT_BURST=100

# For low-traffic scenarios
SECURITY_RATE_LIMIT_REQUESTS_PER_SECOND=5
SECURITY_RATE_LIMIT_BURST=10
```

### Logging

For high-traffic production:

```bash
# Reduce logging overhead
LOGGING_LEVEL=warn
LOGGING_ENABLE_HTTP=false
```

For development and debugging:

```bash
# Verbose logging
LOGGING_LEVEL=debug
LOGGING_ENABLE_HTTP=true
```
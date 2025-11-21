# go-listen REST API Documentation

This document describes the REST API endpoints available in the go-listen application.

## Base URL

All API endpoints are relative to your server's base URL:
```
http://localhost:8080
```

## Authentication

The go-listen API does not require authentication for internal network usage. All endpoints are publicly accessible within your network.

## Content Type

All API endpoints expect and return JSON data with the content type:
```
Content-Type: application/json
```

## Rate Limiting

The API implements rate limiting to prevent abuse:
- **Default Limit**: 10 requests per second per IP address
- **Burst Capacity**: 20 requests
- **Response**: HTTP 429 (Too Many Requests) when limit exceeded

## Security Features

The API includes several security protections:
- **CSRF Protection**: Required for state-changing operations
- **Input Validation**: All inputs are validated and sanitized
- **Security Headers**: Standard security headers are applied
- **Rate Limiting**: Per-IP rate limiting prevents abuse

## Error Responses

All endpoints return errors in a consistent format:

```json
{
  "success": false,
  "error": "Error message describing what went wrong"
}
```

Common HTTP status codes:
- `200` - Success
- `400` - Bad Request (invalid input)
- `405` - Method Not Allowed
- `429` - Too Many Requests (rate limited)
- `500` - Internal Server Error

## Endpoints

### 1. Get CSRF Token

Get a CSRF token required for state-changing operations.

**Endpoint:** `GET /api/csrf-token`

**Response:**
```json
{
  "csrf_token": "generated-csrf-token-string"
}
```

**Example:**
```bash
curl -X GET http://localhost:8080/api/csrf-token
```

---

### 2. Get Playlists

Retrieve playlists from the "Incoming" folder with optional search filtering.

**Endpoint:** `GET /api/playlists`

**Query Parameters:**
- `search` (optional): Filter playlists by name containing this term

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "playlist_id_string",
      "name": "Playlist Name",
      "uri": "spotify:playlist:playlist_id",
      "track_count": 25,
      "embed_url": "https://open.spotify.com/embed/playlist/playlist_id",
      "is_incoming": true
    }
  ]
}
```

**Examples:**

Get all incoming playlists:
```bash
curl -X GET http://localhost:8080/api/playlists
```

Search for playlists containing "rock":
```bash
curl -X GET "http://localhost:8080/api/playlists?search=rock"
```

---

### 3. Add Artist to Playlist

Add an artist's top 5 tracks to a specified playlist with duplicate detection.

**Endpoint:** `POST /api/add-artist`

**Request Headers:**
```
Content-Type: application/json
X-CSRF-Token: your-csrf-token (required)
```

**Request Body:**
```json
{
  "artist_name": "Artist Name",
  "playlist_id": "spotify_playlist_id",
  "force": false
}
```

**Request Parameters:**
- `artist_name` (required): Name of the artist to search for (1-100 characters)
- `playlist_id` (required): Spotify playlist ID where tracks should be added
- `force` (optional): Set to `true` to bypass duplicate detection (default: `false`)

**Success Response:**
```json
{
  "success": true,
  "message": "Successfully added 5 tracks from Artist Name to Playlist Name",
  "data": {
    "success": true,
    "artist": {
      "id": "artist_spotify_id",
      "name": "Artist Name",
      "uri": "spotify:artist:artist_id",
      "genres": ["rock", "alternative"]
    },
    "tracks_added": [
      {
        "id": "track_id",
        "name": "Track Name",
        "uri": "spotify:track:track_id",
        "artists": [
          {
            "id": "artist_id",
            "name": "Artist Name",
            "uri": "spotify:artist:artist_id",
            "genres": []
          }
        ],
        "duration_ms": 240000
      }
    ],
    "playlist": {
      "id": "playlist_id",
      "name": "Playlist Name",
      "uri": "spotify:playlist:playlist_id",
      "track_count": 30,
      "embed_url": "https://open.spotify.com/embed/playlist/playlist_id",
      "is_incoming": true
    },
    "was_duplicate": false,
    "message": "Successfully added 5 tracks from Artist Name to Playlist Name"
  }
}
```

**Duplicate Detection Response:**
When tracks already exist and `force` is `false`:
```json
{
  "success": false,
  "message": "Artist Name's tracks were already added to Playlist Name on 2024-01-15T10:30:00Z. Use force=true to add anyway.",
  "is_duplicate": true,
  "last_added": "2024-01-15T10:30:00Z",
  "data": {
    "success": false,
    "was_duplicate": true,
    "message": "Duplicate tracks detected"
  }
}
```

**Examples:**

Add artist without force (will detect duplicates):
```bash
curl -X POST http://localhost:8080/api/add-artist \
  -H "Content-Type: application/json" \
  -H "X-CSRF-Token: EXAMPLE" \
  -d '{
    "artist_name": "Radiohead",
    "playlist_id": "37i9dQZF1DX0XUsuxWHRQd"
  }' #gitleaks:allow
```

Add artist with force (bypass duplicate detection):
```bash
curl -X POST http://localhost:8080/api/add-artist \
  -H "Content-Type: application/json" \
  -H "X-CSRF-Token: EXAMPLE" \
  -d '{
    "artist_name": "Radiohead",
    "playlist_id": "37i9dQZF1DX0XUsuxWHRQd",
    "force": true
  }'
```

## Data Models

### Artist
```json
{
  "id": "string",           // Spotify artist ID
  "name": "string",         // Artist display name
  "uri": "string",          // Spotify URI
  "genres": ["string"]      // Array of genre strings
}
```

### Track
```json
{
  "id": "string",           // Spotify track ID
  "name": "string",         // Track title
  "uri": "string",          // Spotify URI
  "artists": [Artist],      // Array of artist objects
  "duration_ms": number     // Track duration in milliseconds
}
```

### Playlist
```json
{
  "id": "string",           // Spotify playlist ID
  "name": "string",         // Playlist name
  "uri": "string",          // Spotify URI
  "track_count": number,    // Number of tracks in playlist
  "embed_url": "string",    // Spotify embed URL
  "is_incoming": boolean    // Whether playlist is in "Incoming" folder
}
```

### Add Result
```json
{
  "success": boolean,       // Whether operation succeeded
  "artist": Artist,         // Artist that was processed
  "tracks_added": [Track],  // Array of tracks that were added
  "playlist": Playlist,     // Target playlist
  "was_duplicate": boolean, // Whether duplicates were detected
  "message": "string"       // Human-readable result message
}
```

## Error Handling

### Validation Errors

**Status Code:** `400 Bad Request`

Common validation errors:
- Missing required fields (`artist_name`, `playlist_id`)
- Invalid field lengths (artist name must be 1-100 characters)
- Invalid playlist ID format
- Missing CSRF token for POST requests

### Rate Limiting

**Status Code:** `429 Too Many Requests`

```json
{
  "success": false,
  "error": "Rate limit exceeded. Please try again later."
}
```

### Spotify API Errors

**Status Code:** `500 Internal Server Error`

Common scenarios:
- Artist not found
- Playlist not accessible
- Spotify API rate limits
- Network connectivity issues

```json
{
  "success": false,
  "error": "Failed to add artist: Artist 'Unknown Artist' not found"
}
```

## Usage Examples

### Complete Workflow Example

1. **Get CSRF Token:**
```bash
CSRF_TOKEN=$(curl -s http://localhost:8080/api/csrf-token | jq -r '.csrf_token')
```

2. **Get Available Playlists:**
```bash
curl -s http://localhost:8080/api/playlists | jq '.data[].name'
```

3. **Add Artist to Playlist:**
```bash
curl -X POST http://localhost:8080/api/add-artist \
  -H "Content-Type: application/json" \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -d '{
    "artist_name": "The Beatles",
    "playlist_id": "your_playlist_id_here"
  }' | jq '.'
```

### JavaScript Example

```javascript
// Get CSRF token
const csrfResponse = await fetch('/api/csrf-token');
const { csrf_token } = await csrfResponse.json();

// Add artist to playlist
const response = await fetch('/api/add-artist', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'X-CSRF-Token': csrf_token
  },
  body: JSON.stringify({
    artist_name: 'Pink Floyd',
    playlist_id: 'your_playlist_id',
    force: false
  })
});

const result = await response.json();
console.log(result);
```

## Configuration

API behavior can be configured through environment variables:

- `SECURITY_RATE_LIMIT_REQUESTS_PER_SECOND`: Rate limit per IP (default: 10)
- `SECURITY_RATE_LIMIT_BURST`: Burst capacity (default: 20)
- `SERVER_HOST`: Server bind address (default: localhost)
- `SERVER_PORT`: Server port (default: 8080)

## Troubleshooting

### Common Issues

1. **CSRF Token Missing/Invalid**
   - Ensure you get a fresh CSRF token before making POST requests
   - Include the token in the `X-CSRF-Token` header

2. **Rate Limited**
   - Reduce request frequency
   - Implement exponential backoff in your client

3. **Artist Not Found**
   - Check artist name spelling
   - Try variations of the artist name
   - The fuzzy matching handles some typos but may not catch all variations

4. **Playlist Access Issues**
   - Ensure the playlist exists in your "Incoming" folder
   - Verify Spotify credentials are configured correctly
   - Check that the playlist is not private or restricted

5. **Network Errors**
   - Verify server is running and accessible
   - Check firewall settings
   - Ensure Spotify API credentials are valid

For more detailed troubleshooting, check the server logs which include structured logging for all operations.
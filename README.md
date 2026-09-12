# spotify-stats

A lightweight Go backend that authenticates with the Spotify Web API (OAuth
Authorization Code flow) and exposes simple REST endpoints for a user's
listening data: profile info, top artists/tracks, and recently played tracks.

Built with [Fiber](https://gofiber.io/) (v3) for HTTP routing and a small
custom HTTP client for calling the Spotify API.

## Features

- Spotify OAuth login flow (`/api/web-auth` → Spotify → `/api/callback`)
- Automatic access token caching/refresh
- Fetch user profile, top artists, top tracks, recently played tracks, and
  saved (liked) tracks
- Configurable time range (`short_term` / `medium_term` / `long_term`) and
  result limit on "top" queries
- Derived listening stats built on top of that data:
  - **Taste evolution** — overlap between top artists/tracks across time
    ranges, and which artists are new in the more recent range
  - **Listening patterns** — hour-of-day and day-of-week histograms from
    recently played tracks
  - **Album art colors** — average color of each top track's album artwork

> Note: Spotify deprecated the `audio-features`, `audio-analysis`,
> `recommendations`, and `related-artists` endpoints (and stopped returning
> `genres`/`popularity` on artists and tracks) for apps created after
> November 2024. The stats above intentionally avoid those fields/endpoints.

## Project structure

```
cmd/server            entry point, route wiring
internal/auth         Spotify OAuth flow and token management
internal/spotify-data Spotify API client (profile, top artists/tracks, recent, saved tracks)
internal/stats        derived listening stats (taste evolution, listening patterns, album colors)
internal/httpclient   small wrapper around net/http
docs                  generated Swagger/OpenAPI spec (swag init)
```

## Prerequisites

- Go 1.26+
- A [Spotify Developer](https://developer.spotify.com/dashboard) app

## Setup

1. Create an app on the Spotify Developer Dashboard and add a Redirect URI of:

   ```
   http://127.0.0.1:{PORT}/api/callback
   ```

   where `{PORT}` matches the `PORT` you set below. The URI must be
   registered **exactly** (including the port) under the app's Basic
   Information page, or Spotify will return `error=server_error` during
   login.

2. Copy the example env file and fill in your credentials:

   ```bash
   cp .env.example .env
   ```

   ```env
   CLIENT_ID=your-spotify-client-id
   CLIENT_SECRET=your-spotify-client-secret
   PORT=8080
   ```

3. Run the server:

   ```bash
   make run
   ```

## API documentation (Swagger)

Once the server is running, browse the interactive Swagger UI at:

```
http://127.0.0.1:{PORT}/swagger/index.html
```

The OpenAPI spec is generated from code annotations using [swaggo](https://github.com/swaggo/swag)
and served via [gofiber/contrib/swaggo](https://github.com/gofiber/contrib/tree/main/v3/swaggo).
If you add or change endpoints, regenerate the docs with:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
make swagger
```

## Authentication flow

1. `GET /api/web-auth` — returns the Spotify authorization URL to open in a
   browser. Requests the `user-top-read`, `user-read-recently-played`, and
   `user-library-read` scopes.
2. After the user approves access, Spotify redirects to `/api/callback` with
   an authorization `code`.
3. Subsequent requests to the data endpoints automatically exchange/refresh
   the access token as needed.

## API

| Method | Endpoint                        | Description                                     |
|--------|----------------------------------|--------------------------------------------------|
| GET    | `/api/web-auth`                 | Get the Spotify login/authorize URL               |
| GET    | `/api/callback`                 | OAuth redirect target (used by Spotify)           |
| GET    | `/api/token`                    | Get a valid (cached or refreshed) token           |
| GET    | `/api/me`                       | Current user's Spotify profile                    |
| GET    | `/api/top/artists`              | User's top artists                                |
| GET    | `/api/top/tracks`               | User's top tracks                                 |
| GET    | `/api/recently-played`          | User's recently played tracks                     |
| GET    | `/api/saved-tracks`             | User's saved (liked) tracks                       |
| GET    | `/api/stats/taste-evolution`    | Top artist/track overlap across time ranges       |
| GET    | `/api/stats/listening-patterns` | Hour-of-day / day-of-week listening histograms    |
| GET    | `/api/stats/album-colors`       | Average color of each top track's album artwork   |

`top/*`, `recently-played`, `saved-tracks`, and `album-colors` endpoints
accept `time_range` and/or `limit` query parameters (see `TopOptions`).

# Viewport Forge

A DevTools screenshot automation platform for responsive QA.

Upload a URL and Viewport Forge queues a capture job, runs Playwright across multiple device classes, and stores the generated screenshots for fast visual checks.

## Why this project

- Real-world developer tooling for responsive validation.
- Practical automation with clean service boundaries.
- Great foundation for visual diffing and layout-break detection.

## Core stack

- Frontend: React + TypeScript + Vite
- API: Go (`net/http`) + Redis queue producer
- Worker: Node.js + Playwright + Redis queue consumer
- Queue: Redis list (`vf:capture_jobs`)

## What works right now

- Queue screenshot jobs from the frontend or API.
- Capture default viewport set:
  - iPhone (`390x844`)
  - tablet (`834x1112`)
  - laptop (`1440x900`)
  - ultrawide (`2560x1080`)
  - 4K (`3840x2160`)
- Poll job state (`queued`, `processing`, `completed`, `failed`).
- Write screenshots to `artifacts/<job_id>/`.

## Project layout

```text
viewport-forge/
├── frontend/             # React UI
├── backend/              # Go API
├── worker/               # Playwright worker
├── artifacts/            # Generated screenshots (runtime)
├── docker-compose.yml    # Redis service
└── Makefile
```

## Quick start

### 1) Install dependencies

```bash
make setup
```

### 2) Start Redis

```bash
make redis
```

### 3) Install Playwright Chromium

```bash
cd worker && npm run install:browsers
```

### 4) Run services

In separate terminals:

```bash
make api
make worker
make web
```

Frontend: `http://localhost:5173`  
API health: `http://localhost:8080/health`

## API

### `POST /api/v1/captures`

Request body:

```json
{
  "url": "https://example.com"
}
```

Response (`202 Accepted`):

```json
{
  "id": "e4f2f6f9a2474c11",
  "state": "queued",
  "status_url": "/api/v1/captures/e4f2f6f9a2474c11"
}
```

### `GET /api/v1/captures/:id`

Response (`200 OK`):

```json
{
  "id": "e4f2f6f9a2474c11",
  "url": "https://example.com",
  "state": "completed",
  "requested_at": "2026-03-04T16:00:00Z",
  "started_at": "2026-03-04T16:00:04Z",
  "finished_at": "2026-03-04T16:00:19Z",
  "output_dir": "/.../artifacts/e4f2f6f9a2474c11",
  "screenshots": "5"
}
```

## Near-term roadmap

- Visual diffing against baseline snapshots.
- Layout-break detection rules (overflow, clipping, overlap).
- Lighthouse score collection per viewport.
- Object storage uploads + signed URLs.
- Team/project scoping + auth.

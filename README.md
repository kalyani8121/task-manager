# Task Management System

Production-style REST API in Go with JWT authentication, PostgreSQL, Docker and Kubernetes.

## Run

git clone https://github.com/kalyani8121/task-manager.git
cd task-manager
docker-compose up --build

## API Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | /api/v1/auth/register | No | Create account |
| POST | /api/v1/auth/login | No | Login get token |
| POST | /api/v1/tasks | Yes | Create task |
| GET | /api/v1/tasks | Yes | Get all tasks |
| PUT | /api/v1/tasks/:id | Yes | Update task |
| DELETE | /api/v1/tasks/:id | Yes | Delete task |
| GET | /api/v1/analytics/summary | Yes | Overall stats |
| GET | /api/v1/analytics/by-status | Yes | Count by status |
| GET | /api/v1/analytics/overdue | Yes | Overdue tasks |

## Tech Stack
Go, Gin, PostgreSQL, JWT, Docker, Kubernetes

## Architecture
Handler → Service → Repository → PostgreSQL
# Task Management System

Production-style REST API in Go with JWT authentication, PostgreSQL, Docker and Kubernetes.

## Live URLs
| Service | URL |
|---------|-----|
| Health Check | https://task-manager-production-4677.up.railway.app/health |
| Swagger API Docs | https://task-manager-production-4677.up.railway.app/swagger/index.html |
| Frontend App | https://task-manager-frontend-git-main-kalyani19.vercel.app |

## Unique Features
- Smart Priority Queue — auto-ranks tasks by urgency score
- Analytics Dashboard — completion rates and overdue tracking
- Daily Email Notifications — 9AM IST via Resend API

## Run

git clone https://github.com/kalyani8121/task-manager.git
cd task-manager
docker-compose up --build

## API Endpoints
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | /api/v1/auth/register | No | Register + send verification email |
| GET | /api/v1/auth/verify | No | Verify email with token |
| POST | /api/v1/auth/login | No | Login (requires verified email) |
| POST | /api/v1/tasks | Yes | Create task |
| GET | /api/v1/tasks | Yes | Get all tasks |
| GET | /api/v1/tasks/:id | Yes | Get one task |
| PUT | /api/v1/tasks/:id | Yes | Update task |
| DELETE | /api/v1/tasks/:id | Yes | Delete task |
| GET | /api/v1/tasks/priority-queue | Yes | Smart priority ranking |
| GET | /api/v1/analytics/summary | Yes | Overall stats |
| GET | /api/v1/analytics/by-status | Yes | Count by status |
| GET | /api/v1/analytics/overdue | Yes | Overdue tasks |

## Tech Stack
| Technology | Purpose |
|------------|---------|
| Go 1.22 + Gin | REST API Framework |
| PostgreSQL | Database |
| JWT + bcrypt | Authentication and Security |
| Docker | Containerization (60MB image) |
| Kubernetes | Deployment (2 replicas) |
| Railway | Cloud Hosting |
| Resend API | Email Delivery |
| Swagger | API Documentation |
| Zap | Structured Logging |

## Email Notifications
Every day at 9AM IST:
1. Scans all verified users
2. Gets their pending tasks
3. Calculates priority scores
4. Sends email with ranked task list

## Architecture
Request → Handler → Service → Repository → PostgreSQL

## How to Run Locally
```bash
git clone https://github.com/kalyani8121/task-manager.git
cd task-manager
docker-compose up --build
```
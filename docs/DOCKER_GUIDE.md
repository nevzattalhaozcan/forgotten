# QA Testing Guide - Forgotten App

Simple guide for QA testers to run and test the Forgotten application using Docker.

## Prerequisites

- Docker and Docker Compose installed
- Access to the GitHub repository or container image

## Quick Start (2 Options)

### Option A: Using Pre-built Container (Recommended)

```bash
# 1. Login to GitHub Container Registry
docker login ghcr.io
# Username: your-github-username
# Password: ask developer for access token

# 2. Download and run the application
docker pull ghcr.io/nevzattalhaozcan/forgotten-app:latest
docker run -d --name forgotten-postgres -e POSTGRES_DB=forgotten_db -e POSTGRES_USER=forgotten_user -e POSTGRES_PASSWORD=123456 -p 5432:5432 postgres:15-alpine
docker run -d --name forgotten-app -p 8080:8080 --link forgotten-postgres -e DB_URL=postgres://forgotten_user:123456@forgotten-postgres:5432/forgotten_db ghcr.io/nevzattalhaozcan/forgotten-app:latest
```

### Option B: Using Repository

```bash
# 1. Clone the repository
git clone https://github.com/nevzattalhaozcan/forgotten.git
cd forgotten

# 2. Start all services
make dev
```

## Application Access

Once running, you can access:

- **Main Application**: http://localhost:8080
- **Database Admin (PgAdmin)**: http://localhost:5050
  - Email: `admin@forgotten.com`
  - Password: `admin123`

## Testing the Application

### Health Check
```bash
curl http://localhost:8080/health
```

### Common Test Endpoints
- `GET http://localhost:8080/health` - Application health status
- `GET http://localhost:8080/api/users` - User endpoints (if available)
- Check with developer for specific API endpoints

## Managing the Application

### View Application Logs
```bash
# If using Option A
docker logs forgotten-app

# If using Option B (repository)
make docker-logs
```

### Stop the Application
```bash
# If using Option A
docker stop forgotten-app forgotten-postgres
docker rm forgotten-app forgotten-postgres

# If using Option B (repository)
make docker-down
```

### Restart the Application
```bash
# If using Option A
docker start forgotten-postgres forgotten-app

# If using Option B (repository)
make docker-up
```

## Database Access for Testing

### Using PgAdmin (Web Interface)
1. Go to http://localhost:5050
2. Login with `admin@forgotten.com` / `admin123`
3. Add new server:
   - Name: `Forgotten DB`
   - Host: `postgres` (or `forgotten-postgres` for Option A)
   - Port: `5432`
   - Database: `forgotten_db`
   - Username: `forgotten_user`
   - Password: `123456`

### Using Command Line
```bash
# Connect directly to database
docker exec -it forgotten-postgres psql -U forgotten_user -d forgotten_db
```

## Troubleshooting

### Application Won't Start
```bash
# Check if ports are in use
lsof -i :8080
lsof -i :5432

# View error logs
docker logs forgotten-app
```

### Can't Access Database
```bash
# Ensure database container is running
docker ps | grep postgres

# Check database logs
docker logs forgotten-postgres
```

### Reset Everything
```bash
# Stop and remove all containers
docker stop $(docker ps -a -q --filter="name=forgotten")
docker rm $(docker ps -a -q --filter="name=forgotten")

# Start fresh
# Then repeat Quick Start steps
```

## Test Data

The application starts with a clean database. Ask the developer for:
- Sample data scripts
- Test user credentials
- API documentation
- Postman collections (if available)

## Reporting Issues

When reporting bugs, please include:
1. **Steps to reproduce**
2. **Expected vs actual behavior**
3. **Application logs**: `docker logs forgotten-app`
4. **Database logs** (if relevant): `docker logs forgotten-postgres`
5. **Browser console errors** (for UI issues)

## Quick Commands Reference

| Action | Command |
|--------|---------|
| Start application | `make dev` (Option B) |
| Stop application | `make docker-down` (Option B) |
| View logs | `make docker-logs` (Option B) |
| Reset database | Stop containers → restart |
| Check app health | `curl http://localhost:8080/health` |

---

**Need Help?** Contact the development team with any questions or issues.
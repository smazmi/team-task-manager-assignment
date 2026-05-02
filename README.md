# Team Task Manager

A full-stack collaborative web application where multiple users can create projects, assign tasks, and track their progress. It includes role-based access control (Admin and Member) and a dashboard for monitoring project status.

## Features

- **User Authentication**: Signup and secure login.
- **Project Management**: Create projects, manage team members.
- **Task Management**: Create tasks (Title, Description, Due Date, Priority), assign users, and update statuses.
- **Dashboard**: Track total tasks, tasks by status, tasks per user, and overdue tasks.
- **Role-Based Access**: Admins can manage tasks and users; Members can view and update their assigned tasks.

## Tech Stack

- **Frontend**: React, TypeScript, Vite
- **Backend**: Go (Golang)
- **Database**: PostgreSQL
- **Infrastructure**: Docker, GitHub Actions, GitHub Container Registry (GHCR)

## Local Setup

### Option 1: Using Docker (Recommended)

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd team-task-manager-assignment
   ```

2. Copy the environment file:
   ```bash
   cp .env.example .env.docker
   ```

3. Start the application using Docker Compose:
   ```bash
   docker compose -f docker-compose.dev.yml up --build
   ```

4. The application will be available at:
   - Frontend: `http://localhost:5173`
   - Backend API: `http://localhost:8080`

### Option 2: Manual Setup

**Backend**
1. Ensure Go and PostgreSQL are installed.
2. Navigate to the backend directory:
   ```bash
   cd backend
   ```
3. Copy the environment file and configure your database credentials:
   ```bash
   cp .env.example .env
   ```
4. Run the Go server:
   ```bash
   go run ./cmd/server
   ```

**Frontend**
1. Ensure Bun (or Node.js) is installed.
2. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
3. Install dependencies:
   ```bash
   bun install
   ```
4. Start the development server:
   ```bash
   bun run dev
   ```

## Deployment

The application is configured to be deployed on [Railway](https://railway.app/).

1. **Database**: Create a new PostgreSQL database instance on Railway.
2. **Backend**:
   - Create a new service from your GitHub repository.
   - Point the root directory to `/backend` (or deploy the root directory if Railway detects the Dockerfile automatically).
   - Set the necessary environment variables (`DATABASE_URL`, `JWT_SECRET`, etc.) using the variables from your Railway PostgreSQL instance.
3. **Frontend**:
   - Create another service from your GitHub repository.
   - Point the root directory to `/frontend`.
   - Set `VITE_API_BASE_URL` to the public URL of your deployed backend service.
   - Ensure Railway uses the correct build command (`bun run build`) and output directory (`dist`).

Alternatively, if using the included Docker images from GHCR:
- You can provide the `docker-compose.yml` file to a VPS or Docker environment.
- Fill out a `.env` file on the server with production secrets.
- Run `docker compose up -d` to pull and run the pre-built images. Watchtower will automatically keep them updated.

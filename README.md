# Pulse — Modern Project Management Workspace & Client Portal

A modern, clean, and high-converting UI & Backend Architecture for a web-based Project Management Application inspired by Linear and Notion.

---

## 🌟 Architecture Overview

```
egn-planning/
├── backend/                      # Production Golang RBAC & Project Authorization Server
│   ├── cmd/server/main.go        # Server entrypoint (Port 8080)
│   ├── internal/auth/            # JWT tokens & Central Permission Matrix
│   ├── internal/database/        # SQLite pure-Go driver & auto-seed engine
│   ├── internal/middleware/      # AuthMiddleware, RequireRole, RequirePermission, RequireProjectAccess
│   ├── internal/models/          # User, Project, ProjectMember pivot, Milestone, Task, SubTask, Comment
│   ├── internal/handlers/        # Auth, User, Project, Task, and Milestone REST handlers
│   └── internal/routes/          # Route definitions & CORS configuration
│
├── frontend/                     # Vue 3 Production Workspace (Vite + TS + Pinia + Tailwind)
│   ├── src/types/project.ts      # TypeScript interfaces aligned with Golang REST models
│   ├── src/services/api.ts       # Axios client with JWT Bearer interceptor & mock fallback
│   ├── src/stores/projectStore.ts# Pinia store managing optimistic updates & RBAC UI sync
│   └── src/views/                # Dashboard, Project Workspace, and Client Portal
│
├── UX_BLUEPRINT.md               # Detailed UX Architecture & Information Architecture Blueprint
└── index.html                    # Standalone zero-dependency browser prototype
```

---

## 🔐 Role-Based Access Control (RBAC) & Permissions

The system implements 4 strictly bounded roles:

| Role | Scope | Key Capabilities | Restricted Actions |
|---|---|---|---|
| `super_admin` | Global System | • Full system bypass<br>• Access all projects without membership<br>• Update user roles (`PATCH /users/:id/role`) | None |
| `pm` (Project Manager) | Project Owner | • Create, edit, and archive projects<br>• Manage project milestones<br>• Create, assign, edit, move, delete tasks<br>• Invite/remove members (`/projects/:id/members`) | Cannot update global system roles |
| `member` (Team Member) | Project Contributor | • View assigned projects<br>• Create and move tasks on Kanban<br>• Manage subtask checklists<br>• Upload attachments & post comments | Cannot delete projects<br>Cannot manage project members<br>Cannot approve milestones |
| `client` (Client Stakeholder) | Read-Only & Gating | • View assigned projects & milestones<br>• Inspect client-visible deliverables<br>• Post feedback comments<br>• **Approve or reject milestones** (`PATCH /milestones/:id/approve`) | Cannot create/edit/move tasks<br>Cannot access unassigned projects<br>Cannot view internal-only tasks |

### Project-Level Membership Authorization (`project_members`)
Users (`pm`, `member`, `client`) can **only** access projects they are explicitly assigned to in the `project_members` pivot table. Accessing an unassigned project immediately returns `403 Forbidden`.

---

## 🔑 Test Accounts & Credentials

The backend automatically seeds 4 pre-configured test accounts (1 per role):

> **Default Password for all test accounts:** `Password123!`

- **Super Admin**: `admin@pulse.workspace` / `Password123!`
- **Project Manager**: `pm@pulse.workspace` / `Password123!`
- **Team Member**: `member@pulse.workspace` / `Password123!`
- **Client**: `client@pulse.workspace` / `Password123!` *(Assigned to SUSTAGRI, unassigned to SMART)*

---

## 🚀 Quickstart Guide

### 1. Run the Golang Backend Server
```powershell
cd backend
go run ./cmd/server
```
API server starts at [http://localhost:8080](http://localhost:8080).
- Health Check: `GET http://localhost:8080/health`
- Auth Me: `GET http://localhost:8080/api/v1/auth/me` (requires Bearer token)

### 2. Run the Vue 3 Frontend
```powershell
cd frontend
npm install
npm run dev
```
Open [http://localhost:5173](http://localhost:5173) in your browser.

### 3. Standalone Prototype
Double-click [index.html](file:///c:/Users/agraphies/Documents/GitHub/egn-planning/index.html) to open the interactive UI directly in any browser with zero installation needed.

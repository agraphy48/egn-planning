# Pulse Backend — Role-Based Access Control (RBAC) & Project Authorization Engine

A high-performance Golang backend implementing a strict **Role-Based Access Control (RBAC)** and **Project-Level Membership Authorization** system for the Pulse Project Management Workspace.

---

## 🔐 1. Role & Permission Matrix

The system enforces 4 core roles with strict boundaries:

| Role | Scope | Allowed Actions & Permissions | Restricted / Prohibited Actions |
|---|---|---|---|
| `super_admin` | Global System | • Full system bypass<br>• Access all projects without membership<br>• Manage users & global roles (`PATCH /users/:id/role`)<br>• System configuration | None |
| `pm` (Project Manager) | Project Owner / Manager | • Create/edit/archive projects<br>• Manage project milestones<br>• Create, assign, edit, move, delete tasks<br>• Manage subtasks & upload attachments<br>• Invite/remove `member` and `client` to their projects (`/projects/:id/members`) | Cannot alter system-wide user roles |
| `member` (Team Member) | Assigned Contributor | • View assigned projects<br>• Create tasks within assigned projects<br>• Update status & move cards on Kanban<br>• Manage subtask checklists<br>• Upload attachments & post comments | Cannot delete projects<br>Cannot manage project members<br>Cannot approve milestones |
| `client` (Client Stakeholder) | Read-Only & Gating | • Read-only access to assigned projects & milestones<br>• View client-visible deliverables<br>• Post feedback & comments<br>• **Approve or reject contractual milestones** (`PATCH /milestones/:id/approve`) | Cannot create, edit, or move tasks<br>Cannot access unassigned projects<br>Cannot view internal-only tasks |

---

## 🗄️ 2. Database Schema & Project-Level Authorization

### Pivot Table: `project_members`
Enforces that `pm`, `member`, and `client` users can **only access projects they are explicitly assigned to**:
```sql
CREATE TABLE project_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id VARCHAR(64) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    project_role VARCHAR(32) NOT NULL DEFAULT 'member', -- pm, member, client
    joined_at DATETIME,
    UNIQUE(project_id, user_id)
);
```

- When a request touches `/projects/:id`, `/projects/:id/members`, or `/tasks?project_id=...`, the `RequireProjectAccess()` middleware queries `project_members`.
- If the user is `super_admin`, access is granted globally.
- If the user is not found in `project_members` for the target project, the API immediately returns `403 Forbidden`.

---

## 🛡️ 3. Middleware & Security Layer

Implemented in `internal/middleware/auth.go` and `internal/middleware/rbac.go`:

1. **`AuthMiddleware()`**:
   - Extracts and validates JWT `Bearer <token>` from the HTTP `Authorization` header.
   - Extracts `user_id`, `email`, and `role` into the Gin request context (`c.Set(...)`).
2. **`RequireRole(allowedRoles ...string)`**:
   - Route-level guard verifying user has one of the allowed global roles (e.g. `super_admin`).
3. **`RequirePermission(permission string)`**:
   - Granular permission check (e.g., `project:create`, `task:move`, `milestone:approve`, `user:manage`) evaluated against the central permission registry in `internal/auth/permissions.go`.
4. **`RequireProjectAccess()`**:
   - Inspects `:project_id`, `:id`, or query parameters to verify membership in `project_members`. Returns clean `403 Forbidden` if unauthorized.

---

## 🔑 4. Seeded Test Credentials

The database auto-seeds 4 pre-configured test accounts (1 per role) with the default password:

> **Default Password for all test accounts:** `Password123!`

| Role | Email | Password | Assigned Projects | Key Capabilities |
|---|---|---|---|---|
| **Super Admin** | `admin@pulse.workspace` | `Password123!` | All (Global Bypass) | Change global roles, delete any project |
| **Project Manager** | `pm@pulse.workspace` | `Password123!` | `SUSTAGRI`, `SMART` | Invite members, manage milestones, delete tasks |
| **Team Member** | `member@pulse.workspace` | `Password123!` | `SUSTAGRI`, `SMART` | Move Kanban cards, check subtasks, post comments |
| **Client** | `client@pulse.workspace` | `Password123!` | `SUSTAGRI` only | Approve Milestone 1, view deliverables |

*Notice: The `client@pulse.workspace` account is assigned to `SUSTAGRI`, but NOT `SMART`. Calling `GET /api/v1/projects/proj_smart` with the client token automatically returns `403 Forbidden`.*

---

## 🚀 5. API Endpoints Reference

### Authentication
- `POST /api/v1/auth/login` — Login with email & password, returns JWT token + user profile + permissions.
- `GET /api/v1/auth/me` — Returns authenticated user profile, `role`, and flat array of `permissions` (used by Vue frontend to dynamically adjust UI buttons).

### User Management (Super Admin)
- `GET /api/v1/users` — List all workspace users (`RequirePermission("user:manage")`).
- `PATCH /api/v1/users/:id/role` — Update a user's global role (`RequireRole("super_admin")`).

### Project Management
- `GET /api/v1/projects` — List accessible projects (filtered to user's assigned projects unless super_admin).
- `POST /api/v1/projects` — Create new project (`RequirePermission("project:create")`). Creator is automatically added as `pm` in `project_members`.
- `GET /api/v1/projects/:id` — Get project details (`RequireProjectAccess()`).
- `PATCH /api/v1/projects/:id` — Update project metadata (`RequirePermission("project:edit")`, `RequireProjectAccess()`).
- `DELETE /api/v1/projects/:id` — Delete project (`RequirePermission("project:delete")`).

### Project Member Management
- `GET /api/v1/projects/:id/members` — List members assigned to project (`RequireProjectAccess()`).
- `POST /api/v1/projects/:id/members` — Add/invite member to project (`RequirePermission("project:member:manage")`).
- `DELETE /api/v1/projects/:id/members/:user_id` — Remove member from project (`RequirePermission("project:member:manage")`).

### Task Management
- `GET /api/v1/tasks?project_id=:id` — List tasks for project. For `client`, automatically filters out internal tasks (`client_visible == true` only).
- `POST /api/v1/tasks` — Create task (`RequirePermission("task:create")`).
- `PATCH /api/v1/tasks/:id/status` — Move task on Kanban (`RequirePermission("task:move")`).
- `PATCH /api/v1/tasks/:id` — Edit task details (`RequirePermission("task:edit")`).
- `DELETE /api/v1/tasks/:id` — Delete task (`RequirePermission("task:delete")`).
- `POST /api/v1/tasks/:id/subtasks` — Add subtask checklist item (`RequirePermission("subtask:manage")`).
- `PATCH /api/v1/tasks/:id/subtasks/:sub_id` — Toggle subtask done/undone (`RequirePermission("subtask:manage")`).
- `POST /api/v1/tasks/:id/comments` — Post comment with @mentions (`RequirePermission("comment:create")`).

### Milestone Governance
- `GET /api/v1/milestones?project_id=:id` — List milestones for project.
- `PATCH /api/v1/milestones/:id/approve` — Formally sign off milestone (`RequirePermission("milestone:approve")`). Records client name, ISO timestamp, and locks milestone status.

---

## 💻 6. Running the Backend Server

```powershell
cd backend
go run ./cmd/server
```

Server will start on `http://localhost:8080` with auto-migrated SQLite database and seeded accounts.

package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"pulse-backend/internal/auth"
	"pulse-backend/internal/database"
	"pulse-backend/internal/routes"
)

func setupTestRouter(t *testing.T) *http.Handler {
	// Use temporary test SQLite database
	dbFile := "test_pulse.db"
	_ = os.Remove(dbFile)

	_, err := database.InitDB(dbFile)
	if err != nil {
		t.Fatalf("Failed to init test DB: %v", err)
	}

	database.SeedDatabase()
	r := routes.SetupRouter()
	var handler http.Handler = r
	return &handler
}

func getAuthToken(t *testing.T, handler http.Handler, email, password string) string {
	payload, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})

	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Login failed for %s, status: %d, body: %s", email, w.Code, w.Body.String())
	}

	var res struct {
		Token string `json:"token"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	return res.Token
}

// 1. Verify Permission Matrix
func TestPermissionMatrix(t *testing.T) {
	// super_admin has all permissions
	if !auth.HasPermission(auth.RoleSuperAdmin, auth.PermUserManage) {
		t.Errorf("super_admin should have user:manage permission")
	}

	// pm can create tasks and manage members
	if !auth.HasPermission(auth.RolePM, auth.PermTaskCreate) {
		t.Errorf("pm should have task:create permission")
	}
	if !auth.HasPermission(auth.RolePM, auth.PermProjectMemberManage) {
		t.Errorf("pm should have project:member:manage permission")
	}

	// member can move tasks on kanban, but CANNOT manage members or delete project
	if !auth.HasPermission(auth.RoleMember, auth.PermTaskMove) {
		t.Errorf("member should have task:move permission")
	}
	if auth.HasPermission(auth.RoleMember, auth.PermProjectMemberManage) {
		t.Errorf("member should NOT have project:member:manage permission")
	}
	if auth.HasPermission(auth.RoleMember, auth.PermProjectDelete) {
		t.Errorf("member should NOT have project:delete permission")
	}

	// client can approve milestones and read tasks, but CANNOT create or move tasks
	if !auth.HasPermission(auth.RoleClient, auth.PermMilestoneApprove) {
		t.Errorf("client should have milestone:approve permission")
	}
	if auth.HasPermission(auth.RoleClient, auth.PermTaskCreate) {
		t.Errorf("client should NOT have task:create permission")
	}
	if auth.HasPermission(auth.RoleClient, auth.PermTaskMove) {
		t.Errorf("client should NOT have task:move permission")
	}
}

// 2. Test /api/v1/auth/me returns profile + role + flat array of permissions
func TestAuthMeEndpoint(t *testing.T) {
	handlerPtr := setupTestRouter(t)
	handler := *handlerPtr
	defer os.Remove("test_pulse.db")

	token := getAuthToken(t, handler, "pm@pulse.workspace", "Password123!")

	req, _ := http.NewRequest("GET", "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}

	var res struct {
		Role        string   `json:"role"`
		Permissions []string `json:"permissions"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &res)

	if res.Role != "pm" {
		t.Errorf("Expected role 'pm', got '%s'", res.Role)
	}

	if len(res.Permissions) == 0 {
		t.Errorf("Expected non-empty permissions array for PM")
	}
}

// 3. Test Project Access Control (RBAC & Project Member Isolation)
func TestProjectLevelAuthorization(t *testing.T) {
	handlerPtr := setupTestRouter(t)
	handler := *handlerPtr
	defer os.Remove("test_pulse.db")

	clientToken := getAuthToken(t, handler, "client@pulse.workspace", "Password123!")
	adminToken := getAuthToken(t, handler, "admin@pulse.workspace", "Password123!")
	pmToken := getAuthToken(t, handler, "pm@pulse.workspace", "Password123!")
	memberToken := getAuthToken(t, handler, "member@pulse.workspace", "Password123!")

	// A. Client accesses assigned project (proj_sustagri) -> Should be 200 OK
	reqA, _ := http.NewRequest("GET", "/api/v1/projects/proj_sustagri", nil)
	reqA.Header.Set("Authorization", "Bearer "+clientToken)
	wA := httptest.NewRecorder()
	handler.ServeHTTP(wA, reqA)
	if wA.Code != http.StatusOK {
		t.Errorf("Client should access assigned project proj_sustagri, got: %d", wA.Code)
	}

	// B. Client accesses unassigned project (proj_smart) -> MUST BE 403 Forbidden!
	reqB, _ := http.NewRequest("GET", "/api/v1/projects/proj_smart", nil)
	reqB.Header.Set("Authorization", "Bearer "+clientToken)
	wB := httptest.NewRecorder()
	handler.ServeHTTP(wB, reqB)
	if wB.Code != http.StatusForbidden {
		t.Errorf("Client should be 403 Forbidden on unassigned project proj_smart, got: %d", wB.Code)
	}

	// C. Super Admin accesses unassigned project (proj_smart) -> Should BYPASS and return 200 OK!
	reqC, _ := http.NewRequest("GET", "/api/v1/projects/proj_smart", nil)
	reqC.Header.Set("Authorization", "Bearer "+adminToken)
	wC := httptest.NewRecorder()
	handler.ServeHTTP(wC, reqC)
	if wC.Code != http.StatusOK {
		t.Errorf("Super Admin should bypass project check and get 200 OK, got: %d", wC.Code)
	}

	// D. Client tries to create task -> MUST BE 403 Forbidden (missing task:create permission)
	taskPayload, _ := json.Marshal(map[string]interface{}{
		"project_id": "proj_sustagri",
		"title":      "Unauthorized Client Task",
	})
	reqD, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(taskPayload))
	reqD.Header.Set("Authorization", "Bearer "+clientToken)
	reqD.Header.Set("Content-Type", "application/json")
	wD := httptest.NewRecorder()
	handler.ServeHTTP(wD, reqD)
	if wD.Code != http.StatusForbidden {
		t.Errorf("Client must NOT be allowed to create tasks, expected 403, got: %d", wD.Code)
	}

	// E. Client approves milestone -> Should be 200 OK (has milestone:approve)
	approvePayload, _ := json.Marshal(map[string]string{
		"client_name": "Ir. H. Gunawan",
	})
	reqE, _ := http.NewRequest("PATCH", "/api/v1/milestones/ms_sust_f1/approve", bytes.NewBuffer(approvePayload))
	reqE.Header.Set("Authorization", "Bearer "+clientToken)
	reqE.Header.Set("Content-Type", "application/json")
	wE := httptest.NewRecorder()
	handler.ServeHTTP(wE, reqE)
	if wE.Code != http.StatusOK {
		t.Errorf("Client should be allowed to approve milestone, expected 200, got: %d", wE.Code)
	}

	// F. Member tries to invite member to project -> MUST BE 403 Forbidden (restricted to PM and Super Admin)
	memberPayload, _ := json.Marshal(map[string]string{
		"user_id":      "usr_client",
		"project_role": "client",
	})
	reqF, _ := http.NewRequest("POST", "/api/v1/projects/proj_sustagri/members", bytes.NewBuffer(memberPayload))
	reqF.Header.Set("Authorization", "Bearer "+memberToken)
	reqF.Header.Set("Content-Type", "application/json")
	wF := httptest.NewRecorder()
	handler.ServeHTTP(wF, reqF)
	if wF.Code != http.StatusForbidden {
		t.Errorf("Member must NOT manage project members, expected 403, got: %d", wF.Code)
	}

	// G. PM invites member to project -> Should be 200 OK
	reqG, _ := http.NewRequest("POST", "/api/v1/projects/proj_smart/members", bytes.NewBuffer(memberPayload))
	reqG.Header.Set("Authorization", "Bearer "+pmToken)
	reqG.Header.Set("Content-Type", "application/json")
	wG := httptest.NewRecorder()
	handler.ServeHTTP(wG, reqG)
	if wG.Code != http.StatusOK {
		t.Errorf("PM should be allowed to add members to project, expected 200, got: %d", wG.Code)
	}

	// H. Non-admin tries to change global user role -> MUST BE 403 Forbidden (RequireRole("super_admin"))
	rolePayload, _ := json.Marshal(map[string]string{
		"role": "pm",
	})
	reqH, _ := http.NewRequest("PATCH", "/api/v1/users/usr_member/role", bytes.NewBuffer(rolePayload))
	reqH.Header.Set("Authorization", "Bearer "+pmToken)
	reqH.Header.Set("Content-Type", "application/json")
	wH := httptest.NewRecorder()
	handler.ServeHTTP(wH, reqH)
	if wH.Code != http.StatusForbidden {
		t.Errorf("Non-super_admin must NOT be allowed to update user roles, expected 403, got: %d", wH.Code)
	}

	// I. Super Admin changes user role -> Should be 200 OK
	reqI, _ := http.NewRequest("PATCH", "/api/v1/users/usr_member/role", bytes.NewBuffer(rolePayload))
	reqI.Header.Set("Authorization", "Bearer "+adminToken)
	reqI.Header.Set("Content-Type", "application/json")
	wI := httptest.NewRecorder()
	handler.ServeHTTP(wI, reqI)
	if wI.Code != http.StatusOK {
		t.Errorf("Super admin should update user role, expected 200, got: %d", wI.Code)
	}
}

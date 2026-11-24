package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"pulse-backend/internal/auth"
	"pulse-backend/internal/database"
	"pulse-backend/internal/mailer"
	"pulse-backend/internal/middleware"
	"pulse-backend/internal/models"
)

type CreateInvitationRequest struct {
	Email      string   `json:"email" binding:"required,email"`
	Name       string   `json:"name"`
	Role       string   `json:"role" binding:"required"`
	ProjectIDs []string `json:"project_ids"`
}

type AcceptInvitationRequest struct {
	Token    string `json:"token" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// getPublicBaseURL resolves the canonical public URL for external client links
func getPublicBaseURL(c *gin.Context) string {
	// 1. Priority: Explicit environment variable APP_BASE_URL or PUBLIC_URL
	if base := strings.TrimSpace(os.Getenv("APP_BASE_URL")); base != "" {
		return strings.TrimRight(base, "/")
	}
	if base := strings.TrimSpace(os.Getenv("PUBLIC_URL")); base != "" {
		return strings.TrimRight(base, "/")
	}

	// 2. Check X-Forwarded-Host from reverse proxy (Cloudflare Tunnel / Node Gateway)
	forwardedHost := c.GetHeader("X-Forwarded-Host")
	if forwardedHost != "" && !strings.Contains(forwardedHost, "127.0.0.1") && !strings.Contains(forwardedHost, "localhost") {
		proto := c.GetHeader("X-Forwarded-Proto")
		if proto == "" {
			proto = "https"
		}
		return fmt.Sprintf("%s://%s", proto, forwardedHost)
	}

	// 3. Check Host header if it is already a public domain
	host := c.Request.Host
	if host != "" && !strings.Contains(host, "127.0.0.1") && !strings.Contains(host, "localhost") {
		proto := "https"
		if c.Request.TLS == nil && !strings.Contains(host, "kedukanaksara.com") {
			proto = "http"
		}
		return fmt.Sprintf("%s://%s", proto, host)
	}

	// 4. Default production subdomain fallback
	return "https://project.kedukanaksara.com"
}

// CreateInvitation allows Admin or PM to issue an invitation to a client, member, or PM
func CreateInvitation(c *gin.Context) {
	currentUserID, currentUserRole := middleware.GetCurrentUser(c)
	if currentUserRole != auth.RoleSuperAdmin && currentUserRole != auth.RolePM {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya Administrator atau Project Manager yang dapat mengundang pengguna"})
		return
	}

	var req CreateInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format permintaan tidak valid: email dan peran wajib diisi"})
		return
	}

	cleanEmail := strings.TrimSpace(strings.ToLower(req.Email))
	cleanName := strings.TrimSpace(req.Name)

	if !auth.IsValidRole(req.Role) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Peran tidak valid. Pilihan: client, member, pm, super_admin"})
		return
	}

	// 1. Check if user already exists in users table and is not deleted
	var existingUser models.User
	if err := database.DB.Where("LOWER(email) = ?", cleanEmail).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Pengguna dengan email ini sudah terdaftar dan aktif di sistem"})
		return
	}

	// 2. Check if there is an active pending invitation for this email; if yes, revoke previous one
	database.DB.Model(&models.Invitation{}).
		Where("LOWER(email) = ? AND status = ?", cleanEmail, "pending").
		Update("status", "revoked")

	// 3. Serialize project IDs
	var cleanProjectIDs []string
	for _, pid := range req.ProjectIDs {
		pid = strings.TrimSpace(pid)
		if pid != "" {
			cleanProjectIDs = append(cleanProjectIDs, pid)
		}
	}
	projectsJSON, _ := json.Marshal(cleanProjectIDs)

	// 4. Generate unique secure token
	token := fmt.Sprintf("inv_%s%s", strings.ReplaceAll(uuid.New().String(), "-", ""), strings.ReplaceAll(uuid.New().String()[:8], "-", ""))

	invID := fmt.Sprintf("inv_%s", strings.ReplaceAll(uuid.New().String()[:12], "-", ""))
	now := time.Now()
	expiresAt := now.Add(7 * 24 * time.Hour) // Valid for 7 days

	invitation := models.Invitation{
		ID:         invID,
		Email:      cleanEmail,
		Name:       cleanName,
		Role:       req.Role,
		Token:      token,
		Status:     "pending",
		ProjectIDs: string(projectsJSON),
		InvitedBy:  currentUserID,
		ExpiresAt:  expiresAt,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := database.DB.Create(&invitation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data undangan: " + err.Error()})
		return
	}

	// Determine client-accessible base URL (defaults to production subdomain)
	baseURL := getPublicBaseURL(c)
	inviteURL := fmt.Sprintf("%s/accept-invite?token=%s", baseURL, token)

	invitation.ProjectList = cleanProjectIDs
	invitation.InviteURL = inviteURL

	// 5. Fetch inviter display name
	var inviter models.User
	database.DB.Select("name").Where("id = ?", currentUserID).First(&inviter)
	inviterName := inviter.Name
	if inviterName == "" {
		inviterName = "Administrator"
	}

	// 6. Fetch project names for email summary
	var projectNames []string
	if len(cleanProjectIDs) > 0 {
		database.DB.Model(&models.Project{}).Where("id IN ?", cleanProjectIDs).Pluck("name", &projectNames)
	}

	// 7. Dispatch invitation email via SMTP if configured
	emailDispatched := false
	if mailer.IsConfigured() {
		emailDispatched = true
		go func(toEmail, toName, invName, role string, projs []string, link string) {
			if err := mailer.SendInvitationEmail(toEmail, toName, invName, role, projs, link); err != nil {
				log.Printf("[Invitation] SMTP dispatch error to %s: %v\n", toEmail, err)
			}
		}(cleanEmail, cleanName, inviterName, req.Role, projectNames, inviteURL)
	} else {
		log.Printf("[Invitation] Layanan SMTP belum dikonfigurasi. Tautan aktivasi untuk %s: %s\n", cleanEmail, inviteURL)
	}

	msg := "Undangan berhasil dibuat dan tautan aktivasi siap dibagikan"
	if emailDispatched {
		msg = fmt.Sprintf("Undangan berhasil dibuat dan email aktivasi sedang dikirimkan ke %s", cleanEmail)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":          msg,
		"invitation":       invitation,
		"invite_url":       inviteURL,
		"email_dispatched": emailDispatched,
	})
}

// ListInvitations returns pending and historical invitations for admin review
func ListInvitations(c *gin.Context) {
	_, currentUserRole := middleware.GetCurrentUser(c)
	if currentUserRole != auth.RoleSuperAdmin && currentUserRole != auth.RolePM {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak"})
		return
	}

	var invitations []models.Invitation
	if err := database.DB.Order("created_at desc").Find(&invitations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat daftar undangan: " + err.Error()})
		return
	}

	// Fetch inviters map
	var users []models.User
	database.DB.Select("id, name").Find(&users)
	userNameMap := make(map[string]string)
	for _, u := range users {
		userNameMap[u.ID] = u.Name
	}

	baseURL := getPublicBaseURL(c)

	now := time.Now()
	for i := range invitations {
		// Auto update expired status in view
		if invitations[i].Status == "pending" && invitations[i].ExpiresAt.Before(now) {
			invitations[i].Status = "expired"
		}

		var pids []string
		if invitations[i].ProjectIDs != "" {
			_ = json.Unmarshal([]byte(invitations[i].ProjectIDs), &pids)
		}
		if pids == nil {
			pids = []string{}
		}
		invitations[i].ProjectList = pids
		invitations[i].InviterName = userNameMap[invitations[i].InvitedBy]
		invitations[i].InviteURL = fmt.Sprintf("%s/accept-invite?token=%s", baseURL, invitations[i].Token)
	}

	c.JSON(http.StatusOK, invitations)
}

// RevokeInvitation cancels a pending invitation
func RevokeInvitation(c *gin.Context) {
	_, currentUserRole := middleware.GetCurrentUser(c)
	if currentUserRole != auth.RoleSuperAdmin && currentUserRole != auth.RolePM {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak"})
		return
	}

	invID := c.Param("id")
	if invID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID undangan wajib diisi"})
		return
	}

	var inv models.Invitation
	if err := database.DB.Where("id = ?", invID).First(&inv).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Undangan tidak ditemukan"})
		return
	}

	inv.Status = "revoked"
	inv.UpdatedAt = time.Now()
	database.DB.Save(&inv)

	c.JSON(http.StatusOK, gin.H{
		"message": "Undangan berhasil dibatalkan",
		"status":  "revoked",
	})
}

// VerifyInvitation verifies the token from the public accept-invite page
func VerifyInvitation(c *gin.Context) {
	token := strings.TrimSpace(c.Query("token"))
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter token undangan wajib disertakan"})
		return
	}

	var inv models.Invitation
	if err := database.DB.Where("token = ?", token).First(&inv).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tautan undangan tidak ditemukan atau tidak valid"})
		return
	}

	if inv.Status == "accepted" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Undangan ini sudah pernah digunakan untuk mendaftar akun",
			"code":      "ALREADY_ACCEPTED",
			"email":     inv.Email,
			"is_usable": false,
		})
		return
	}

	if inv.Status == "revoked" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Undangan ini telah dibatalkan oleh Administrator",
			"code":      "REVOKED",
			"is_usable": false,
		})
		return
	}

	if inv.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Tautan undangan ini telah kedaluwarsa. Silakan hubungi Administrator untuk meminta undangan baru.",
			"code":      "EXPIRED",
			"is_usable": false,
		})
		return
	}

	// Fetch project metadata
	var pids []string
	if inv.ProjectIDs != "" {
		_ = json.Unmarshal([]byte(inv.ProjectIDs), &pids)
	}

	var projects []models.Project
	if len(pids) > 0 {
		database.DB.Where("id IN (?)", pids).Select("id, name, slug, client_name, color_dot").Find(&projects)
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":      true,
		"email":      inv.Email,
		"name":       inv.Name,
		"role":       inv.Role,
		"projects":   projects,
		"expires_at": inv.ExpiresAt,
	})
}

type ConfirmInvitationRequest struct {
	Token string `json:"token" binding:"required"`
}

// ConfirmInvitation processes 1-click confirmation: generates password, activates user account,
// sends credentials to email, and returns details for the confirmation page.
func ConfirmInvitation(c *gin.Context) {
	var req ConfirmInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter token undangan wajib disertakan"})
		return
	}

	token := strings.TrimSpace(req.Token)
	var inv models.Invitation
	if err := database.DB.Where("token = ?", token).First(&inv).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tautan undangan tidak ditemukan atau tidak valid"})
		return
	}

	// Fetch project metadata
	var pids []string
	if inv.ProjectIDs != "" {
		_ = json.Unmarshal([]byte(inv.ProjectIDs), &pids)
	}
	var projectNames []string
	if len(pids) > 0 {
		database.DB.Model(&models.Project{}).Where("id IN ?", pids).Pluck("name", &projectNames)
	}

	// Check if already accepted
	if inv.Status == "accepted" {
		c.JSON(http.StatusOK, gin.H{
			"status":    "already_confirmed",
			"email":     inv.Email,
			"name":      inv.Name,
			"role":      inv.Role,
			"projects":  projectNames,
			"message":   "Akun Anda sudah aktif sebelumnya. Kredensial login (username & password) telah dikirimkan ke email Anda.",
			"login_url": "/login",
		})
		return
	}

	// If revoked, check if newer active invite exists for this email
	if inv.Status == "revoked" {
		var activeInv models.Invitation
		if err := database.DB.Where("LOWER(email) = ? AND status = ?", strings.ToLower(inv.Email), "pending").Order("created_at desc").First(&activeInv).Error; err == nil {
			inv = activeInv
		} else {
			var acceptedInv models.Invitation
			if err := database.DB.Where("LOWER(email) = ? AND status = ?", strings.ToLower(inv.Email), "accepted").Order("updated_at desc").First(&acceptedInv).Error; err == nil {
				c.JSON(http.StatusOK, gin.H{
					"status":    "already_confirmed",
					"email":     acceptedInv.Email,
					"name":      acceptedInv.Name,
					"role":      acceptedInv.Role,
					"projects":  projectNames,
					"message":   "Akun Anda sudah aktif sebelumnya. Kredensial login (username & password) telah dikirimkan ke email Anda.",
					"login_url": "/login",
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Tautan undangan ini telah dibatalkan oleh Administrator. Silakan hubungi Administrator untuk meminta undangan baru.",
			})
			return
		}
	}

	if inv.ExpiresAt.Before(time.Now()) {
		var acceptedInv models.Invitation
		if err := database.DB.Where("LOWER(email) = ? AND status = ?", strings.ToLower(inv.Email), "accepted").Order("updated_at desc").First(&acceptedInv).Error; err == nil {
			c.JSON(http.StatusOK, gin.H{
				"status":    "already_confirmed",
				"email":     acceptedInv.Email,
				"name":      acceptedInv.Name,
				"role":      acceptedInv.Role,
				"projects":  projectNames,
				"message":   "Akun Anda sudah aktif sebelumnya. Kredensial login (username & password) telah dikirimkan ke email Anda.",
				"login_url": "/login",
			})
			return
		}
		inv.Status = "expired"
		database.DB.Save(&inv)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tautan undangan ini telah kedaluwarsa. Silakan hubungi Administrator untuk meminta undangan baru.",
		})
		return
	}

	// Generate secure readable temporary password: Pulse@<random 6 digits>
	rand.Seed(time.Now().UnixNano())
	randomCode := 100000 + rand.Intn(900000)
	generatedPassword := fmt.Sprintf("Pulse@%d", randomCode)

	cleanEmail := strings.TrimSpace(strings.ToLower(inv.Email))
	cleanName := strings.TrimSpace(inv.Name)
	if cleanName == "" {
		cleanName = strings.Split(cleanEmail, "@")[0]
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(generatedPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi kata sandi"})
		return
	}

	now := time.Now()
	sessionID := uuid.New().String()
	avatar := fmt.Sprintf("https://api.dicebear.com/7.x/avataaars/svg?seed=%s", url.QueryEscape(cleanName))

	// Upsert User
	var user models.User
	if err := database.DB.Unscoped().Where("LOWER(email) = ?", cleanEmail).First(&user).Error; err == nil {
		if user.DeletedAt.Valid {
			database.DB.Unscoped().Delete(&user)
			user = models.User{
				ID:               fmt.Sprintf("usr_%s", strings.ReplaceAll(uuid.New().String()[:12], "-", "")),
				Name:             cleanName,
				Email:            cleanEmail,
				Password:         string(hashedPassword),
				Role:             inv.Role,
				AvatarURL:        avatar,
				CurrentSessionID: sessionID,
				LastActiveAt:     now,
				CreatedAt:        now,
				UpdatedAt:        now,
			}
			database.DB.Create(&user)
		} else {
			// Update password and role
			user.Password = string(hashedPassword)
			user.Role = inv.Role
			if cleanName != "" {
				user.Name = cleanName
			}
			user.UpdatedAt = now
			database.DB.Save(&user)
		}
	} else {
		user = models.User{
			ID:               fmt.Sprintf("usr_%s", strings.ReplaceAll(uuid.New().String()[:12], "-", "")),
			Name:             cleanName,
			Email:            cleanEmail,
			Password:         string(hashedPassword),
			Role:             inv.Role,
			AvatarURL:        avatar,
			CurrentSessionID: sessionID,
			LastActiveAt:     now,
			CreatedAt:        now,
			UpdatedAt:        now,
		}
		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat akun: " + err.Error()})
			return
		}
	}

	// Assign project memberships
	projRole := "member"
	if user.Role == "pm" {
		projRole = "pm"
	} else if user.Role == "client" {
		projRole = "client"
	}

	for _, pid := range pids {
		pid = strings.TrimSpace(pid)
		if pid == "" {
			continue
		}
		pm := models.ProjectMember{
			ProjectID:   pid,
			UserID:      user.ID,
			ProjectRole: projRole,
			JoinedAt:    now,
		}
		database.DB.Where(models.ProjectMember{ProjectID: pid, UserID: user.ID}).FirstOrCreate(&pm)
	}
	user.ProjectIDs = pids

	// Mark invitation accepted
	inv.Status = "accepted"
	inv.UpdatedAt = now
	database.DB.Save(&inv)

	// Send email with credentials via SMTP
	loginURL := fmt.Sprintf("%s/login", getPublicBaseURL(c))
	if mailer.IsConfigured() {
		go func(toEmail, toName, rawPass, role string, projs []string, loginLink string) {
			if err := mailer.SendWelcomeCredentialsEmail(toEmail, toName, rawPass, role, projs, loginLink); err != nil {
				log.Printf("[Invitation] Gagal mengirimkan kredensial login ke %s: %v\n", toEmail, err)
			}
		}(user.Email, user.Name, generatedPassword, user.Role, projectNames, loginURL)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"status":    "confirmed",
		"email":     user.Email,
		"name":      user.Name,
		"password":  generatedPassword,
		"role":      user.Role,
		"projects":  projectNames,
		"message":   "Konfirmasi akun berhasil! Username dan password telah dikirimkan ke email Anda.",
		"login_url": "/login",
	})
}

// AcceptInvitation allows the recipient to confirm their name and password, creating their account
func AcceptInvitation(c *gin.Context) {
	var req AcceptInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap: Nama dan kata sandi minimal 6 karakter"})
		return
	}

	cleanToken := strings.TrimSpace(req.Token)
	cleanName := strings.TrimSpace(req.Name)

	var inv models.Invitation
	if err := database.DB.Where("token = ? AND status = ?", cleanToken, "pending").First(&inv).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tautan undangan tidak valid atau sudah digunakan"})
		return
	}

	if inv.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tautan undangan telah kedaluwarsa"})
		return
	}

	cleanEmail := strings.TrimSpace(strings.ToLower(inv.Email))

	// Check if user already exists
	var existing models.User
	if err := database.DB.Unscoped().Where("LOWER(email) = ?", cleanEmail).First(&existing).Error; err == nil {
		if existing.DeletedAt.Valid {
			// Permanently remove soft-deleted ghost user
			database.DB.Unscoped().Delete(&existing)
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": "Email ini sudah memiliki akun aktif di sistem. Silakan langsung login."})
			return
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi kata sandi"})
		return
	}

	// Create user
	newID := fmt.Sprintf("usr_%s", strings.ReplaceAll(uuid.New().String()[:12], "-", ""))
	sessionID := uuid.New().String()
	avatar := fmt.Sprintf("https://api.dicebear.com/7.x/avataaars/svg?seed=%s", url.QueryEscape(cleanName))

	now := time.Now()
	user := models.User{
		ID:               newID,
		Name:             cleanName,
		Email:            cleanEmail,
		Password:         string(hashedPassword),
		Role:             inv.Role,
		AvatarURL:        avatar,
		CurrentSessionID: sessionID,
		LastActiveAt:     now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat akun pengguna: " + err.Error()})
		return
	}

	// Assign project memberships
	var pids []string
	if inv.ProjectIDs != "" {
		_ = json.Unmarshal([]byte(inv.ProjectIDs), &pids)
	}

	projRole := "member"
	if user.Role == "pm" {
		projRole = "pm"
	} else if user.Role == "client" {
		projRole = "client"
	}

	for _, pid := range pids {
		pid = strings.TrimSpace(pid)
		if pid == "" {
			continue
		}
		pm := models.ProjectMember{
			ProjectID:   pid,
			UserID:      user.ID,
			ProjectRole: projRole,
			JoinedAt:    now,
		}
		database.DB.Where(models.ProjectMember{ProjectID: pid, UserID: user.ID}).FirstOrCreate(&pm)
	}
	user.ProjectIDs = pids

	// Mark invitation accepted
	inv.Status = "accepted"
	inv.UpdatedAt = now
	database.DB.Save(&inv)

	// Fetch project titles for welcome credentials email
	var projectTitles []string
	if len(pids) > 0 {
		database.DB.Model(&models.Project{}).Where("id IN ?", pids).Pluck("title", &projectTitles)
	}

	// Dispatch username and password email via SMTP
	if mailer.IsConfigured() {
		loginURL := fmt.Sprintf("%s/login", getPublicBaseURL(c))
		go func(toEmail, toName, rawPassword, role string, projs []string, loginLink string) {
			if err := mailer.SendWelcomeCredentialsEmail(toEmail, toName, rawPassword, role, projs, loginLink); err != nil {
				log.Printf("[Invitation] Gagal mengirimkan kredensial login ke %s: %v\n", toEmail, err)
			}
		}(user.Email, user.Name, req.Password, user.Role, projectTitles, loginURL)
	}

	// Issue JWT token and register active session
	token, err := auth.GenerateToken(user.ID, user.Email, user.Role, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Akun berhasil dibuat namun gagal membuat sesi login"})
		return
	}
	auth.SetActiveSession(user.ID, sessionID)

	permissions := auth.GetPermissionsForRole(user.Role)

	redirectURL := "/dashboard"
	if user.Role == auth.RoleClient {
		redirectURL = "/client-portal"
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Akun Anda berhasil dikonfirmasi dan diaktifkan! Kredensial login telah dikirimkan ke email Anda.",
		"token":        token,
		"user":         user,
		"permissions":  permissions,
		"redirect_url": redirectURL,
	})
}

// ResendInvitationEmail manually triggers email dispatch for a pending invitation
func ResendInvitationEmail(c *gin.Context) {
	_, currentUserRole := middleware.GetCurrentUser(c)
	if currentUserRole != auth.RoleSuperAdmin && currentUserRole != auth.RolePM {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya Administrator atau Project Manager yang dapat mengirim ulang undangan"})
		return
	}

	invID := c.Param("id")
	var inv models.Invitation
	if err := database.DB.Where("id = ?", invID).First(&inv).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data undangan tidak ditemukan"})
		return
	}

	if inv.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Tidak dapat mengirim ulang email. Status undangan saat ini: %s", inv.Status)})
		return
	}

	if time.Now().After(inv.ExpiresAt) {
		inv.Status = "expired"
		database.DB.Save(&inv)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tautan undangan sudah kadaluarsa (melewati 7 hari). Silakan buat undangan baru."})
		return
	}

	// Fetch inviter
	var inviter models.User
	database.DB.Select("name").Where("id = ?", inv.InvitedBy).First(&inviter)
	inviterName := inviter.Name
	if inviterName == "" {
		inviterName = "Administrator"
	}

	// Fetch project titles
	var pids []string
	_ = json.Unmarshal([]byte(inv.ProjectIDs), &pids)
	var projectTitles []string
	if len(pids) > 0 {
		database.DB.Model(&models.Project{}).Where("id IN ?", pids).Pluck("title", &projectTitles)
	}

	baseURL := getPublicBaseURL(c)
	inviteURL := fmt.Sprintf("%s/accept-invite?token=%s", baseURL, inv.Token)

	if !mailer.IsConfigured() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Layanan SMTP belum dikonfigurasi pada server. Silakan atur kredensial SMTP terlebih dahulu."})
		return
	}

	if err := mailer.SendInvitationEmail(inv.Email, inv.Name, inviterName, inv.Role, projectTitles, inviteURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirimkan email via SMTP: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    fmt.Sprintf("Email undangan berhasil dikirimkan ulang ke %s", inv.Email),
		"invite_url": inviteURL,
	})
}

// GetSMTPStatusHandler returns current SMTP configuration state
func GetSMTPStatusHandler(c *gin.Context) {
	c.JSON(http.StatusOK, mailer.GetStatus())
}

// TestSMTPHandler sends a test email to verify SMTP credentials
func TestSMTPHandler(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Alamat email penerima pengujian wajib diisi"})
		return
	}

	if err := mailer.SendTestEmail(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengirim email pengujian: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Email pengujian SMTP sukses terkirim ke %s!", req.Email),
	})
}


package handlers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"pulse-backend/internal/auth"
	"pulse-backend/internal/database"
	"pulse-backend/internal/middleware"
	"pulse-backend/internal/models"
)

// sanitizeCSVField mitigates CSV Formula Injection (OWASP ASVS 5.1.5)
func sanitizeCSVField(s string) string {
	if len(s) == 0 {
		return s
	}
	// Characters that can trigger formula execution in spreadsheet software
	if s[0] == '=' || s[0] == '+' || s[0] == '-' || s[0] == '@' || s[0] == '\t' || s[0] == '\r' {
		return "'" + s
	}
	return s
}

// ExportProjectsCSV exports all projects accessible to the user as a downloadable CSV file
func ExportProjectsCSV(c *gin.Context) {
	userID, role := middleware.GetCurrentUser(c)

	var projects []models.Project
	if role == auth.RoleSuperAdmin {
		if err := database.DB.Order("created_at asc").Find(&projects).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects for export"})
			return
		}
	} else {
		subQuery := database.DB.Table("project_members").Select("project_id").Where("user_id = ?", userID)
		if err := database.DB.Where("id IN (?)", subQuery).Order("created_at asc").Find(&projects).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects for export"})
			return
		}
	}

	// Calculate tasks stats
	for i := range projects {
		var total, completed int64
		database.DB.Model(&models.Task{}).Where("project_id = ?", projects[i].ID).Count(&total)
		database.DB.Model(&models.Task{}).Where("project_id = ? AND status = ?", projects[i].ID, "Done").Count(&completed)
		projects[i].TotalTasks = int(total)
		projects[i].CompletedTasks = int(completed)
		if total > 0 {
			projects[i].ProgressPercentage = int((completed * 100) / total)
		}
	}

	buf := &bytes.Buffer{}
	// UTF-8 BOM for Microsoft Excel compatibility
	buf.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(buf)

	// CSV Header
	header := []string{
		"ID Proyek",
		"Nama Proyek",
		"Kode Slug",
		"Nama Klien",
		"Status",
		"Progress (%)",
		"Tenggat Waktu",
		"Total Tugas",
		"Tugas Selesai",
		"Warna Tag",
		"Tanggal Dibuat",
	}
	_ = writer.Write(header)

	// CSV Rows (sanitized against CSV Formula Injection)
	for _, p := range projects {
		row := []string{
			sanitizeCSVField(p.ID),
			sanitizeCSVField(p.Name),
			sanitizeCSVField(p.Slug),
			sanitizeCSVField(p.ClientName),
			sanitizeCSVField(p.Status),
			strconv.Itoa(p.ProgressPercentage),
			sanitizeCSVField(p.Deadline),
			strconv.Itoa(p.TotalTasks),
			strconv.Itoa(p.CompletedTasks),
			sanitizeCSVField(p.ColorDot),
			p.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		_ = writer.Write(row)
	}

	writer.Flush()

	filename := fmt.Sprintf("pulse_projects_export_%s.csv", time.Now().Format("20060102_150405"))
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

// ExportTasksCSV exports tasks as a downloadable CSV file
func ExportTasksCSV(c *gin.Context) {
	userID, role := middleware.GetCurrentUser(c)
	projectIDFilter := c.Query("project_id")

	var tasks []models.Task
	query := database.DB.Model(&models.Task{})

	if projectIDFilter != "" {
		query = query.Where("project_id = ?", projectIDFilter)
	} else if role != auth.RoleSuperAdmin {
		subQuery := database.DB.Table("project_members").Select("project_id").Where("user_id = ?", userID)
		query = query.Where("project_id IN (?)", subQuery)
	}

	if err := query.Order("project_id asc, created_at desc").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks for export"})
		return
	}

	// Fetch project map for readable project names
	var projects []models.Project
	database.DB.Select("id, name").Find(&projects)
	projNameMap := make(map[string]string)
	for _, p := range projects {
		projNameMap[p.ID] = p.Name
	}

	buf := &bytes.Buffer{}
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	writer := csv.NewWriter(buf)

	header := []string{
		"ID Tugas",
		"ID Proyek",
		"Nama Proyek",
		"Kode Tugas",
		"Judul Tugas",
		"Kategori",
		"Prioritas",
		"Status",
		"Assignee ID",
		"Tenggat Waktu",
		"Deskripsi",
		"Tanggal Dibuat",
	}
	_ = writer.Write(header)

	for _, t := range tasks {
		projName := projNameMap[t.ProjectID]
		if projName == "" {
			projName = t.ProjectID
		}
		row := []string{
			sanitizeCSVField(t.ID),
			sanitizeCSVField(t.ProjectID),
			sanitizeCSVField(projName),
			sanitizeCSVField(t.Key),
			sanitizeCSVField(t.Title),
			sanitizeCSVField(t.Category),
			sanitizeCSVField(t.Priority),
			sanitizeCSVField(t.Status),
			sanitizeCSVField(t.AssigneeID),
			sanitizeCSVField(t.DueDate),
			sanitizeCSVField(t.Description),
			t.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		_ = writer.Write(row)
	}

	writer.Flush()

	filename := fmt.Sprintf("pulse_tasks_export_%s.csv", time.Now().Format("20060102_150405"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

// GetProjectImportTemplate generates a sample CSV template for users to fill in Microsoft Excel or Calc
func GetProjectImportTemplate(c *gin.Context) {
	buf := &bytes.Buffer{}
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	writer := csv.NewWriter(buf)

	header := []string{
		"Nama Proyek",
		"Kode Slug",
		"Nama Klien",
		"Status",
		"Tenggat Waktu",
		"Warna Tag",
	}
	_ = writer.Write(header)

	// Two sample rows for user guidance
	sample1 := []string{
		"Sustagri Fase 3 (Ekspansi Lapangan)",
		"sustagri-fase-3",
		"PT Agro Lestari Nusantara",
		"active",
		"2026-12-31",
		"#10B981",
	}
	sample2 := []string{
		"Smart IoT Tracking Device",
		"smart-iot-tracking",
		"PT Smart Solusindo",
		"active",
		"2026-11-30",
		"#6366F1",
	}
	_ = writer.Write(sample1)
	_ = writer.Write(sample2)

	writer.Flush()

	c.Header("Content-Disposition", "attachment; filename=template_import_proyek.csv")
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

type ProjectImportItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	ClientName string `json:"client_name"`
	Status     string `json:"status"`
	Deadline   string `json:"deadline"`
	ColorDot   string `json:"color_dot"`
}

type ImportProjectsJSONRequest struct {
	Projects []ProjectImportItem `json:"projects"`
}

// ImportProjects parses uploaded CSV file or JSON array and creates/updates projects
func ImportProjects(c *gin.Context) {
	userID, _ := middleware.GetCurrentUser(c)

	var items []ProjectImportItem

	// Check if uploaded as multipart file
	fileHeader, err := c.FormFile("file")
	if err == nil && fileHeader != nil {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read uploaded file: " + err.Error()})
			return
		}
		defer file.Close()

		// Read CSV bytes and strip potential UTF-8 BOM
		rawBytes, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file content"})
			return
		}
		rawBytes = bytes.TrimPrefix(rawBytes, []byte{0xEF, 0xBB, 0xBF})

		r := csv.NewReader(bytes.NewReader(rawBytes))
		r.LazyQuotes = true
		r.TrimLeadingSpace = true

		records, err := r.ReadAll()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CSV file format: " + err.Error()})
			return
		}

		if len(records) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File CSV harus memiliki baris header dan minimal 1 baris data"})
			return
		}

		// Detect header indices
		header := records[0]
		nameIdx, slugIdx, clientIdx, statusIdx, deadlineIdx, colorIdx := -1, -1, -1, -1, -1, -1

		for idx, col := range header {
			norm := strings.ToLower(strings.TrimSpace(col))
			if strings.Contains(norm, "nama") && !strings.Contains(norm, "klien") {
				nameIdx = idx
			} else if strings.Contains(norm, "slug") || strings.Contains(norm, "kode") {
				slugIdx = idx
			} else if strings.Contains(norm, "klien") || strings.Contains(norm, "client") {
				clientIdx = idx
			} else if strings.Contains(norm, "status") {
				statusIdx = idx
			} else if strings.Contains(norm, "tenggat") || strings.Contains(norm, "deadline") {
				deadlineIdx = idx
			} else if strings.Contains(norm, "warna") || strings.Contains(norm, "color") {
				colorIdx = idx
			}
		}

		if nameIdx == -1 {
			// Fallback: column 0 is Name
			nameIdx = 0
		}

		for lineNo, row := range records[1:] {
			if len(row) <= nameIdx {
				continue
			}
			name := strings.TrimSpace(row[nameIdx])
			if name == "" {
				continue
			}

			item := ProjectImportItem{
				Name: name,
			}
			if slugIdx != -1 && len(row) > slugIdx {
				item.Slug = strings.TrimSpace(row[slugIdx])
			}
			if clientIdx != -1 && len(row) > clientIdx {
				item.ClientName = strings.TrimSpace(row[clientIdx])
			}
			if statusIdx != -1 && len(row) > statusIdx {
				item.Status = strings.TrimSpace(row[statusIdx])
			}
			if deadlineIdx != -1 && len(row) > deadlineIdx {
				item.Deadline = strings.TrimSpace(row[deadlineIdx])
			}
			if colorIdx != -1 && len(row) > colorIdx {
				item.ColorDot = strings.TrimSpace(row[colorIdx])
			}

			_ = lineNo
			items = append(items, item)
		}
	} else {
		// Fallback: Parse JSON request
		var jsonReq ImportProjectsJSONRequest
		if err := c.ShouldBindJSON(&jsonReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kirim file CSV multipart atau payload JSON 'projects'"})
			return
		}
		items = jsonReq.Projects
	}

	if len(items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tidak ada baris data proyek valid yang ditemukan untuk diimpor"})
		return
	}

	createdCount := 0
	updatedCount := 0
	var importedProjects []models.Project

	slugRegex := regexp.MustCompile(`[^a-z0-9\-]+`)

	for _, item := range items {
		cleanName := strings.TrimSpace(item.Name)
		if cleanName == "" {
			continue
		}

		slug := strings.TrimSpace(strings.ToLower(item.Slug))
		if slug == "" {
			slug = strings.ToLower(cleanName)
			slug = strings.ReplaceAll(slug, " ", "-")
			slug = slugRegex.ReplaceAllString(slug, "")
		}

		status := strings.ToLower(strings.TrimSpace(item.Status))
		if status != "completed" && status != "on_hold" && status != "archived" {
			status = "active"
		}

		color := strings.TrimSpace(item.ColorDot)
		if color == "" || !strings.HasPrefix(color, "#") {
			color = "#4F46E5"
		}

		deadline := strings.TrimSpace(item.Deadline)
		if deadline == "" {
			deadline = time.Now().AddDate(0, 3, 0).Format("2006-01-02")
		}

		// Check if project with slug already exists (including re-activating soft-deleted if any)
		var existing models.Project
		err := database.DB.Unscoped().Where("slug = ?", slug).First(&existing).Error

		if err == nil {
			// Update existing project
			existing.Name = cleanName
			if item.ClientName != "" {
				existing.ClientName = item.ClientName
			}
			existing.Status = status
			existing.Deadline = deadline
			existing.ColorDot = color
			existing.DeletedAt = gorm.DeletedAt{}
			database.DB.Save(&existing)
			updatedCount++
			importedProjects = append(importedProjects, existing)
		} else {
			// Create new project
			newID := fmt.Sprintf("proj_%s", strings.ReplaceAll(uuid.New().String()[:8], "-", ""))
			newProj := models.Project{
				ID:                 newID,
				Name:               cleanName,
				Slug:               slug,
				ClientName:         item.ClientName,
				Status:             status,
				ProgressPercentage: 0,
				ColorDot:           color,
				Deadline:           deadline,
				OwnerID:            userID,
			}
			database.DB.Create(&newProj)

			// Assign current user as PM member
			if userID != "" {
				member := models.ProjectMember{
					ProjectID:   newProj.ID,
					UserID:      userID,
					ProjectRole: auth.RolePM,
					JoinedAt:    time.Now(),
				}
				database.DB.Create(&member)
			}

			createdCount++
			importedProjects = append(importedProjects, newProj)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       fmt.Sprintf("Impor selesai: %d proyek baru dibuat, %d proyek diperbarui", createdCount, updatedCount),
		"total_rows":    len(items),
		"created_count": createdCount,
		"updated_count": updatedCount,
		"projects":      importedProjects,
	})
}

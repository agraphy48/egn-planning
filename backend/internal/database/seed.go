package database

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"pulse-backend/internal/auth"
	"pulse-backend/internal/models"
)

// SeedDatabase inserts initial test accounts (1 per role), projects, and membership pivot records
func SeedDatabase() {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("[Seed] Database already seeded. Skipping initial seed.")
		return
	}

	defaultPassword := "Password123!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("[Seed] Failed to hash password: %v", err)
	}

	// 1. Create 4 Test Accounts (1 per role)
	users := []models.User{
		{
			ID:        "usr_super_admin",
			Name:      "System Super Admin",
			Email:     "admin@pulse.workspace",
			Password:  string(hashedPassword),
			Role:      auth.RoleSuperAdmin,
			AvatarURL: "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80",
		},
		{
			ID:        "usr_pm",
			Name:      "Alex Rivera (PM)",
			Email:     "pm@pulse.workspace",
			Password:  string(hashedPassword),
			Role:      auth.RolePM,
			AvatarURL: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&auto=format&fit=crop&q=80",
		},
		{
			ID:        "usr_member",
			Name:      "Sarah Connor (Member)",
			Email:     "member@pulse.workspace",
			Password:  string(hashedPassword),
			Role:      auth.RoleMember,
			AvatarURL: "https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=100&auto=format&fit=crop&q=80",
		},
		{
			ID:        "usr_client",
			Name:      "Ir. H. Gunawan (Client)",
			Email:     "client@pulse.workspace",
			Password:  string(hashedPassword),
			Role:      auth.RoleClient,
			AvatarURL: "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100&auto=format&fit=crop&q=80",
		},
	}

	for _, u := range users {
		DB.Create(&u)
	}

	// 2. Create Projects
	p1 := models.Project{
		ID:                 "proj_sustagri",
		Name:               "SUSTAGRI Platform & CMS Modular",
		Slug:               "sustagri",
		ClientName:         "Kementerian Pertanian / Stakeholder Agraria",
		Status:             "active",
		ProgressPercentage: 52,
		ColorDot:           "#10B981",
		Deadline:           "2027-01-26",
		OwnerID:            "usr_pm",
	}

	p2 := models.Project{
		ID:                 "proj_smart",
		Name:               "SMART (Repositori Riset & AI Text Mining)",
		Slug:               "smart",
		ClientName:         "Badan Riset & Inovasi (BRIN) / R&D Institute",
		Status:             "active",
		ProgressPercentage: 43,
		ColorDot:           "#6366F1",
		Deadline:           "2027-01-26",
		OwnerID:            "usr_pm",
	}

	DB.Create(&p1)
	DB.Create(&p2)

	// 3. Project Members Pivot Table (Assign memberships)
	// SUSTAGRI: pm, member, and client
	members := []models.ProjectMember{
		{
			ProjectID:   p1.ID,
			UserID:      "usr_pm",
			ProjectRole: auth.RolePM,
			JoinedAt:    time.Now(),
		},
		{
			ProjectID:   p1.ID,
			UserID:      "usr_member",
			ProjectRole: auth.RoleMember,
			JoinedAt:    time.Now(),
		},
		{
			ProjectID:   p1.ID,
			UserID:      "usr_client",
			ProjectRole: auth.RoleClient,
			JoinedAt:    time.Now(),
		},
		// SMART: pm and member only (client does NOT have access to SMART)
		{
			ProjectID:   p2.ID,
			UserID:      "usr_pm",
			ProjectRole: auth.RolePM,
			JoinedAt:    time.Now(),
		},
		{
			ProjectID:   p2.ID,
			UserID:      "usr_member",
			ProjectRole: auth.RoleMember,
			JoinedAt:    time.Now(),
		},
	}

	for _, m := range members {
		DB.Create(&m)
	}

	// 4. Milestones
	ms1 := models.Milestone{
		ID:          "ms_sust_f1",
		ProjectID:   p1.ID,
		Title:       "Fase 1: CMS Landing Page, Knowledge Base & Fertinnovation CMS",
		Description: "Pengembangan landing page modular, 5 detail topic knowledge, dan modul dinamis Fertinnovation.",
		DueDate:     "2026-12-24",
		Status:      "in_review",
	}
	ms2 := models.Milestone{
		ID:          "ms_smart_f1",
		ProjectID:   p2.ID,
		Title:       "Fase 1: Repository Riset Text Mining, Ekstraksi Data & Dashboard",
		Description: "Pipeline text mining indexing PDF/DOCX, faceted search, engine ekspor multi-format.",
		DueDate:     "2026-11-25",
		Status:      "in_review",
	}
	DB.Create(&ms1)
	DB.Create(&ms2)

	// 5. Initial Tasks
	tasks := []models.Task{
		{
			ID:            "task_01",
			ProjectID:     p1.ID,
			Key:           "SUST-F1-01.1",
			Title:         "Perancangan struktur skema database & arsitektur CMS modular",
			Description:   "Perancangan ERD & Skema Database CMS Sustagri",
			Category:      "Backend",
			Priority:      "High",
			Status:        "Done",
			AssigneeID:    "usr_member",
			DueDate:       "2026-11-05",
			ClientVisible: true,
		},
		{
			ID:            "task_02",
			ProjectID:     p1.ID,
			Key:           "SUST-F1-01.2",
			Title:         "Implementasi Frontend Landing Page Sustagri responsif",
			Description:   "Master Header, Global Search, Notifikasi, Auth UI, Contact Footer",
			Category:      "Frontend",
			Priority:      "High",
			Status:        "Review",
			AssigneeID:    "usr_member",
			DueDate:       "2026-11-12",
			ClientVisible: true,
		},
		{
			ID:            "task_03",
			ProjectID:     p1.ID,
			Key:           "SUST-F1-01.3",
			Title:         "Pembuatan Panel Admin CMS fleksibel",
			Description:   "CRUD dinamis banner, section, menu navigasi, & preview konten",
			Category:      "Backend",
			Priority:      "High",
			Status:        "In Progress",
			AssigneeID:    "usr_pm",
			DueDate:       "2026-11-19",
			ClientVisible: false, // Internal only
		},
		{
			ID:            "task_04",
			ProjectID:     p2.ID,
			Key:           "SMRT-F1-01.1",
			Title:         "Perancangan arsitektur Repositori Hasil Riset & pipeline ekstraksi teks",
			Description:   "PDF/DOCX parsing & indexing text mining",
			Category:      "AI/ML",
			Priority:      "Urgent",
			Status:        "In Progress",
			AssigneeID:    "usr_member",
			DueDate:       "2026-11-06",
			ClientVisible: true,
		},
	}

	for _, t := range tasks {
		DB.Create(&t)
	}

	log.Println("==================================================================")
	log.Println("[Seed] RBAC TEST ACCOUNTS SEEDED SUCCESSFULLY")
	log.Println("Default Password for all accounts: " + defaultPassword)
	log.Println("------------------------------------------------------------------")
	log.Println("1. SUPER ADMIN : admin@pulse.workspace  (Role: super_admin)")
	log.Println("2. PROJECT MGR : pm@pulse.workspace     (Role: pm)")
	log.Println("3. TEAM MEMBER : member@pulse.workspace (Role: member)")
	log.Println("4. CLIENT      : client@pulse.workspace (Role: client)")
	log.Println("==================================================================")
}

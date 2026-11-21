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
		{
			ID:        "usr_agus",
			Name:      "Agus Setiawan",
			Email:     "agus.setiawan@pulse.workspace",
			Password:  string(hashedPassword),
			Role:      auth.RolePM,
			AvatarURL: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&auto=format&fit=crop&q=80",
		},
		{
			ID:        "usr_adi",
			Name:      "Adi Saputra",
			Email:     "adi.saputra@pulse.workspace",
			Password:  string(hashedPassword),
			Role:      auth.RolePM,
			AvatarURL: "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=100&auto=format&fit=crop&q=80",
		},
		{
			ID:        "usr_yudistira",
			Name:      "Yudistira Sira Permana",
			Email:     "yudistira@pulse.workspace",
			Password:  string(hashedPassword),
			Role:      auth.RoleMember,
			AvatarURL: "https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?w=100&auto=format&fit=crop&q=80",
		},
		{
			ID:        "usr_hardianto",
			Name:      "Hardianto",
			Email:     "hardianto@pulse.workspace",
			Password:  string(hashedPassword),
			Role:      auth.RoleMember,
			AvatarURL: "https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=100&auto=format&fit=crop&q=80",
		},
		{
			ID:        "usr_dedi",
			Name:      "Dedi Irawan",
			Email:     "dedi.irawan@pulse.workspace",
			Password:  string(hashedPassword),
			Role:      auth.RoleMember,
			AvatarURL: "https://images.unsplash.com/photo-1522075469751-3a6694fb2f61?w=100&auto=format&fit=crop&q=80",
		},
	}

	for _, u := range users {
		DB.Create(&u)
	}

	// 2. Create Projects
	p1 := models.Project{
		ID:                 "proj_sustagri",
		Name:               "SUSTAGRI Fase 1",
		Slug:               "sustagri-fase-1",
		ClientName:         "Kementerian Pertanian / Stakeholder Agraria",
		Status:             "active",
		ProgressPercentage: 52,
		ColorDot:           "#10B981",
		Deadline:           "2027-01-26",
		OwnerID:            "usr_pm",
	}

	p1_2 := models.Project{
		ID:                 "proj_sustagri_f2",
		Name:               "SUSTAGRI Fase 2",
		Slug:               "sustagri-fase-2",
		ClientName:         "Kementerian Pertanian / Stakeholder Agraria",
		Status:             "active",
		ProgressPercentage: 25,
		ColorDot:           "#059669",
		Deadline:           "2027-04-15",
		OwnerID:            "usr_pm",
	}

	p2 := models.Project{
		ID:                 "proj_smart",
		Name:               "SMART Fase 1",
		Slug:               "smart-fase-1",
		ClientName:         "Badan Riset & Inovasi (BRIN) / R&D Institute",
		Status:             "active",
		ProgressPercentage: 43,
		ColorDot:           "#6366F1",
		Deadline:           "2027-01-26",
		OwnerID:            "usr_pm",
	}

	p2_2 := models.Project{
		ID:                 "proj_smart_f2",
		Name:               "SMART Fase 2",
		Slug:               "smart-fase-2",
		ClientName:         "Badan Riset & Inovasi (BRIN) / R&D Institute",
		Status:             "active",
		ProgressPercentage: 20,
		ColorDot:           "#8B5CF6",
		Deadline:           "2027-05-20",
		OwnerID:            "usr_pm",
	}

	DB.Create(&p1)
	DB.Create(&p1_2)
	DB.Create(&p2)
	DB.Create(&p2_2)

	// 3. Project Members Pivot Table (Assign memberships for all 4 projects)
	targetProjects := []string{p1.ID, p1_2.ID, p2.ID, p2_2.ID}
	assignedUsers := []struct {
		userID string
		role   string
	}{
		{"usr_pm", auth.RolePM},
		{"usr_member", auth.RoleMember},
		{"usr_client", auth.RoleClient},
		{"usr_agus", auth.RolePM},
		{"usr_adi", auth.RolePM},
		{"usr_yudistira", auth.RoleMember},
		{"usr_hardianto", auth.RoleMember},
		{"usr_dedi", auth.RoleMember},
	}

	for _, pid := range targetProjects {
		for _, u := range assignedUsers {
			DB.Create(&models.ProjectMember{
				ProjectID:   pid,
				UserID:      u.userID,
				ProjectRole: u.role,
				JoinedAt:    time.Now(),
			})
		}
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
	ms1_2 := models.Milestone{
		ID:          "ms_sust_f2",
		ProjectID:   p1_2.ID,
		Title:       "Fase 2: Integrasi IoT Sensor Tanah, Marketplace Pupuk & Mobile App",
		Description: "Implementasi telemetri MQTT sensor tanah, modul distribusi pupuk, dan mobile app Flutter.",
		DueDate:     "2027-04-15",
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
	ms2_2 := models.Milestone{
		ID:          "ms_smart_f2",
		ProjectID:   p2_2.ID,
		Title:       "Fase 2: AI Multi-Agent Semantic Search & Deep Graph Analytics",
		Description: "Vector search embeddings, ekstraksi abstrak otomatis, dan visualisasi relasi graph riset.",
		DueDate:     "2027-05-20",
		Status:      "in_review",
	}
	DB.Create(&ms1)
	DB.Create(&ms1_2)
	DB.Create(&ms2)
	DB.Create(&ms2_2)

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

// SeedAuxiliaryData populates default records for new modules if empty without modifying existing projects or users
func SeedAuxiliaryData() {
	// 1. Activity Logs
	var actCount int64
	DB.Model(&models.ActivityLog{}).Count(&actCount)
	if actCount == 0 {
		logs := []models.ActivityLog{
			{
				ID:         "act_01",
				ProjectID:  "proj_sustagri",
				UserID:     "usr_agus",
				UserName:   "Agus Setiawan",
				UserAvatar: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&auto=format&fit=crop&q=80",
				Action:     "membuat task baru",
				Target:     "SUST-F1-01.1",
				Details:    "Inisialisasi arsitektur database",
				CreatedAt:  time.Now().Add(-2 * time.Hour),
			},
			{
				ID:         "act_02",
				ProjectID:  "proj_sustagri",
				UserID:     "usr_yudistira",
				UserName:   "Yudistira Sira Permana",
				UserAvatar: "https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?w=100&auto=format&fit=crop&q=80",
				Action:     "menyelesaikan code review",
				Target:     "SUST-F1-01.2",
				Details:    "Frontend Landing Page Sustagri responsif",
				CreatedAt:  time.Now().Add(-1 * time.Hour),
			},
			{
				ID:         "act_03",
				ProjectID:  "proj_smart",
				UserID:     "usr_hardianto",
				UserName:   "Hardianto",
				UserAvatar: "https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=100&auto=format&fit=crop&q=80",
				Action:     "mengunggah deliverable",
				Target:     "Milestone 1",
				Details:    "Repositori Hasil Riset & UI",
				CreatedAt:  time.Now().Add(-30 * time.Minute),
			},
		}
		for _, l := range logs {
			DB.Create(&l)
		}
	}

	// 2. Wiki Pages
	var wikiCount int64
	DB.Model(&models.WikiPage{}).Count(&wikiCount)
	if wikiCount == 0 {
		wikis := []models.WikiPage{
			{
				ID:            "wiki_01",
				ProjectID:     "proj_sustagri",
				Title:         "Spesifikasi Arsitektur Sistem & API RBAC",
				Slug:          "spesifikasi-arsitektur",
				Content:       "Aplikasi SustAgri menggunakan backend Golang Gin dengan enkripsi JWT HS256 dan RBAC 4 role (Super Admin, PM, Member, Client). Keamanan jaringan diperkuat dengan Cloudflare Tunnel TLS 1.3.",
				AuthorID:      "usr_agus",
				AuthorName:    "Agus Setiawan",
				ClientVisible: true,
				CreatedAt:     time.Now().Add(-24 * time.Hour),
				UpdatedAt:     time.Now().Add(-24 * time.Hour),
			},
			{
				ID:            "wiki_02",
				ProjectID:     "proj_sustagri",
				Title:         "Panduan Deployment Windows Server Host",
				Slug:          "panduan-deploy",
				Content:       "Daemon produksi berjalan pada host Windows lokal (Node.js gateway port 5080 + Go binary port 5081) dengan auto-restart dan logging log rotasi.",
				AuthorID:      "usr_dedi",
				AuthorName:    "Dedi Irawan",
				ClientVisible: false,
				CreatedAt:     time.Now().Add(-12 * time.Hour),
				UpdatedAt:     time.Now().Add(-12 * time.Hour),
			},
		}
		for _, w := range wikis {
			DB.Create(&w)
		}
	}

	// 3. Time Entries (Timesheets)
	var timeCount int64
	DB.Model(&models.TimeEntry{}).Count(&timeCount)
	if timeCount == 0 {
		times := []models.TimeEntry{
			{
				ID:          "time_01",
				ProjectID:   "proj_sustagri",
				TaskID:      "task_01",
				UserID:      "usr_yudistira",
				Hours:       6.5,
				Description: "Implementasi skema database SQLite dan migrasi GORM",
				Date:        time.Now().Format("2006-01-02"),
				CreatedAt:   time.Now().Add(-4 * time.Hour),
			},
			{
				ID:          "time_02",
				ProjectID:   "proj_sustagri",
				TaskID:      "task_02",
				UserID:      "usr_hardianto",
				Hours:       8.0,
				Description: "Pembuatan komponen UI Landing Page Vue 3 dan Tailwind",
				Date:        time.Now().Format("2006-01-02"),
				CreatedAt:   time.Now().Add(-3 * time.Hour),
			},
			{
				ID:          "time_03",
				ProjectID:   "proj_smart",
				TaskID:      "task_04",
				UserID:      "usr_dedi",
				Hours:       5.0,
				Description: "Pengujian pipeline ekstraksi teks dan verifikasi format",
				Date:        time.Now().Format("2006-01-02"),
				CreatedAt:   time.Now().Add(-2 * time.Hour),
			},
		}
		for _, te := range times {
			DB.Create(&te)
		}
	}

	// 4. Notifications
	var notifCount int64
	DB.Model(&models.Notification{}).Count(&notifCount)
	if notifCount == 0 {
		notifs := []models.Notification{
			{
				ID:        "notif_01",
				UserID:    "usr_agus",
				Title:     "Milestone Memerlukan Review",
				Message:   "Milestone 'Inisiasi & Analisis Kebutuhan' telah diserahkan untuk review klien.",
				Type:      "milestone",
				Link:      "/client-portal",
				IsRead:    false,
				CreatedAt: time.Now().Add(-1 * time.Hour),
			},
			{
				ID:        "notif_02",
				UserID:    "usr_agus",
				Title:     "Pencatatan Jam Kerja",
				Message:   "Yudistira mencatatkan 6.5 jam kerja pada task SUST-F1-01.1.",
				Type:      "timesheet",
				Link:      "/projects",
				IsRead:    false,
				CreatedAt: time.Now().Add(-2 * time.Hour),
			},
			{
				ID:        "notif_03",
				UserID:    "usr_client",
				Title:     "Deliverable Baru Siap Ditinjau",
				Message:   "Deliverable SustAgri Sprint 1 siap untuk verifikasi dan persetujuan Anda.",
				Type:      "deliverable",
				Link:      "/client-portal",
				IsRead:    false,
				CreatedAt: time.Now().Add(-30 * time.Minute),
			},
		}
		for _, n := range notifs {
			DB.Create(&n)
		}
	}

	// 5. Budgets
	var budgetCount int64
	DB.Model(&models.Budget{}).Count(&budgetCount)
	if budgetCount == 0 {
		budgets := []models.Budget{
			{
				ID:        "bgt_01",
				ProjectID: "proj_sustagri",
				Category:  "Engineering & Development",
				Allocated: 50000000,
				Spent:     32500000,
				Currency:  "IDR",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "bgt_02",
				ProjectID: "proj_sustagri",
				Category:  "Cloud & Infrastructure (Cloudflare/Host)",
				Allocated: 15000000,
				Spent:     4200000,
				Currency:  "IDR",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "bgt_03",
				ProjectID: "proj_sustagri",
				Category:  "UI/UX & Prototyping",
				Allocated: 10000000,
				Spent:     9000000,
				Currency:  "IDR",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "bgt_04",
				ProjectID: "proj_smart",
				Category:  "AI Model Fine-Tuning",
				Allocated: 40000000,
				Spent:     21000000,
				Currency:  "IDR",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		for _, b := range budgets {
			DB.Create(&b)
		}
	}
}

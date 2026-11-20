package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a system user with an explicit global Role
type User struct {
	ID        string         `gorm:"primaryKey;size:64" json:"id"`
	Name      string         `gorm:"size:128;not null" json:"name"`
	Email     string         `gorm:"size:128;uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Role      string         `gorm:"size:32;not null;default:'member'" json:"role"` // super_admin, pm, member, client
	AvatarURL string         `gorm:"size:255" json:"avatar_url"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ProjectMember is the pivot table enforcing project-level authorization
type ProjectMember struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID   string    `gorm:"size:64;index;uniqueIndex:idx_project_user,priority:1;not null" json:"project_id"`
	UserID      string    `gorm:"size:64;index;uniqueIndex:idx_project_user,priority:2;not null" json:"user_id"`
	ProjectRole string    `gorm:"size:32;not null;default:'member'" json:"project_role"` // pm, member, client
	JoinedAt    time.Time `json:"joined_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// Project represents a workspace project
type Project struct {
	ID                 string         `gorm:"primaryKey;size:64" json:"id"`
	Name               string         `gorm:"size:255;not null" json:"name"`
	Slug               string         `gorm:"size:255;uniqueIndex" json:"slug"`
	ClientName         string         `gorm:"size:255" json:"client_name,omitempty"`
	Status             string         `gorm:"size:32;default:'active'" json:"status"` // active, completed, archived
	ProgressPercentage int            `gorm:"default:0" json:"progress_percentage"`
	ColorDot           string         `gorm:"size:32;default:'#4F46E5'" json:"color_dot"`
	Deadline           string         `gorm:"size:64" json:"deadline"`
	OwnerID            string         `gorm:"size:64;not null" json:"owner_id"`
	TotalTasks         int            `gorm:"-" json:"total_tasks"`
	CompletedTasks     int            `gorm:"-" json:"completed_tasks"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Members    []ProjectMember `gorm:"foreignKey:ProjectID" json:"members,omitempty"`
	Milestones []Milestone     `gorm:"foreignKey:ProjectID" json:"milestones,omitempty"`
	Tasks      []Task          `gorm:"foreignKey:ProjectID" json:"tasks,omitempty"`
}

// Milestone represents formal contractual gates & deliverables
type Milestone struct {
	ID          string     `gorm:"primaryKey;size:64" json:"id"`
	ProjectID   string     `gorm:"size:64;index;not null" json:"project_id"`
	Title       string     `gorm:"size:255;not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	DueDate     string     `gorm:"size:64" json:"due_date"`
	Status      string     `gorm:"size:32;default:'pending'" json:"status"` // pending, in_review, approved, rejected
	ApprovedBy  string     `gorm:"size:128" json:"approved_by,omitempty"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Task represents a kanban task
type Task struct {
	ID            string         `gorm:"primaryKey;size:64" json:"id"`
	ProjectID     string         `gorm:"size:64;index;not null" json:"project_id"`
	Key           string         `gorm:"size:64;index" json:"key"` // e.g. #SUST-F1-01.1
	Title         string         `gorm:"size:255;not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	Category      string         `gorm:"size:64;default:'General'" json:"category"`
	Priority      string         `gorm:"size:32;default:'Medium'" json:"priority"` // Urgent, High, Medium, Low
	Status        string         `gorm:"size:32;default:'Backlog'" json:"status"`  // Backlog, In Progress, Review, Done
	AssigneeID    string         `gorm:"size:64;index" json:"assignee_id"`
	DueDate       string         `gorm:"size:64" json:"due_date"`
	StartDate     string         `gorm:"size:64" json:"start_date,omitempty"`
	TimeEstimate  string         `gorm:"size:32" json:"time_estimate,omitempty"`
	Fase          string         `gorm:"size:64" json:"fase,omitempty"`
	EpicModul     string         `gorm:"size:255" json:"epic_modul,omitempty"`
	Deliverable   string         `gorm:"size:255" json:"deliverable,omitempty"`
	ClientVisible bool           `gorm:"default:true" json:"client_visible"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Assignee *User     `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
	Subtasks []SubTask `gorm:"foreignKey:TaskID" json:"subtasks"`
	Comments []Comment `gorm:"foreignKey:TaskID" json:"comments"`
}

// SubTask represents an interactive checklist item
type SubTask struct {
	ID        string    `gorm:"primaryKey;size:64" json:"id"`
	TaskID    string    `gorm:"size:64;index;not null" json:"task_id"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	IsDone    bool      `gorm:"default:false" json:"is_done"`
	CreatedAt time.Time `json:"created_at"`
}

// Comment represents team and client discussion
type Comment struct {
	ID        string    `gorm:"primaryKey;size:64" json:"id"`
	TaskID    string    `gorm:"size:64;index;not null" json:"task_id"`
	UserID    string    `gorm:"size:64;index;not null" json:"user_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`

	// Relation
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

package auth

// Role constants
const (
	RoleSuperAdmin = "super_admin"
	RolePM         = "pm"
	RoleMember     = "member"
	RoleClient     = "client"
)

// Permission constants
const (
	PermUserManage          = "user:manage"
	PermUserRoleUpdate      = "user:role:update"
	PermProjectCreate       = "project:create"
	PermProjectRead         = "project:read"
	PermProjectEdit         = "project:edit"
	PermProjectDelete       = "project:delete"
	PermProjectArchive      = "project:archive"
	PermProjectMemberManage = "project:member:manage"
	PermMilestoneCreate     = "milestone:create"
	PermMilestoneRead       = "milestone:read"
	PermMilestoneEdit       = "milestone:edit"
	PermMilestoneApprove    = "milestone:approve"
	PermTaskCreate          = "task:create"
	PermTaskRead            = "task:read"
	PermTaskEdit            = "task:edit"
	PermTaskMove            = "task:move"
	PermTaskDelete          = "task:delete"
	PermSubtaskManage       = "subtask:manage"
	PermAttachmentUpload    = "attachment:upload"
	PermCommentCreate       = "comment:create"
	PermCommentRead         = "comment:read"
)

// AllPermissions contains all valid permissions in the system
var AllPermissions = []string{
	PermUserManage,
	PermUserRoleUpdate,
	PermProjectCreate,
	PermProjectRead,
	PermProjectEdit,
	PermProjectDelete,
	PermProjectArchive,
	PermProjectMemberManage,
	PermMilestoneCreate,
	PermMilestoneRead,
	PermMilestoneEdit,
	PermMilestoneApprove,
	PermTaskCreate,
	PermTaskRead,
	PermTaskEdit,
	PermTaskMove,
	PermTaskDelete,
	PermSubtaskManage,
	PermAttachmentUpload,
	PermCommentCreate,
	PermCommentRead,
}

// RolePermissions defines the strict permission map for each role
var RolePermissions = map[string][]string{
	RoleSuperAdmin: AllPermissions,
	RolePM: {
		PermProjectCreate,
		PermProjectRead,
		PermProjectEdit,
		PermProjectArchive,
		PermProjectMemberManage,
		PermMilestoneCreate,
		PermMilestoneRead,
		PermMilestoneEdit,
		PermMilestoneApprove,
		PermTaskCreate,
		PermTaskRead,
		PermTaskEdit,
		PermTaskMove,
		PermTaskDelete,
		PermSubtaskManage,
		PermAttachmentUpload,
		PermCommentCreate,
		PermCommentRead,
	},
	RoleMember: {
		PermProjectRead,
		PermMilestoneRead,
		PermTaskCreate,
		PermTaskRead,
		PermTaskEdit,
		PermTaskMove,
		PermSubtaskManage,
		PermAttachmentUpload,
		PermCommentCreate,
		PermCommentRead,
	},
	RoleClient: {
		PermProjectRead,
		PermMilestoneRead,
		PermMilestoneApprove,
		PermTaskRead,
		PermCommentCreate,
		PermCommentRead,
	},
}

// GetPermissionsForRole returns the flat array of permissions for a role
func GetPermissionsForRole(role string) []string {
	perms, exists := RolePermissions[role]
	if !exists {
		return []string{}
	}
	return perms
}

// HasPermission checks if a role contains a specific permission
func HasPermission(role string, permission string) bool {
	if role == RoleSuperAdmin {
		return true
	}
	perms := GetPermissionsForRole(role)
	for _, p := range perms {
		if p == permission {
			return true
		}
	}
	return false
}

// IsValidRole validates if a role string exists in the system
func IsValidRole(role string) bool {
	switch role {
	case RoleSuperAdmin, RolePM, RoleMember, RoleClient:
		return true
	default:
		return false
	}
}

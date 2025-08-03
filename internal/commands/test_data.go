package commands

import (
	"github.com/lunchboxsushi/jit/pkg/types"
)

// GenerateTestTickets creates dummy tickets for testing the log command
func GenerateTestTickets() []*types.Ticket {
	return []*types.Ticket{
		// Epic 1: User Authentication
		{
			Key:         "PROJ-100",
			Type:        types.TicketTypeEpic,
			Title:       "User Authentication",
			Status:      "In Progress",
			Description: "Modernize our authentication system",
		},
		{
			Key:         "PROJ-101",
			Type:        types.TicketTypeTask,
			Title:       "OAuth Implementation",
			Status:      "To Do",
			Description: "Implement OAuth 2.0 providers",
			Relationships: types.TicketRelationships{
				ParentKey: "PROJ-100",
			},
		},
		{
			Key:         "PROJ-102",
			Type:        types.TicketTypeSubtask,
			Title:       "Google OAuth",
			Status:      "Done",
			Description: "Integrate Google OAuth provider",
			Relationships: types.TicketRelationships{
				ParentKey: "PROJ-101",
			},
		},
		{
			Key:         "PROJ-103",
			Type:        types.TicketTypeSubtask,
			Title:       "GitHub OAuth",
			Status:      "In Progress",
			Description: "Integrate GitHub OAuth provider",
			Relationships: types.TicketRelationships{
				ParentKey: "PROJ-101",
			},
		},
		{
			Key:         "PROJ-104",
			Type:        types.TicketTypeTask,
			Title:       "MFA Setup",
			Status:      "Blocked",
			Description: "Implement multi-factor authentication",
			Relationships: types.TicketRelationships{
				ParentKey: "PROJ-100",
			},
		},
		{
			Key:         "PROJ-105",
			Type:        types.TicketTypeSubtask,
			Title:       "TOTP Implementation",
			Status:      "To Do",
			Description: "Implement TOTP-based MFA",
			Relationships: types.TicketRelationships{
				ParentKey: "PROJ-104",
			},
		},

		// Epic 2: Database Migration
		{
			Key:         "PROJ-200",
			Type:        types.TicketTypeEpic,
			Title:       "Database Migration",
			Status:      "Done",
			Description: "Migrate to new database schema",
		},
		{
			Key:         "PROJ-201",
			Type:        types.TicketTypeTask,
			Title:       "Schema Updates",
			Status:      "In Progress",
			Description: "Update database schema",
			Relationships: types.TicketRelationships{
				ParentKey: "PROJ-200",
			},
		},

		// Orphaned Tasks
		{
			Key:         "PROJ-300",
			Type:        types.TicketTypeTask,
			Title:       "Standalone Feature",
			Status:      "To Do",
			Description: "A task without a parent epic",
		},
		{
			Key:         "PROJ-301",
			Type:        types.TicketTypeSubtask,
			Title:       "Implementation Details",
			Status:      "In Progress",
			Description: "Details for standalone feature",
			Relationships: types.TicketRelationships{
				ParentKey: "PROJ-300",
			},
		},
		{
			Key:         "PROJ-302",
			Type:        types.TicketTypeTask,
			Title:       "Another Orphan Task",
			Status:      "Done",
			Description: "Another task without parent",
		},
	}
}

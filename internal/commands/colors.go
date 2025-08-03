package commands

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lunchboxsushi/jit/pkg/types"
)

// Color scheme for the log command
var (
	// Ticket type colors
	EpicColor    = lipgloss.NewStyle().Foreground(lipgloss.Color("#8B5CF6"))            // Purple
	TaskColor    = lipgloss.NewStyle().Foreground(lipgloss.Color("#3B82F6"))            // Blue
	SubtaskColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#60A5FA"))            // Light Blue
	FocusColor   = lipgloss.NewStyle().Foreground(lipgloss.Color("#F97316"))            // Orange
	HeaderColor  = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Bold(true) // Gray, Bold

	// Status colors
	StatusDone       = lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981")) // Green
	StatusInProgress = lipgloss.NewStyle().Foreground(lipgloss.Color("#F59E0B")) // Yellow
	StatusBlocked    = lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444")) // Red
	StatusToDo       = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")) // Gray

	// Tree structure colors
	TreeColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")) // Gray
)

// GetStatusColor returns the appropriate color for a status
func GetStatusColor(status string) lipgloss.Style {
	switch strings.ToLower(status) {
	case "done", "completed", "closed":
		return StatusDone
	case "in progress", "in-progress", "progress":
		return StatusInProgress
	case "blocked", "block":
		return StatusBlocked
	default:
		return StatusToDo
	}
}

// GetTicketTypeColor returns the appropriate color for a ticket type
func GetTicketTypeColor(ticketType string) lipgloss.Style {
	switch ticketType {
	case "epic":
		return EpicColor
	case "task":
		return TaskColor
	case "subtask":
		return SubtaskColor
	default:
		return TaskColor
	}
}

// ColorizeTicket formats a ticket with appropriate colors
func ColorizeTicket(ticket *types.Ticket, isFocused bool, showOrphan bool) string {
	var parts []string

	// Focus indicator
	if isFocused {
		parts = append(parts, FocusColor.Render("@"))
	} else {
		parts = append(parts, " ")
	}

	// Ticket type
	typeColor := GetTicketTypeColor(ticket.Type)
	parts = append(parts, typeColor.Render(strings.Title(ticket.Type)))

	// Ticket key
	parts = append(parts, fmt.Sprintf("[%s]", ticket.Key))

	// Status
	statusColor := GetStatusColor(ticket.Status)
	parts = append(parts, statusColor.Render(fmt.Sprintf("<%s>", ticket.Status)))

	// Title
	parts = append(parts, "-")
	parts = append(parts, ticket.Title)

	return strings.Join(parts, " ")
}

// ColorizeHeader formats a section header
func ColorizeHeader(header string) string {
	return HeaderColor.Render(header)
}

// ColorizeTreeLine formats a tree line with proper indentation
func ColorizeTreeLine(prefix, connector string, isLast bool) string {
	if isLast {
		return TreeColor.Render(prefix + "└─")
	}
	return TreeColor.Render(prefix + "├─")
}

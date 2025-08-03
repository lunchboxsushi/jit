package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/lunchboxsushi/jit/pkg/types"
	"github.com/spf13/cobra"
)

var (
	logAllFlag    bool
	logStatusFlag string
	logJSONFlag   bool
	logOrphanFlag bool
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Display ticket tree view",
	Long: `Show a hierarchical view of tracked tickets with current focus highlighted.
	
Examples:
  jit log                    # Show current context tree
  jit log --all             # Show all tracked tickets
  jit log --status "In Progress"  # Filter by status
  jit log --json            # Output as JSON`,
	Run: func(cmd *cobra.Command, args []string) {
		// For testing, use dummy data
		tickets := GenerateTestTickets()

		// Set current focus for testing (only one should be active)
		currentEpic := ""
		currentTask := ""
		currentSubtask := "PROJ-103"

		if len(tickets) == 0 {
			if logJSONFlag {
				fmt.Println("{\n  \"current_focus\": {\n    \"epic\": \"\",\n    \"task\": \"\",\n    \"subtask\": \"\"\n  },\n  \"tickets\": []\n}")
			} else {
				fmt.Println("No tickets found. Use 'jit track <ticket>' to start tracking tickets.")
			}
			return
		}

		// Filter by status if specified
		if logStatusFlag != "" {
			tickets = filterTicketsByStatus(tickets, logStatusFlag)
		}

		// Sort tickets by type and key
		sortTickets(tickets)

		// Display based on format
		if logJSONFlag {
			displayJSON(tickets, currentEpic, currentTask, currentSubtask)
		} else {
			displayTree(tickets, currentEpic, currentTask, currentSubtask)
		}
	},
}

func init() {
	logCmd.Flags().BoolVar(&logAllFlag, "all", false, "Show all tickets (not just current context)")
	logCmd.Flags().StringVar(&logStatusFlag, "status", "", "Filter by status (e.g., 'In Progress', 'Done')")
	logCmd.Flags().BoolVar(&logJSONFlag, "json", false, "Output as JSON")
	logCmd.Flags().BoolVar(&logOrphanFlag, "orphan", false, "Show orphaned tasks")
}

// filterTicketsByStatus filters tickets by status
func filterTicketsByStatus(tickets []*types.Ticket, status string) []*types.Ticket {
	var filtered []*types.Ticket
	for _, ticket := range tickets {
		if strings.EqualFold(ticket.Status, status) {
			filtered = append(filtered, ticket)
		}
	}
	return filtered
}

// sortTickets sorts tickets by type (epic, task, subtask) and then by key
func sortTickets(tickets []*types.Ticket) {
	sort.Slice(tickets, func(i, j int) bool {
		// Sort by type priority: epic > task > subtask
		typeOrder := map[string]int{
			types.TicketTypeEpic:    1,
			types.TicketTypeTask:    2,
			types.TicketTypeSubtask: 3,
		}

		orderI := typeOrder[tickets[i].Type]
		orderJ := typeOrder[tickets[j].Type]

		if orderI != orderJ {
			return orderI < orderJ
		}

		// If same type, sort by key
		return tickets[i].Key < tickets[j].Key
	})
}

// displayTree displays tickets in the new enhanced tree format
func displayTree(tickets []*types.Ticket, currentEpic, currentTask, currentSubtask string) {
	// Ensure only one focus item is set (prioritize most specific: subtask > task > epic)
	var actualFocus string
	var focusType string

	if currentSubtask != "" {
		actualFocus = currentSubtask
		focusType = "subtask"
	} else if currentTask != "" {
		actualFocus = currentTask
		focusType = "task"
	} else if currentEpic != "" {
		actualFocus = currentEpic
		focusType = "epic"
	}

	// Group tickets by type
	epics := filterTicketsByType(tickets, types.TicketTypeEpic)
	tasks := filterTicketsByType(tickets, types.TicketTypeTask)
	subtasks := filterTicketsByType(tickets, types.TicketTypeSubtask)

	// Display epics and their hierarchy with custom tree formatting
	for i, epic := range epics {
		// Epic is focused only if it's the actual focus item
		isEpicFocused := focusType == "epic" && epic.Key == actualFocus

		epicText := formatTicketText(epic, isEpicFocused)
		fmt.Println(epicText)

		// Show tasks under this epic
		epicTasks := filterTasksByParent(tasks, epic.Key)
		for j, task := range epicTasks {
			// Task is focused only if it's the actual focus item
			isTaskFocused := focusType == "task" && task.Key == actualFocus

			taskText := formatTicketText(task, isTaskFocused)
			isLastTask := j == len(epicTasks)-1

			// Use rounded tree characters
			if isLastTask {
				fmt.Printf("   ╰─── %s\n", taskText)
			} else {
				fmt.Printf("   ├─── %s\n", taskText)
			}

			// Check if task has subtasks
			taskSubtasks := filterSubtasksByParent(subtasks, task.Key)
			if len(taskSubtasks) > 0 {
				for k, subtask := range taskSubtasks {
					// Subtask is focused only if it's the actual focus item
					isSubtaskFocused := focusType == "subtask" && subtask.Key == actualFocus
					subtaskText := formatTicketText(subtask, isSubtaskFocused)
					isLastSubtask := k == len(taskSubtasks)-1

					// Use rounded tree characters with proper indentation
					if isLastTask {
						if isLastSubtask {
							fmt.Printf("       ╰─── %s\n", subtaskText)
						} else {
							fmt.Printf("       ├─── %s\n", subtaskText)
						}
					} else {
						if isLastSubtask {
							fmt.Printf("   │   ╰─── %s\n", subtaskText)
						} else {
							fmt.Printf("   │   ├─── %s\n", subtaskText)
						}
					}
				}
			}
		}

		// Add spacing between epics
		if i < len(epics)-1 {
			fmt.Println()
		}
	}

	// Show orphaned tasks if there are any
	orphanTasks := filterOrphanTasks(tasks)
	if len(orphanTasks) > 0 && (logOrphanFlag || logAllFlag) {
		fmt.Println()
		fmt.Println("== ORPHAN TASKS ==")
		fmt.Println()

		for i, task := range orphanTasks {
			// Orphan task is focused only if it's the actual focus item
			isTaskFocused := focusType == "task" && task.Key == actualFocus
			taskText := formatTicketText(task, isTaskFocused)

			// Check if orphan task has subtasks
			taskSubtasks := filterSubtasksByParent(subtasks, task.Key)
			if len(taskSubtasks) > 0 {
				fmt.Printf("   ├─── %s\n", taskText)

				// Add subtasks
				for j, subtask := range taskSubtasks {
					// Subtask is focused only if it's the actual focus item
					isSubtaskFocused := focusType == "subtask" && subtask.Key == actualFocus
					subtaskText := formatTicketText(subtask, isSubtaskFocused)
					isLastSubtask := j == len(taskSubtasks)-1

					if isLastSubtask {
						fmt.Printf("   │   ╰─── %s\n", subtaskText)
					} else {
						fmt.Printf("   │   ├─── %s\n", subtaskText)
					}
				}
			} else {
				if i == len(orphanTasks)-1 {
					fmt.Printf("   ╰─── %s\n", taskText)
				} else {
					fmt.Printf("   ├─── %s\n", taskText)
				}
			}
		}
	}
}

// formatTicketText formats a ticket with colors and focus indicators
func formatTicketText(ticket *types.Ticket, isFocused bool) string {
	var parts []string

	// Focus indicator - only show @ for the actual focus item and its parent chain
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

// displayJSON displays tickets in JSON format
func displayJSON(tickets []*types.Ticket, currentEpic, currentTask, currentSubtask string) {
	fmt.Println("{")
	fmt.Printf("  \"current_focus\": {\n")
	fmt.Printf("    \"epic\": \"%s\",\n", currentEpic)
	fmt.Printf("    \"task\": \"%s\",\n", currentTask)
	fmt.Printf("    \"subtask\": \"%s\"\n", currentSubtask)
	fmt.Printf("  },\n")
	fmt.Printf("  \"tickets\": [\n")

	for i, ticket := range tickets {
		fmt.Printf("    {\n")
		fmt.Printf("      \"key\": \"%s\",\n", ticket.Key)
		fmt.Printf("      \"type\": \"%s\",\n", ticket.Type)
		fmt.Printf("      \"title\": \"%s\",\n", ticket.Title)
		fmt.Printf("      \"status\": \"%s\",\n", ticket.Status)
		fmt.Printf("      \"description\": \"%s\"", ticket.Description)

		if ticket.Relationships.ParentKey != "" {
			fmt.Printf(",\n      \"parent_key\": \"%s\"", ticket.Relationships.ParentKey)
		}

		if i < len(tickets)-1 {
			fmt.Printf("\n    },\n")
		} else {
			fmt.Printf("\n    }\n")
		}
	}

	fmt.Printf("  ]\n")
	fmt.Printf("}\n")
}

// filterTicketsByType filters tickets by type
func filterTicketsByType(tickets []*types.Ticket, ticketType string) []*types.Ticket {
	var filtered []*types.Ticket
	for _, ticket := range tickets {
		if ticket.Type == ticketType {
			filtered = append(filtered, ticket)
		}
	}
	return filtered
}

// filterTasksByParent filters tasks by parent epic
func filterTasksByParent(tasks []*types.Ticket, parentEpic string) []*types.Ticket {
	var filtered []*types.Ticket
	for _, task := range tasks {
		if task.Relationships.ParentKey == parentEpic {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

// filterSubtasksByParent filters subtasks by parent task
func filterSubtasksByParent(subtasks []*types.Ticket, parentTask string) []*types.Ticket {
	var filtered []*types.Ticket
	for _, subtask := range subtasks {
		if subtask.Relationships.ParentKey == parentTask {
			filtered = append(filtered, subtask)
		}
	}
	return filtered
}

// filterOrphanTasks filters tasks that have no parent (orphan tasks)
func filterOrphanTasks(tasks []*types.Ticket) []*types.Ticket {
	var filtered []*types.Ticket
	for _, task := range tasks {
		if task.Relationships.ParentKey == "" {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

// hasTaskWithParent checks if any task under a given epic has the specified task key
func hasTaskWithParent(tasks []*types.Ticket, epicKey string, taskKey string) bool {
	for _, task := range tasks {
		if task.Relationships.ParentKey == epicKey && task.Key == taskKey {
			return true
		}
	}
	return false
}

// hasSubtaskWithParent checks if any subtask under a given task has the specified subtask key
func hasSubtaskWithParent(subtasks []*types.Ticket, taskKey string, subtaskKey string) bool {
	for _, subtask := range subtasks {
		if subtask.Relationships.ParentKey == taskKey && subtask.Key == subtaskKey {
			return true
		}
	}
	return false
}

// GetLogCmd returns the log command
func GetLogCmd() *cobra.Command {
	return logCmd
}

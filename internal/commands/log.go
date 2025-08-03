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
		// Initialize command context
		ctx, err := InitializeCommand()
		if err != nil {
			HandleError(err, "Failed to initialize")
			return
		}

		// Get current focus
		currentEpic, _ := ctx.ContextManager.GetCurrentEpic()
		currentTask, _ := ctx.ContextManager.GetCurrentTask()
		currentSubtask, _ := ctx.ContextManager.GetCurrentSubtask()

		// Get all tickets
		ticketKeys, err := ctx.Storage.ListTickets()
		if err != nil {
			HandleError(err, "Failed to list tickets")
			return
		}

		if len(ticketKeys) == 0 {
			if logJSONFlag {
				fmt.Println("{\n  \"current_focus\": {\n    \"epic\": \"\",\n    \"task\": \"\",\n    \"subtask\": \"\"\n  },\n  \"tickets\": []\n}")
			} else {
				fmt.Println("No tickets found. Use 'jit track <ticket>' to start tracking tickets.")
			}
			return
		}

		// Load tickets
		var tickets []*types.Ticket
		for _, key := range ticketKeys {
			ticket, err := ctx.Storage.LoadTicket(key)
			if err != nil {
				continue // Skip tickets that can't be loaded
			}
			tickets = append(tickets, ticket)
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

// displayTree displays tickets in a tree format
func displayTree(tickets []*types.Ticket, currentEpic, currentTask, currentSubtask string) {
	fmt.Println("Ticket Tree:")
	fmt.Println()

	// Group tickets by type
	epics := filterTicketsByType(tickets, types.TicketTypeEpic)
	tasks := filterTicketsByType(tickets, types.TicketTypeTask)
	subtasks := filterTicketsByType(tickets, types.TicketTypeSubtask)

	// Display epics
	for _, epic := range epics {
		focusMarker := ""
		if epic.Key == currentEpic {
			focusMarker = " *"
		}
		fmt.Printf("📦 %s (%s)%s\n", epic.Key, epic.Status, focusMarker)
		fmt.Printf("   %s\n", epic.Title)

		// Show tasks under this epic
		epicTasks := filterTasksByParent(tasks, epic.Key)
		for _, task := range epicTasks {
			focusMarker = ""
			if task.Key == currentTask {
				focusMarker = " *"
			}
			fmt.Printf("   📋 %s (%s)%s\n", task.Key, task.Status, focusMarker)
			fmt.Printf("      %s\n", task.Title)

			// Show subtasks under this task
			taskSubtasks := filterSubtasksByParent(subtasks, task.Key)
			for _, subtask := range taskSubtasks {
				focusMarker = ""
				if subtask.Key == currentSubtask {
					focusMarker = " *"
				}
				fmt.Printf("      🔧 %s (%s)%s\n", subtask.Key, subtask.Status, focusMarker)
				fmt.Printf("         %s\n", subtask.Title)
			}
		}
		fmt.Println()
	}

	// Show orphan tasks (tasks without epic)
	orphanTasks := filterOrphanTasks(tasks)
	if len(orphanTasks) > 0 {
		fmt.Println("Orphan Tasks:")
		for _, task := range orphanTasks {
			focusMarker := ""
			if task.Key == currentTask {
				focusMarker = " *"
			}
			fmt.Printf("📋 %s (%s)%s\n", task.Key, task.Status, focusMarker)
			fmt.Printf("   %s\n", task.Title)

			// Show subtasks under this task
			taskSubtasks := filterSubtasksByParent(subtasks, task.Key)
			for _, subtask := range taskSubtasks {
				focusMarker = ""
				if subtask.Key == currentSubtask {
					focusMarker = " *"
				}
				fmt.Printf("   🔧 %s (%s)%s\n", subtask.Key, subtask.Status, focusMarker)
				fmt.Printf("      %s\n", subtask.Title)
			}
		}
		fmt.Println()
	}

	// Show current focus summary
	if currentEpic != "" || currentTask != "" || currentSubtask != "" {
		fmt.Println("Current Focus:")
		if currentEpic != "" {
			fmt.Printf("  Epic: %s\n", currentEpic)
		}
		if currentTask != "" {
			fmt.Printf("  Task: %s\n", currentTask)
		}
		if currentSubtask != "" {
			fmt.Printf("  Subtask: %s\n", currentSubtask)
		}
	}
}

// displayJSON displays tickets in JSON format
func displayJSON(tickets []*types.Ticket, currentEpic, currentTask, currentSubtask string) {
	// Simple JSON structure for now
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
		fmt.Printf("      \"parent\": \"%s\"\n", ticket.Relationships.ParentKey)
		if i < len(tickets)-1 {
			fmt.Printf("    },\n")
		} else {
			fmt.Printf("    }\n")
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

// filterOrphanTasks filters tasks that don't have a parent epic
func filterOrphanTasks(tasks []*types.Ticket) []*types.Ticket {
	var filtered []*types.Ticket
	for _, task := range tasks {
		if task.Relationships.ParentKey == "" {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

// GetLogCmd returns the log command
func GetLogCmd() *cobra.Command {
	return logCmd
}

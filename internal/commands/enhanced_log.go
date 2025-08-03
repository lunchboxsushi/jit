package commands

import (
	"fmt"

	"github.com/lunchboxsushi/jit/pkg/types"
	"github.com/spf13/cobra"
)

var (
	enhancedLogAllFlag    bool
	enhancedLogStatusFlag string
	enhancedLogJSONFlag   bool
	enhancedLogOrphanFlag bool
	enhancedLogTuiFlag    bool
	enhancedLogNoTuiFlag  bool
)

var enhancedLogCmd = &cobra.Command{
	Use:   "log-enhanced",
	Short: "Display enhanced ticket tree view with colors",
	Long: `Show a hierarchical view of tracked tickets with enhanced formatting and colors.
	
Examples:
  jit log-enhanced                    # Show current context tree with colors
  jit log-enhanced --all             # Show all tracked tickets
  jit log-enhanced --status "In Progress"  # Filter by status
  jit log-enhanced --orphan          # Show orphaned tasks
  jit log-enhanced --json            # Output as JSON`,
	Run: func(cmd *cobra.Command, args []string) {
		// For testing, use dummy data
		tickets := GenerateTestTickets()

		// Set current focus for testing
		currentEpic := "PROJ-100"
		currentTask := "PROJ-101"
		currentSubtask := "PROJ-103"

		if len(tickets) == 0 {
			if enhancedLogJSONFlag {
				fmt.Println("{\n  \"current_focus\": {\n    \"epic\": \"\",\n    \"task\": \"\",\n    \"subtask\": \"\"\n  },\n  \"tickets\": []\n}")
			} else {
				fmt.Println("No tickets found. Use 'jit track <ticket>' to start tracking tickets.")
			}
			return
		}

		// Filter by status if specified
		if enhancedLogStatusFlag != "" {
			tickets = filterTicketsByStatus(tickets, enhancedLogStatusFlag)
		}

		// Sort tickets by type and key
		sortTickets(tickets)

		// Display based on format
		if enhancedLogJSONFlag {
			displayEnhancedJSON(tickets, currentEpic, currentTask, currentSubtask)
		} else {
			displayEnhancedTree(tickets, currentEpic, currentTask, currentSubtask)
		}
	},
}

func init() {
	enhancedLogCmd.Flags().BoolVar(&enhancedLogAllFlag, "all", false, "Show all tickets (not just current context)")
	enhancedLogCmd.Flags().StringVar(&enhancedLogStatusFlag, "status", "", "Filter by status (e.g., 'In Progress', 'Done')")
	enhancedLogCmd.Flags().BoolVar(&enhancedLogJSONFlag, "json", false, "Output as JSON")
	enhancedLogCmd.Flags().BoolVar(&enhancedLogOrphanFlag, "orphan", false, "Show orphaned tasks")
	enhancedLogCmd.Flags().BoolVar(&enhancedLogTuiFlag, "tui", true, "Enable TUI mode (default)")
	enhancedLogCmd.Flags().BoolVar(&enhancedLogNoTuiFlag, "no-tui", false, "Force text-only output")
}

// displayEnhancedTree displays tickets in the new enhanced tree format
func displayEnhancedTree(tickets []*types.Ticket, currentEpic, currentTask, currentSubtask string) {
	fmt.Println("Enhanced Ticket Tree:")
	fmt.Println()

	// Group tickets by type
	epics := filterTicketsByType(tickets, types.TicketTypeEpic)
	tasks := filterTicketsByType(tickets, types.TicketTypeTask)
	subtasks := filterTicketsByType(tickets, types.TicketTypeSubtask)

	// Display epics and their hierarchy
	for i, epic := range epics {
		isLastEpic := i == len(epics)-1
		isFocused := epic.Key == currentEpic

		// Epic line
		epicLine := ColorizeTicket(epic, isFocused, false)
		fmt.Println(epicLine)

		// Show tasks under this epic
		epicTasks := filterTasksByParent(tasks, epic.Key)
		for _, task := range epicTasks {
			isTaskFocused := task.Key == currentTask

			// Task line with proper indentation
			taskPrefix := "   "
			if isLastEpic {
				taskPrefix = "   "
			}
			taskLine := taskPrefix + ColorizeTicket(task, isTaskFocused, false)
			fmt.Println(taskLine)

			// Show subtasks under this task
			taskSubtasks := filterSubtasksByParent(subtasks, task.Key)
			for _, subtask := range taskSubtasks {
				isSubtaskFocused := subtask.Key == currentSubtask

				// Subtask line with proper indentation
				subtaskPrefix := taskPrefix + "   "
				subtaskLine := subtaskPrefix + ColorizeTicket(subtask, isSubtaskFocused, false)
				fmt.Println(subtaskLine)
			}
		}
		fmt.Println()
	}

	// Show orphaned tasks if requested or if there are any
	orphanTasks := filterOrphanTasks(tasks)
	if len(orphanTasks) > 0 && (enhancedLogOrphanFlag || enhancedLogAllFlag) {
		fmt.Println(ColorizeHeader("== ORPHAN TASKS =="))
		fmt.Println()

		for _, task := range orphanTasks {
			isTaskFocused := task.Key == currentTask

			// Orphan task line
			taskLine := ColorizeTicket(task, isTaskFocused, true)
			fmt.Println(taskLine)

			// Show subtasks under this orphan task
			taskSubtasks := filterSubtasksByParent(subtasks, task.Key)
			for _, subtask := range taskSubtasks {
				isSubtaskFocused := subtask.Key == currentSubtask

				// Subtask line with proper indentation
				subtaskPrefix := "   "
				subtaskLine := subtaskPrefix + ColorizeTicket(subtask, isSubtaskFocused, true)
				fmt.Println(subtaskLine)
			}
		}
		fmt.Println()
	}

	// Show current focus summary
	if currentEpic != "" || currentTask != "" || currentSubtask != "" {
		fmt.Println(ColorizeHeader("Current Focus:"))
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

// displayEnhancedJSON displays tickets in JSON format
func displayEnhancedJSON(tickets []*types.Ticket, currentEpic, currentTask, currentSubtask string) {
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

// GetEnhancedLogCmd returns the enhanced log command
func GetEnhancedLogCmd() *cobra.Command {
	return enhancedLogCmd
}

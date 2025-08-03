package commands

import (
	"fmt"

	"github.com/lunchboxsushi/jit/pkg/types"
	"github.com/spf13/cobra"
)

var (
	taskNoEnrichFlag bool
	taskNoCreateFlag bool
	taskOrphanFlag   bool
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Create a new task",
	Long: `Create a new task by opening an editor with a template.
	
The command will:
1. Open an editor with a task template
2. Optionally enhance the content with AI (unless --no-enrich is used)
3. Create the task in Jira (unless --no-create is used)
4. Save the task locally
5. Set focus to the new task

By default, tasks are created under the current epic context.
Use --orphan (-o) to create a task without a parent epic.`,
	Run: func(cmd *cobra.Command, args []string) {
		flags := CreateFlags{
			NoEnrich: taskNoEnrichFlag,
			NoCreate: taskNoCreateFlag,
			Orphan:   taskOrphanFlag,
		}

		// Override context validation for orphan tasks
		validateContext := ValidateTaskContext
		if taskOrphanFlag {
			validateContext = ValidateEpicContext // No validation needed for orphan tasks
		}

		options := CreateOptions{
			TicketType:       types.TicketTypeTask,
			TemplateName:     "task.md",
			ValidateContext:  validateContext,
			SetRelationships: SetTaskRelationships,
			SuccessMessage:   "Task created successfully",
		}

		if err := CreateTicket(cmd, options, flags); err != nil {
			fmt.Printf("❌ %v\n", err)
			return
		}

		// Add orphan-specific success message
		if taskOrphanFlag {
			fmt.Println("📦 Orphan task (no parent epic)")
		}
	},
}

func init() {
	taskCmd.Flags().BoolVar(&taskNoEnrichFlag, "no-enrich", false, "Skip AI enrichment")
	taskCmd.Flags().BoolVar(&taskNoCreateFlag, "no-create", false, "Save locally only, don't create in Jira")
	taskCmd.Flags().BoolVarP(&taskOrphanFlag, "orphan", "o", false, "Create an orphan task (no parent epic)")
}

// GetTaskCmd returns the task command
func GetTaskCmd() *cobra.Command {
	return taskCmd
}

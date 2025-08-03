package commands

import (
	"fmt"

	"github.com/lunchboxsushi/jit/pkg/types"
	"github.com/spf13/cobra"
)

var (
	subtaskNoEnrichFlag bool
	subtaskNoCreateFlag bool
)

var subtaskCmd = &cobra.Command{
	Use:   "subtask",
	Short: "Create a new subtask",
	Long: `Create a new subtask by opening an editor with a template.
	
The command will:
1. Open an editor with a subtask template
2. Optionally enhance the content with AI (unless --no-enrich is used)
3. Create the subtask in Jira (unless --no-create is used)
4. Save the subtask locally
5. Set focus to the new subtask

Subtasks must be created under a task context.`,
	Run: func(cmd *cobra.Command, args []string) {
		flags := CreateFlags{
			NoEnrich: subtaskNoEnrichFlag,
			NoCreate: subtaskNoCreateFlag,
		}

		options := CreateOptions{
			TicketType:       types.TicketTypeSubtask,
			TemplateName:     "subtask.md",
			ValidateContext:  ValidateSubtaskContext,
			SetRelationships: SetSubtaskRelationships,
			SuccessMessage:   "Subtask created successfully",
		}

		if err := CreateTicket(cmd, options, flags); err != nil {
			fmt.Printf("❌ %v\n", err)
			return
		}
	},
}

func init() {
	subtaskCmd.Flags().BoolVar(&subtaskNoEnrichFlag, "no-enrich", false, "Skip AI enrichment")
	subtaskCmd.Flags().BoolVar(&subtaskNoCreateFlag, "no-create", false, "Save locally only, don't create in Jira")
}

// GetSubtaskCmd returns the subtask command
func GetSubtaskCmd() *cobra.Command {
	return subtaskCmd
}

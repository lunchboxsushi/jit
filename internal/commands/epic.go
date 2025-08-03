package commands

import (
	"fmt"

	"github.com/lunchboxsushi/jit/pkg/types"
	"github.com/spf13/cobra"
)

var (
	epicNoEnrichFlag bool
	epicNoCreateFlag bool
)

var epicCmd = &cobra.Command{
	Use:   "epic",
	Short: "Create a new epic",
	Long: `Create a new epic by opening an editor with a template.
	
The command will:
1. Open an editor with an epic template
2. Optionally enhance the content with AI (unless --no-enrich is used)
3. Create the epic in Jira (unless --no-create is used)
4. Save the epic locally
5. Set focus to the new epic`,
	Run: func(cmd *cobra.Command, args []string) {
		flags := CreateFlags{
			NoEnrich: epicNoEnrichFlag,
			NoCreate: epicNoCreateFlag,
		}

		options := CreateOptions{
			TicketType:       types.TicketTypeEpic,
			TemplateName:     "epic.md",
			ValidateContext:  ValidateEpicContext,
			SetRelationships: SetEpicRelationships,
			SuccessMessage:   "Epic created successfully",
			ParentInfo:       "Epics are top-level tickets that contain tasks and subtasks.",
		}

		if err := CreateTicket(cmd, options, flags); err != nil {
			fmt.Printf("❌ %v\n", err)
			return
		}
	},
}

func init() {
	epicCmd.Flags().BoolVar(&epicNoEnrichFlag, "no-enrich", false, "Skip AI enrichment")
	epicCmd.Flags().BoolVar(&epicNoCreateFlag, "no-create", false, "Save locally only, don't create in Jira")
}

// GetEpicCmd returns the epic command
func GetEpicCmd() *cobra.Command {
	return epicCmd
}

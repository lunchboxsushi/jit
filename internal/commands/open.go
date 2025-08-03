package commands

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [ticket-key]",
	Short: "Open Jira ticket in browser",
	Long: `Open a Jira ticket in your default web browser. If no ticket is specified, uses current focus.
	
Examples:
  jit open                    # Open current focus in browser
  jit open SRE-1234          # Open specific ticket in browser`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize command context
		ctx, err := InitializeCommand()
		if err != nil {
			HandleError(err, "Failed to initialize")
			return
		}

		var ticketKey string

		// Get ticket key from args or current focus
		if len(args) > 0 {
			ticketKey = args[0]
		} else {
			// Try to get from current focus
			currentEpic, _ := ctx.ContextManager.GetCurrentEpic()
			currentTask, _ := ctx.ContextManager.GetCurrentTask()
			currentSubtask, _ := ctx.ContextManager.GetCurrentSubtask()

			// Prefer subtask > task > epic
			if currentSubtask != "" {
				ticketKey = currentSubtask
			} else if currentTask != "" {
				ticketKey = currentTask
			} else if currentEpic != "" {
				ticketKey = currentEpic
			} else {
				fmt.Println("No ticket specified and no current focus.")
				fmt.Println("Use 'jit focus <ticket>' to set focus or specify a ticket key.")
				return
			}
		}

		// Validate ticket exists
		if !ctx.Storage.Exists(ticketKey) {
			fmt.Printf("Ticket %s not found in local storage.\n", ticketKey)
			fmt.Printf("Use 'jit track %s' to track it first.\n", ticketKey)
			return
		}

		// Generate URL
		url := fmt.Sprintf("%s/browse/%s", ctx.Config.Jira.URL, ticketKey)

		// Open in browser
		if err := openBrowser(url); err != nil {
			HandleError(err, "Failed to open browser")
			return
		}

		fmt.Printf("Opened %s in browser\n", ticketKey)
	},
}

// openBrowser opens a URL in the default browser
func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", url)
	case "linux":
		// Try xdg-open first, then sensible-browser
		if _, err := exec.LookPath("xdg-open"); err == nil {
			cmd = exec.Command("xdg-open", url)
		} else if _, err := exec.LookPath("sensible-browser"); err == nil {
			cmd = exec.Command("sensible-browser", url)
		} else {
			return fmt.Errorf("no browser opener available (install xdg-utils)")
		}
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Run()
}

// GetOpenCmd returns the open command
func GetOpenCmd() *cobra.Command {
	return openCmd
}

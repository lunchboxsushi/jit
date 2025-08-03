package commands

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	linkShortFlag bool
)

var linkCmd = &cobra.Command{
	Use:   "link [ticket-key]",
	Short: "Get Jira URL for ticket",
	Long: `Output the Jira URL for a ticket. If no ticket is specified, uses current focus.
	
Examples:
  jit link                    # Get URL for current focus
  jit link SRE-1234          # Get URL for specific ticket
  jit link -s                # Short format (ticket key only)
  jit link SRE-1234 -s       # Short format for specific ticket`,
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

		// Output based on format
		if linkShortFlag {
			fmt.Println(ticketKey)
		} else {
			fmt.Println(url)
		}

		// Try to copy to clipboard
		if err := copyToClipboard(url); err != nil {
			// Non-fatal error, just warn
			PrintWarning("Failed to copy to clipboard (URL still displayed)")
		} else {
			if !linkShortFlag {
				PrintInfo("URL copied to clipboard")
			}
		}
	},
}

func init() {
	linkCmd.Flags().BoolVarP(&linkShortFlag, "short", "s", false, "Output ticket key only")
}

// copyToClipboard copies text to clipboard using platform-specific commands
func copyToClipboard(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("pbcopy")
	case "linux":
		// Try xclip first, then xsel
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if _, err := exec.LookPath("xsel"); err == nil {
			cmd = exec.Command("xsel", "--clipboard", "--input")
		} else {
			return fmt.Errorf("no clipboard tool available (install xclip or xsel)")
		}
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	// Set up pipe for input
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	// Start command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start clipboard command: %v", err)
	}

	// Write text to stdin
	if _, err := stdin.Write([]byte(text)); err != nil {
		return fmt.Errorf("failed to write to clipboard: %v", err)
	}

	// Close stdin and wait for command to finish
	stdin.Close()
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("clipboard command failed: %v", err)
	}

	return nil
}

// GetLinkCmd returns the link command
func GetLinkCmd() *cobra.Command {
	return linkCmd
}

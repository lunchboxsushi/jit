package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script for specified shell",
	Long: `Generate completion script for jit CLI.

To load completions:

Bash:
  $ source <(jit completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ jit completion bash > ~/.local/share/bash-completion/completions/jit
  # macOS:
  $ jit completion bash > /usr/local/etc/bash_completion.d/jit

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ jit completion zsh > "${fpath[1]}/_jit"

  # You will need to start a new shell for this setup to take effect.

Fish:
  $ jit completion fish | source

  # To load completions for each session, execute once:
  $ jit completion fish > ~/.config/fish/completions/jit.fish

PowerShell:
  PS> jit completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> jit completion powershell > jit.ps1
  # and source this file from your PowerShell profile.
`,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shell := args[0]
		rootCmd := cmd.Root()

		var err error
		switch shell {
		case "bash":
			err = rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			err = rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			err = rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			err = rootCmd.GenPowerShellCompletion(os.Stdout)
		default:
			fmt.Fprintf(os.Stderr, "Unsupported shell type: %s\n", shell)
			os.Exit(1)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating completion script: %v\n", err)
			os.Exit(1)
		}
	},
}

// GetCompletionCmd returns the completion command
func GetCompletionCmd() *cobra.Command {
	return completionCmd
}

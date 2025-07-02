package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:
  $ source <(flow completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ flow completion bash > /etc/bash_completion.d/flow
  # macOS:
  $ flow completion bash > /usr/local/etc/bash_completion.d/flow

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ flow completion zsh > "${fpath[1]}/_flow"

  # You will need to start a new shell for this setup to take effect.

Fish:
  $ flow completion fish | source

  # To load completions for each session, execute once:
  $ flow completion fish > ~/.config/fish/completions/flow.fish

PowerShell:
  PS> flow completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> flow completion powershell > flow.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		switch args[0] {
		case "bash":
			err = rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			err = rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			err = rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			err = rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

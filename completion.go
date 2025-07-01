package main

import (
	"fmt"
	"os"
)

func handleCompletion() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: flow completion [bash|zsh]")
		os.Exit(1)
	}

	shell := os.Args[2]
	switch shell {
	case "bash":
		fmt.Print(bashCompletion)
	case "zsh":
		fmt.Print(zshCompletion)
	default:
		fmt.Fprintf(os.Stderr, "Unsupported shell: %s\n", shell)
		os.Exit(1)
	}
}

const bashCompletion = `
_flow_completions() {
    COMPREPLY=()
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local commands="start status pause resume end completion"
    local shells="bash zsh"

    if [ "${COMP_CWORD}" -eq 1 ]; then
        COMPREPLY=( $(compgen -W "${commands}" -- "${cur}") )
    elif [ "${COMP_CWORD}" -eq 2 ] && [ "${COMP_WORDS[1]}" = "completion" ]; then
        COMPREPLY=( $(compgen -W "${shells}" -- "${cur}") )
    fi
}
complete -F _flow_completions flow
`

const zshCompletion = `
#compdef flow

_flow() {
    local -a commands
    commands=(
        'start:Begin a new deep work session'
        'status:Check the current session'
        'pause:Pause the current session'
        'resume:Resume a paused session'
        'end:End the current session'
        'completion:Generate completion script'
    )
    _describe 'command' commands

    case $words[1] in
        completion)
            local -a shells
            shells=(
                'bash:Generate Bash completion script'
                'zsh:Generate Zsh completion script'
            )
            _describe 'shell' shells
            ;;
    esac
}

_flow "$@"
`
